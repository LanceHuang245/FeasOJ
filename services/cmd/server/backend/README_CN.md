[English](README.md) | 简体中文
<p align="center">
    <a href="https://github.com/ClaretWheel1481/FeasOJ-Backend">
        <img src="public/logo.png" height="200"/>
    </a>
</p>

> 为统一项目管理与规范化，JudgeCore将迁移至[FeasOJ](https://github.com/LanceHuang245/FeasOJ)

# FeasOJ-Backend
### 项目简介
FeasOJ 是一个基于 Vue 和 Golang 的在线编程练习平台，支持多国语言、讨论区、竞赛等功能，旨在为用户提供一个方便、高效的学习和练习环境。
<br>
[FeasOJ-Frontend](https://github.com/ClaretWheel1481/FeasOJ)
[FeasOJ-JudgeCore](https://github.com/ClaretWheel1481/FeasOJ-JudgeCore)
[ImageGuard](https://github.com/ClaretWheel1481/ImageGuard)
[Profanity Detector](https://github.com/ClaretWheel1481/ProfanityDetector)
[API Document(CN)](https://claret-feasoj.apifox.cn)
[API Document(EN)](https://claret-feasoj.apifox.cn/en/)

### 注意
该项目正在不断扩大以及微服务化，请注意阅读Readme。
如果你找到任何Bug请发布Issue告诉我。

### 项目结构
```
FeasOJ-Backend
│ 
├─config.json          # 配置文件
├─internal
│  ├─gin
│  ├─global
│  ├─judge
│  ├─middlewares
│  ├─router
│  └─utils
│      ├─locales
│      ├─scheduler
│      └─sql
├─go.mod
└─main.go
```

### 环境
- Golang 1.24.4
- Redis
- MySQL 8.0+
- RabbitMQ
- Consul

### 如何运行
1. 克隆此库以及[ImageGuard](https://github.com/ClaretWheel1481/ImageGuard)(可选)和[Profanity Detector](https://github.com/ClaretWheel1481/ProfanityDetector)(可选)以及[FeasOJ-JudgeCore](https://github.com/ClaretWheel1481/FeasOJ-JudgeCore)(必须)
2. 安装 Docker、MySQL、Redis、Consul和RabbitMQ
3. 运行 `MySQL`、`Redis`、`Docker`、`Consul`、`RabbitMQ`、`ImageGuard`(可选) 和 `Profanity Detector`(可选)
4. 运行 `cd src` 和 `go mod tidy` 下载依赖
5. 首次运行程序会自动创建 `config.json` 配置文件，请编辑该文件配置数据库、Redis、邮件等信息
6. 运行 `go run main.go` 启动后端服务器

### 配置说明
程序使用 JSON 配置文件管理所有配置项。首次运行时会自动创建 `config.json` 文件，包含以下配置：

- **服务器配置**: 监听地址、HTTPS设置、证书路径
- **数据库配置**: MySQL连接信息、连接池设置
- **Redis配置**: Redis连接地址和密码
- **邮件配置**: SMTP服务器信息
- **微服务配置**: RabbitMQ、Consul地址
- **功能开关**: [ImageGuard](https://github.com/ClaretWheel1481/ImageGuard)、[Profanity Detector](https://github.com/ClaretWheel1481/ProfanityDetector)

更多关于配置文件信息，请前往 [CONFIG README.md](/docs/CONFIG_README_EN.md) 查看

### 致谢
- [Go](https://github.com/golang/go)
- [Gin](https://github.com/gin-gonic/gin)
- [gorm](https://github.com/go-gorm/gorm)
- [jwt-go](https://github.com/golang-jwt/jwt)
- [docker](https://github.com/moby/moby)
- [gomail](https://github.com/go-gomail/gomail)
- [go-redis](https://github.com/redis/go-redis)
- [go-i18n](https://github.com/nicksnyder/go-i18n)
- [gocron](https://github.com/go-co-op/gocron)
- [amqp091-go](https://github.com/rabbitmq/amqp091-go)