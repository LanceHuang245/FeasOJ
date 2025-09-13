package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/golang-jwt/jwt/v5"
)

// Config 全局配置结构体
type Config struct {
	Server   ServerConfig   `toml:"server"`
	RabbitMQ RabbitMQConfig `toml:"rabbitmq"`
	Consul   ConsulConfig   `toml:"consul"`
	Features FeaturesConfig `toml:"features"`
	Database DatabaseConfig `toml:"database"`
	Redis    RedisConfig    `toml:"redis"`
	Mail     MailConfig     `toml:"mail"`
	JWT      JWTConfig      `toml:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host        string `toml:"host"`
	EnableHTTPS bool   `toml:"enable_https"`
	CertPath    string `toml:"cert_path"`
	KeyPath     string `toml:"key_path"`
}

// RabbitMQConfig RabbitMQ配置
type RabbitMQConfig struct {
	Host string `toml:"host"`
}

// ConsulConfig Consul配置
type ConsulConfig struct {
	Host string `toml:"host"`
}

// FeaturesConfig 功能开关配置
type FeaturesConfig struct {
	ImageGuardEnabled        bool `toml:"image_guard_enabled"`
	ProfanityDetectorEnabled bool `toml:"profanity_detector_enabled"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type         string `toml:"type"`          // 数据库类型: mysql, postgresql
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	Name         string `toml:"name"`
	User         string `toml:"user"`
	Password     string `toml:"password"`
	SSLMode      string `toml:"ssl_mode"`      // PostgreSQL SSL模式
	MaxOpenConns int    `toml:"max_open_conns"`
	MaxIdleConns int    `toml:"max_idle_conns"`
	MaxLifeTime  int    `toml:"max_life_time"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `toml:"host"`
	Password string `toml:"password"`
}

// MailConfig 邮件配置
type MailConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	SigningMethod    string `toml:"signing_method"`
	TokenExpireHours int    `toml:"token_expire_hours"`
	SecretKey        string `toml:"secret_key"`
}

// 全局配置实例
var GlobalConfig *Config

// 初始化配置
func InitConfig() error {
	configPath := "config.toml"

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Println("[FeasOJ] Configuration file does not exist, creating default configuration file...")
		if err := createDefaultConfig(configPath); err != nil {
			return fmt.Errorf("failed to create default configuration file")
		}
		log.Println("[FeasOJ] Default configuration file created, please edit config.toml file and restart the program")
		return fmt.Errorf("error connection to the database and other services")
	}

	// 读取配置文件
	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Panicln("[FeasOJ] Failed to read configuration file, please check file permissions")
		return fmt.Errorf("failed to read configuration")
	}

	// 解析TOML配置
	GlobalConfig = &Config{}
	if _, err := toml.Decode(string(configData), GlobalConfig); err != nil {
		log.Panicln("[FeasOJ] Failed to parse configuration file, please check the format of the configuration file")
		return fmt.Errorf("failed to read configuration")
	}

	// 验证配置
	if err := validateConfig(GlobalConfig); err != nil {
		log.Panicln("[FeasOJ] Configuration validation failed, please check the configuration file")
		return fmt.Errorf("failed to read configuration")
	}

	log.Println("[FeasOJ] Configuration file loaded successfully")
	return nil
}

// 创建默认配置文件
func createDefaultConfig(configPath string) error {
	defaultConfig := &Config{
		Server: ServerConfig{
			Host:        "127.0.0.1:37882",
			EnableHTTPS: true,
			CertPath:    "./certificate/fullchain.pem",
			KeyPath:     "./certificate/privkey.key",
		},
		RabbitMQ: RabbitMQConfig{
			Host: "amqp://rabbitmq:password@localhost:5672/",
		},
		Consul: ConsulConfig{
			Host: "localhost:8500",
		},
		Features: FeaturesConfig{
			ImageGuardEnabled:        false,
			ProfanityDetectorEnabled: false,
		},
		Database: DatabaseConfig{
			Type:         "mysql",
			Host:         "localhost",
			Port:         3306,
			Name:         "feasoj",
			User:         "feasoj",
			Password:     "password",
			SSLMode:      "disable",
			MaxOpenConns: 240,
			MaxIdleConns: 100,
			MaxLifeTime:  32,
		},
		Redis: RedisConfig{
			Host:     "localhost:6379",
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
			SecretKey:        "default-secret-key",
		},
	}

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(defaultConfig); err != nil {
		return err
	}

	return os.WriteFile(configPath, buf.Bytes(), 0644)
}

// 验证配置
func validateConfig(config *Config) error {
	if config.Server.Host == "" {
		return fmt.Errorf("Server Address cannot be empty")
	}
	if config.Database.Host == "" {
		return fmt.Errorf("Database configuration is incomplete")
	}
	if config.Redis.Host == "" {
		return fmt.Errorf("Redis Address cannot be empty")
	}
	return nil
}

// 获取数据库连接字符串
func GetDatabaseDSN() string {
	if GlobalConfig == nil {
		return ""
	}
	
	// 根据数据库类型生成相应的连接字符串
	switch GlobalConfig.Database.Type {
	case "postgresql", "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			GlobalConfig.Database.Host, GlobalConfig.Database.Port,
			GlobalConfig.Database.User, GlobalConfig.Database.Password,
			GlobalConfig.Database.Name, GlobalConfig.Database.SSLMode)
	default: // 默认为MySQL
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai",
			GlobalConfig.Database.User, GlobalConfig.Database.Password,
			GlobalConfig.Database.Host, GlobalConfig.Database.Port,
			GlobalConfig.Database.Name)
	}
}

// 获取JWT签名方法
func GetJWTSigningMethod() jwt.SigningMethod {
	if GlobalConfig == nil || GlobalConfig.JWT.SigningMethod == "HS256" {
		return jwt.SigningMethodHS256
	}
	return jwt.SigningMethodHS256
}

// 获取JWT签名密钥
func GetJWTSecretKey() []byte {
	if GlobalConfig == nil {
		return []byte("default-secret-key")
	}
	return []byte(GlobalConfig.JWT.SecretKey)
}

// 获取JWT过期时间
func GetJWTExpirePeriod() time.Duration {
	if GlobalConfig == nil {
		return 30 * 24 * time.Hour
	}
	return time.Duration(GlobalConfig.JWT.TokenExpireHours) * time.Hour
}