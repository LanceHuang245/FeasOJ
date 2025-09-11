package main

import (
	"FeasOJ/internal/config"
	"FeasOJ/internal/global"
	"FeasOJ/internal/judge"
	"FeasOJ/internal/utils"
	"FeasOJ/internal/utils/scheduler"
	"FeasOJ/internal/utils/sql"
	"FeasOJ/server"
	"bufio"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("[FeasOJ] The server is starting...")
	global.CurrentDir, _ = os.Getwd()

	// 定义目录映射
	dirs := map[string]*string{
		"certificate": &global.CertDir,
		"avatars":     &global.AvatarsDir,
		"logs":        &global.LogDir,
		"docs":        &global.DocsDir,
	}

	// 遍历map，设置路径并创建不存在的目录
	for name, dir := range dirs {
		*dir = filepath.Join(global.CurrentDir, name)
		if _, err := os.Stat(*dir); os.IsNotExist(err) {
			err := os.Mkdir(*dir, os.ModePerm)
			if err != nil {
				log.Panicf("[FeasOJ] Failed to create directory: %v", err)
				return
			}
		}
	}

	// 初始化Logger
	logFile, err := utils.InitializeLogger()
	if err != nil {
		log.Fatalf("[FeasOJ] Failed to initialize logger: %v", err)
	}
	defer utils.CloseLogger(logFile)

	// 初始化配置
	if err := config.InitConfig(); err != nil {
		log.Fatalf("[FeasOJ] Failed to initialize config: %v", err)
	}

	// 初始化数据库
	global.DB = utils.ConnectSql()
	if utils.InitTable() {
		log.Println("[FeasOJ] Database initialization complete")
	} else {
		log.Println("[FeasOJ] Database initialization failed")
		return
	}

	// 初始化管理员账户
	if sql.GetAdminUser(1) {
		log.Println("[FeasOJ] The administrator account already exists and will continue")
	} else {
		sql.Register(utils.InitAdminAccount())
	}

	// 测试邮箱模块是否正常
	if utils.TestSend(config.GlobalConfig.Mail) {
		log.Println("[FeasOJ] Email service initialization complete")
	} else {
		log.Println("[FeasOJ] Email service startup failed, please check the configuration")
		return
	}

	// 测试Redis连接
	if utils.ConnectRedis() == nil {
		log.Println("[FeasOJ] Redis connection failed, please check the configuration")
		return
	} else {
		log.Println("[FeasOJ] Redis connection successful")
	}

	// 测试ImageGuard连接
	if config.GlobalConfig.Features.ImageGuardEnabled {
		if utils.ImageGuardPing() {
			log.Println("[FeasOJ] ImageGuard service connection successful")
		} else {
			log.Panicln("[FeasOJ] ImageGuard service connection failed")
			return
		}
	}

	// 测试ProfanityDetector连接
	if config.GlobalConfig.Features.ProfanityDetectorEnabled {
		if utils.ProfanityDetectorPing() {
			log.Println("[FeasOJ] ProfanityDetector service connection successful")
		} else {
			log.Panicln("[FeasOJ] ProfanityDetector service connection failed")
			return
		}
	}

	// 测试JudgeCore连接
	if judge.JudgeCorePing() {
		log.Println("[FeasOJ] JudgeCore service connection successful")
	} else {
		log.Panicln("[FeasOJ] JudgeCore service connection failed")
		return
	}

	// 判题结果消息队列
	go judge.ConsumeJudgeResults()

	// 启用竞赛状态调度器
	go scheduler.ScheduleCompetitionStatus()

	// 启动服务器
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	server.LoadRouter(r)

	// 挂载头像文件夹
	r.StaticFS("/api/v1/avatar", http.Dir(global.AvatarsDir))

	// 挂载文档文件夹
	r.StaticFS("/api/v1/docs", http.Dir(global.DocsDir))

	startServer := func(protocol, address, certFile, keyFile string) {
		for {
			var err error
			if protocol == "http" {
				err = r.Run(address)
			} else {
				err = r.RunTLS(address, certFile, keyFile)
			}
			if err != nil {
				log.Printf("[FeasOJ] Server start error: %v\n", err)
				os.Exit(0)
			}
		}
	}

	if config.GlobalConfig.Server.EnableHTTPS {
		go startServer("https", config.GlobalConfig.Server.Address, config.GlobalConfig.Server.CertPath, config.GlobalConfig.Server.KeyPath)
	} else {
		go startServer("http", config.GlobalConfig.Server.Address, "", "")
	}

	log.Println("[FeasOJ] Server is running on", config.GlobalConfig.Server.Address, "Https Status:", config.GlobalConfig.Server.EnableHTTPS)

	// 监听终端输入
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for scanner.Scan() {
			input := scanner.Text()
			if input == "quit" || input == "exit" {
				log.Println("[FeasOJ] Server is shutting down...")
				os.Exit(0)
			}
		}
	}()

	// 监听系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("[FeasOJ] Server is shutting down...")
}
