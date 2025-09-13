package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Consul struct {
	Host        string `toml:"host"`
	ServiceName string `toml:"service_name"`
	ServiceID   string `toml:"service_id"`
}

type RabbitMQ struct {
	Host string `toml:"host"`
}

type Server struct {
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	EnableHTTPS bool   `toml:"enable_https"`
	CertPath    string `toml:"cert_path"`
	KeyPath     string `toml:"key_path"`
}

type Sandbox struct {
	Memory        int64   `toml:"memory"`         // 内存限制 (字节)
	NanoCPUs      float64 `toml:"nano_cpus"`      // CPU限制 (核心数)
	CPUShares     int64   `toml:"cpu_shares"`     // CPU权重
	MaxConcurrent int     `toml:"max_concurrent"` // 最大并发数
}

type Database struct {
	Type     string `toml:"type"` // 数据库类型: mysql, postgresql
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Name     string `toml:"name"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	SSLMode  string `toml:"ssl_mode"` // PostgreSQL SSL模式
}

// AppConfig 配置结构体
type AppConfig struct {
	Consul   Consul   `toml:"consul"`
	RabbitMQ RabbitMQ `toml:"rabbitmq"`
	Server   Server   `toml:"server"`
	Sandbox  Sandbox  `toml:"sandbox"`
	Database Database `toml:"database"`
}

// LoadConfig 加载TOML配置文件
func LoadConfig(currentDir string) (*AppConfig, error) {
	configPath := filepath.Join(currentDir, "config.toml")

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 如果配置文件不存在，创建默认配置
		if err := createDefaultConfig(configPath); err != nil {
			return nil, fmt.Errorf("failed to create default config file: %v", err)
		}
		log.Println("[FeasOJ] Created default config file:", configPath)
	}

	// 读取配置文件
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config AppConfig
	if _, err := toml.Decode(string(configData), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// createDefaultConfig 创建默认配置文件
func createDefaultConfig(configPath string) error {
	defaultConfig := AppConfig{
		Consul: Consul{
			Host:        "127.0.0.1:8500",
			ServiceName: "JudgeCore",
			ServiceID:   "JudgeCore-1",
		},
		RabbitMQ: RabbitMQ{
			Host: "amqp://rabbitmq:password@localhost:5672/",
		},
		Server: Server{
			Host:        "127.0.0.1",
			Port:        37885,
			EnableHTTPS: false,
			CertPath:    "./certificate/fullchain.pem",
			KeyPath:     "./certificate/privkey.key",
		},
		Sandbox: Sandbox{
			Memory:        2 * 1024 * 1024 * 1024,
			NanoCPUs:      0.5,
			CPUShares:     1024,
			MaxConcurrent: 5,
		},
		Database: Database{
			Type:     "mysql",
			Host:     "localhost",
			Port:     3306,
			Name:     "feasoj",
			User:     "feasoj",
			Password: "password",
			SSLMode:  "disable",
		},
	}

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(defaultConfig); err != nil {
		return err
	}

	return os.WriteFile(configPath, buf.Bytes(), 0644)
}

// SaveConfig 保存配置到文件
func SaveConfig(config *AppConfig, currentDir string) error {
	configPath := filepath.Join(currentDir, "config.toml")

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		return err
	}

	return os.WriteFile(configPath, buf.Bytes(), 0644)
}
