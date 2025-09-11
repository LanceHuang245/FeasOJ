package global

import (
	"gorm.io/gorm"
)

// 全局变量 - 本地配置文件路径
var LogDir string
var AvatarsDir string
var CurrentDir string
var CertDir string
var DocsDir string

// 全局变量 - 数据库连接对象
var DB *gorm.DB

// 全局变量 - 微服务地址
var ImageGuardAddr string
var ProfanityDetectorAddr string
var JudgeCoreAddr string
