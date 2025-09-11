English | [简体中文](README_CN.md)
<p align="center">
    <a href="https://github.com/ClaretWheel1481/FeasOJ-Backend">
        <img src="public/logo.png" height="200"/>
    </a>
</p>

> For unified project management and standardization, Backend will be migrating to [FeasOJ](https://github.com/LanceHuang245/FeasOJ)

# FeasOJ
### Project Description
FeasOJ is an online programming practice platform based on Vue and Golang, supporting multi-languages, discussion forums, contests and other features, aiming to provide users with a convenient and efficient learning and practice environment.
<br>
[FeasOJ-Frontend](https://github.com/ClaretWheel1481/FeasOJ)
[FeasOJ-JudgeCore](https://github.com/ClaretWheel1481/FeasOJ-JudgeCore)
[ImageGuard](https://github.com/ClaretWheel1481/ImageGuard)
[Profanity Detector](https://github.com/ClaretWheel1481/ProfanityDetector)
[API Document(CN)](https://claret-feasoj.apifox.cn)
[API Document(EN)](https://claret-feasoj.apifox.cn/en/)

### Notice
The project is expanding as well as being microserviced, so watch out for the Readme.
If you find any bugs, please open an issue.

### Project Structure
```
FeasOJ-Backend
│ 
├─config.json          # Configuration file
├─internal
│  ├─gin
│  ├─global
│  ├─judge
│  ├─middlewares
│  ├─router
│  └─utils
│      ├─locales   # i18n files
│      ├─scheduler
│      └─sql
├─go.mod
└─main.go    # Main entry file
```

### Environment
- Golang 1.25.0
- Redis
- MySQL 8.0+
- RabbitMQ
- Consul

### How to run
1. Clone this repository and [ImageGuard](https://github.com/ClaretWheel1481/ImageGuard)(Optional), [Profanity Detector](https://github.com/ClaretWheel1481/ProfanityDetector)(Optional), and [FeasOJ-JudgeCore](https://github.com/ClaretWheel1481/FeasOJ-JudgeCore)
2. Install Docker, RabbitMQ, Consul, Redis and MySQL
3. Start `MySQL`, `Redis`, `Docker`, `Consul`, `RabbitMQ`, `ImageGuard`(Optional), `ProfanityDetector`(Optional) Services
4. Run `cd src` and `go mod tidy` Install dependencies
5. On first run, the program will automatically create a `config.json` configuration file. Please edit this file to configure database, Redis, email and other information
6. Run `go run main.go` to start the backend server

### Configuration
The program uses JSON configuration files to manage all configuration items. On first run, it will automatically create a `config.json` file containing:

- **Server Configuration**: Listen address, HTTPS settings, certificate paths
- **Database Configuration**: MySQL connection information, connection pool settings
- **Redis Configuration**: Redis connection address and password
- **Email Configuration**: SMTP server information
- **Microservice Configuration**: RabbitMQ, Consul addresses
- **Feature Switches**: [ImageGuard](https://github.com/ClaretWheel1481/ImageGuard), [Profanity Detector](https://github.com/ClaretWheel1481/ProfanityDetector), etc.

For detailed configuration instructions, please refer to [CONFIG README.md](/docs/CONFIG_README_EN.md)

### Acknowledgments
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
