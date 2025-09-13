package judge

import (
	"FeasOJ/app/judgecore/internal/config"
	"FeasOJ/app/judgecore/internal/global"
	"FeasOJ/app/judgecore/internal/utils"
	"FeasOJ/pkg/databases/repository"
	"FeasOJ/pkg/structs"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type Task struct {
	UID  string
	PID  int
	Name string
}

// ProcessJudgeTasks 函数用于处理判题任务
func ProcessJudgeTasks(rmqConfig config.RabbitMQ, db *gorm.DB, pool *JudgePool) {
	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error

	for {
		conn, ch, err = utils.ConnectRabbitMQ(rmqConfig)
		if err != nil {
			log.Println("[FeasOJ] RabbitMQ connect error, retrying in 3s: ", err)
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("[FeasOJ] RabbitMQ connected")
		break
	}
	defer conn.Close()
	defer ch.Close()

	taskChan := make(chan Task)
	var wg sync.WaitGroup

	for i := 0; i < pool.sandboxConfig.MaxConcurrent; i++ {
		wg.Add(1)
		go worker(taskChan, ch, &wg, db, pool)
	}

	for {
		// 获取队列中的任务
		msgs, err := ch.Consume(
			"judgeTask", // 队列名称
			"",          // 消费者标签
			true,        // 自动应答
			false,       // 是否排他
			false,       // 是否持久化
			false,       // 是否等待
			nil,         // 额外参数
		)
		if err != nil {
			log.Println("[FeasOJ] Failed to start consuming, retrying in 3s: ", err)
			time.Sleep(3 * time.Second)
			// 重新连接
			conn.Close()
			ch.Close()
			for {
				conn, ch, err = utils.ConnectRabbitMQ(rmqConfig)
				if err != nil {
					log.Println("[FeasOJ] RabbitMQ reconnect error, retrying in 3s: ", err)
					time.Sleep(3 * time.Second)
					continue
				}
				break
			}
			continue
		}

		for msg := range msgs {
			taskData := string(msg.Body)
			parts := strings.Split(taskData, "_")
			if len(parts) < 2 {
				log.Printf("[FeasOJ] Invalid task data format: %s", taskData)
				continue
			}
			uid := parts[0]
			pidStr := strings.Split(parts[1], ".")[0]
			pid, _ := strconv.Atoi(pidStr)

			taskChan <- Task{UID: uid, PID: pid, Name: taskData}
		}
		log.Println("[FeasOJ] RabbitMQ channel closed. Exiting task processor.")
		break
	}

	close(taskChan)
	wg.Wait()
}

// worker 使用容器池执行任务
func worker(taskChan chan Task, ch *amqp.Channel, wg *sync.WaitGroup, db *gorm.DB, pool *JudgePool) {
	defer wg.Done()

	for task := range taskChan {
		problem, err := repository.SelectProblemByPid(db, task.PID)
		if err != nil {
			log.Printf("[FeasOJ] Failed to get problem info for PID %d: %v", task.PID, err)
			continue
		}

		testCases := repository.SelectTestCasesByPid(db, task.PID)
		if len(testCases) == 0 {
			log.Printf("[FeasOJ] No test cases found for PID %d", task.PID)
			repository.ModifyJudgeStatus(db, task.PID, task.UID, global.SystemError)
			continue
		}

		containerID := pool.AcquireContainer()
		pool.containerIDs.Store(task.Name, containerID)

		result := CompileAndRun(task.Name, containerID, problem, testCases)
		repository.ModifyJudgeStatus(db, task.PID, task.UID, result)

		resultMsg := structs.JudgeResultMessage{
			UserID:    task.UID,
			ProblemID: task.PID,
			Status:    result,
		}

		if err := utils.PublishJudgeResult(ch, resultMsg); err != nil {
			log.Printf("[FeasOJ] Failed to publish result: %v", err)
		}

		pool.ReleaseContainer(containerID)
		pool.containerIDs.Delete(task.Name)
	}
}
