package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Config 全局配置结构体
type Config struct {
	Server   ServerConfig   `json:"server"`
	RabbitMQ RabbitMQConfig `json:"rabbitmq"`
	Consul   ConsulConfig   `json:"consul"`
	Features FeaturesConfig `json:"features"`
	MySQL    MySQLConfig    `json:"mysql"`
	Redis    RedisConfig    `json:"redis"`
	Mail     MailConfig     `json:"mail"`
	JWT      JWTConfig      `json:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Address     string `json:"address"`
	EnableHTTPS bool   `json:"enable_https"`
	CertPath    string `json:"cert_path"`
	KeyPath     string `json:"key_path"`
}

// RabbitMQConfig RabbitMQ配置
type RabbitMQConfig struct {
	Address string `json:"address"`
}

// ConsulConfig Consul配置
type ConsulConfig struct {
	Address string `json:"address"`
}

// FeaturesConfig 功能开关配置
type FeaturesConfig struct {
	ImageGuardEnabled        bool `json:"image_guard_enabled"`
	ProfanityDetectorEnabled bool `json:"profanity_detector_enabled"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
	MaxLifeTime  int    `json:"max_life_time"`
	DbAddress    string `json:"db_address"`
	DbName       string `json:"db_name"`
	DbUser       string `json:"db_user"`
	DbPassword   string `json:"db_password"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

// MailConfig 邮件配置
type MailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	SigningMethod    string `json:"signing_method"`
	TokenExpireHours int    `json:"token_expire_hours"`
}

// 全局配置实例
var GlobalConfig *Config

// 初始化配置
func InitConfig() error {
	configPath := "config.json"

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Println("[FeasOJ] Configuration file does not exist, creating default configuration file...")
		if err := createDefaultConfig(configPath); err != nil {
			return fmt.Errorf("failed to create default configuration file: %v", err)
		}
		log.Println("[FeasOJ] Default configuration file created, please edit config.json file and restart the program")
		return fmt.Errorf("please edit config.json file to configure database and other information")
	}

	// 读取配置文件
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析JSON配置
	GlobalConfig = &Config{}
	if err := json.Unmarshal(configData, GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证配置
	if err := validateConfig(GlobalConfig); err != nil {
		return fmt.Errorf("配置验证失败: %v", err)
	}

	log.Println("[FeasOJ] Configuration file loaded successfully")
	return nil
}

// 创建默认配置文件
func createDefaultConfig(configPath string) error {
	defaultConfig := &Config{
		Server: ServerConfig{
			Address:     "127.0.0.1:37882",
			EnableHTTPS: true,
			CertPath:    "./certificate/fullchain.pem",
			KeyPath:     "./certificate/privkey.key",
		},
		RabbitMQ: RabbitMQConfig{
			Address: "amqp://USERNAME:PASSWORD@IP:PORT/",
		},
		Consul: ConsulConfig{
			Address: "localhost:8500",
		},
		Features: FeaturesConfig{
			ImageGuardEnabled:        true,
			ProfanityDetectorEnabled: true,
		},
		MySQL: MySQLConfig{
			MaxOpenConns: 240,
			MaxIdleConns: 100,
			MaxLifeTime:  32,
			DbAddress:    "localhost:3306",
			DbName:       "feasoj",
			DbUser:       "root",
			DbPassword:   "password",
		},
		Redis: RedisConfig{
			Address:  "localhost:6379",
			Password: "",
		},
		Mail: MailConfig{
			Host:     "smtp.qq.com",
			Port:     465,
			User:     "your-email@qq.com",
			Password: "your-password",
		},
		JWT: JWTConfig{
			SigningMethod:    "HS256",
			TokenExpireHours: 720,
		},
	}

	configData, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, configData, 0644)
}

// 验证配置
func validateConfig(config *Config) error {
	if config.Server.Address == "" {
		return fmt.Errorf("服务器地址不能为空")
	}
	if config.MySQL.DbAddress == "" || config.MySQL.DbName == "" || config.MySQL.DbUser == "" {
		return fmt.Errorf("MySQL配置不完整")
	}
	if config.Redis.Address == "" {
		return fmt.Errorf("Redis地址不能为空")
	}
	return nil
}

// 获取MySQL连接字符串
func GetMySQLDSN() string {
	if GlobalConfig == nil {
		return ""
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai",
		GlobalConfig.MySQL.DbUser, GlobalConfig.MySQL.DbPassword,
		GlobalConfig.MySQL.DbAddress, GlobalConfig.MySQL.DbName)
}

// 获取JWT签名方法
func GetJWTSigningMethod() jwt.SigningMethod {
	if GlobalConfig == nil || GlobalConfig.JWT.SigningMethod == "HS256" {
		return jwt.SigningMethodHS256
	}
	return jwt.SigningMethodHS256
}

// 获取JWT过期时间
func GetJWTExpirePeriod() time.Duration {
	if GlobalConfig == nil {
		return 30 * 24 * time.Hour
	}
	return time.Duration(GlobalConfig.JWT.TokenExpireHours) * time.Hour
}
