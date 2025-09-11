package judge

import (
	"JudgeCore/internal/config"
	"context"
	"log"
	"os/exec"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// JudgePool 容器池结构
type JudgePool struct {
	pool          chan string
	mutex         sync.Mutex
	sandboxConfig config.Sandbox
	codeDir       string
	containerIDs  sync.Map
}

// NewJudgePool 创建一个新的 JudgePool 实例
func NewJudgePool(sandboxConfig config.Sandbox, codeDir string) *JudgePool {
	return &JudgePool{
		sandboxConfig: sandboxConfig,
		codeDir:       codeDir,
	}
}

// Initialize 预热容器池
func (p *JudgePool) Initialize(n int) {
	p.pool = make(chan string, n)
	for i := 0; i < n; i++ {
		containerID, err := p.startContainer()
		if err != nil {
			log.Printf("[FeasOJ] Error starting container during preheat: %v", err)
			continue
		}
		p.pool <- containerID
	}
	log.Printf("[FeasOJ] Preheated %d containers", len(p.pool))
}

// AcquireContainer 从池中获取一个空闲容器（若池为空则阻塞等待）
func (p *JudgePool) AcquireContainer() string {
	containerID := <-p.pool
	return containerID
}

// ReleaseContainer 将容器归还到池中
func (p *JudgePool) ReleaseContainer(containerID string) {
	// 清理容器中所有残留的任务目录
	if err := p.resetContainer(containerID); err != nil {
		log.Printf("[FeasOJ] Reset failed for container %s: %v, terminating it", containerID, err)
		p.containerIDs.Delete(containerID)
		go TerminateContainer(containerID)

		// 尝试启动一个新容器替换
		newContainerID, err := p.startContainer()
		if err != nil {
			log.Printf("[FeasOJ] Failed to start new container to replace failed one: %v", err)
			return
		}
		containerID = newContainerID
	}

	// 将有效容器归还
	p.mutex.Lock()
	defer p.mutex.Unlock()
	select {
	case p.pool <- containerID:
	default:
		log.Printf("[FeasOJ] Pool is full or closed. Terminating extra container %s", containerID)
		p.containerIDs.Delete(containerID)
		go TerminateContainer(containerID)
	}
}

// Shutdown 在服务关闭时终止池中所有容器
func (p *JudgePool) Shutdown() {
	p.mutex.Lock()
	close(p.pool)
	p.mutex.Unlock()

	log.Println("[FeasOJ] Shutting down container pool...")
	p.containerIDs.Range(func(key, value interface{}) bool {
		containerID := key.(string)
		TerminateContainer(containerID)
		log.Printf("[FeasOJ] Terminated container %s", containerID)
		return true
	})
}

// resetContainer 用于在归还容器到池中前清理所有残留的任务目录
func (p *JudgePool) resetContainer(containerID string) error {
	resetCmd := exec.Command("docker", "exec", containerID, "sh", "-c", "find /workspace -maxdepth 1 -type d -name 'task_*' -exec rm -rf {} +")
	if err := resetCmd.Run(); err != nil {
		log.Printf("[FeasOJ] Error resetting container %s: %v", containerID, err)
		return err
	}
	return nil
}

// startContainer 启动一个新的沙盒容器并返回其ID
func (p *JudgePool) startContainer() (string, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	containerConfig := &container.Config{
		Image: "judgecore:latest",
		Cmd:   []string{"sh"},
		Tty:   true,
	}

	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			Memory:    p.sandboxConfig.Memory,
			NanoCPUs:  int64(p.sandboxConfig.NanoCPUs * 1e9),
			CPUShares: p.sandboxConfig.CPUShares,
		},
		Binds: []string{
			p.codeDir + ":/workspace", // 挂载文件夹
		},
		AutoRemove: true, // 容器退出后自动删除
		CapDrop:    []string{"ALL"},
	}

	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	p.containerIDs.Store(resp.ID, true)
	return resp.ID, nil
}
