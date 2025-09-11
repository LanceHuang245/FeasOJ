package judge

import (
	"FeasOJ/internal/global"
	"FeasOJ/internal/utils"
	"FeasOJ/server/handler"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// 发送评测结果到消息队列
func ConsumeJudgeResults() {
	for {
		conn, ch, err := utils.ConnectRabbitMQ()
		if err != nil {
			log.Println("[FeasOJ] Failed to connect to RabbitMQ, retrying in 5 seconds:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// 声明结果队列
		q, err := ch.QueueDeclare(
			"judgeResults", // 队列名称
			true,           // 持久化
			false,          // 自动删除
			false,          // 排他性
			false,          // 不等待
			nil,            // 参数
		)
		if err != nil {
			log.Println("[FeasOJ] Failed to declare a message queue, retrying in 5 seconds:", err)
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		msgs, err := ch.Consume(
			q.Name, // 队列名称
			"",     // 消费者标签
			true,   // 自动应答
			false,  // 排他性
			false,  // 不等待
			false,  // 参数
			nil,
		)
		if err != nil {
			log.Println("[FeasOJ] Failed to register consumer, retrying in 5 seconds:", err)
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		// 消费消息
		errChan := make(chan error, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[FeasOJ] Panic in judge result consumer: %v", r)
					errChan <- fmt.Errorf("panic: %v", r)
				}
			}()
			for d := range msgs {
				var result global.JudgeResultMessage
				if err := json.Unmarshal(d.Body, &result); err != nil {
					log.Printf("[FeasOJ] Error decoding result: %v", err)
					continue
				}
				// 发送 SSE 通知
				if client, ok := handler.Clients[fmt.Sprint(result.UserID)]; ok {
					lang := client.Lang
					tag := language.Make(lang)
					langBundle := utils.InitI18n()
					localizer := i18n.NewLocalizer(langBundle, tag.String())
					message, _ := localizer.Localize(&i18n.LocalizeConfig{
						MessageID: "problem_completed",
						TemplateData: map[string]any{
							"PID": result.ProblemID,
						},
					})
					client.MessageChan <- message
				}
			}
			errChan <- fmt.Errorf("[FeasOJ] RabbitMQ channel closed")
		}()

		err = <-errChan
		log.Println("[FeasOJ] Judge result consumer error, reconnecting:", err)
		ch.Close()
		conn.Close()
		time.Sleep(5 * time.Second)
	}
}
