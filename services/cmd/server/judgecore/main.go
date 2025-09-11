package main

import (
	"JudgeCore/internal/config"
	"JudgeCore/internal/judge"
	"JudgeCore/internal/utils"
	"JudgeCore/server"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("[FeasOJ] Failed to get current directory: %v", err)
	}

	// 定义并创建必要的目录
	logDir := filepath.Join(currentDir, "logs")
	codeDir := filepath.Join(currentDir, "codefiles")

	certDir := filepath.Join(currentDir, "certificate")
	for _, dir := range []string{logDir, codeDir, certDir} {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.Mkdir(dir, os.ModePerm)
		}
	}

	// 初始化Logger
	logFile, err := utils.InitializeLogger(logDir)
	if err != nil {
		log.Fatalf("[FeasOJ] Failed to initialize logger: %v", err)
	}
	defer utils.CloseLogger(logFile)

	// 加载配置
	cfg, err := config.LoadConfig(currentDir)
	if err != nil {
		log.Fatalf("[FeasOJ] Failed to load config: %v", err)
	}

	// 初始化数据库
	db, err := utils.ConnectSql(cfg.Database)
	if err != nil {
		log.Fatalf("[FeasOJ] MySQL initialization failed: %v", err)
	}
	log.Println("[FeasOJ] MySQL initialization complete")

	// 初始化Consul客户端
	consulConfig := api.DefaultConfig()
	consulConfig.Address = cfg.Consul.Address
	log.Println("[FeasOJ] Connecting to Consul...")
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("[FeasOJ] Error connecting to Consul: %v", err)
	}

	// 构建沙盒镜像
	if !judge.BuildImage(currentDir) {
		log.Fatalf("[FeasOJ] SandBox builds fail, please make sure Docker is running and up to date")
	}
	log.Println("[FeasOJ] SandBox builds successfully")

	// 初始化并预热容器池
	judgePool := judge.NewJudgePool(cfg.Sandbox, codeDir)
	judgePool.Initialize(cfg.Sandbox.MaxConcurrent)

	// 启动Judge任务处理协程
	go judge.ProcessJudgeTasks(cfg.RabbitMQ, db, judgePool)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	server.LoadRouter(r, db, judgePool, codeDir)

	go func() {
		serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port)
		var err error
		if cfg.Server.EnableHTTPS {
			certPath := filepath.Join(certDir, filepath.Base(cfg.Server.CertPath))
			keyPath := filepath.Join(certDir, filepath.Base(cfg.Server.KeyPath))
			err = r.RunTLS(serverAddr, certPath, keyPath)
		} else {
			err = r.Run(serverAddr)
		}
		if err != nil {
			log.Fatalf("[FeasOJ] Server start error: %v\n", err)
		}
	}()

	// 注册服务
	if err := utils.RegService(consulClient, cfg.Consul, cfg.Server); err != nil {
		log.Fatalf("[FeasOJ] Failed to register service with Consul: %v", err)
	}

	// 优雅地关闭
	gracefulShutdown(logFile, judgePool)
}

func gracefulShutdown(logFile *os.File, pool *judge.JudgePool) {
	// 监听终端输入或中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if scanner.Text() == "exit" || scanner.Text() == "EXIT" {
				close(quit)
				return
			}
		}
	}()

	log.Println("[FeasOJ] Input 'exit' or Ctrl+C to stop the server")
	<-quit

	log.Println("[FeasOJ] The server is shutting down, please be patient to wait for the container to be closed")
	pool.Shutdown()
	utils.CloseLogger(logFile)
	os.Exit(0)
}
