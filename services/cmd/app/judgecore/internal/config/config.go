package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Consul struct {
	Address     string `json:"address"`
	ServiceName string `json:"service_name"`
	ServiceID   string `json:"service_id"`
}

type RabbitMQ struct {
	Address string `json:"address"`
}

type Server struct {
	Address     string `json:"address"`
	Port        int    `json:"port"`
	EnableHTTPS bool   `json:"enable_https"`
	CertPath    string `json:"cert_path"`
	KeyPath     string `json:"key_path"`
}

type Sandbox struct {
	Memory        int64   `json:"memory"`         // 内存限制 (字节)
	NanoCPUs      float64 `json:"nano_cpus"`      // CPU限制 (核心数)
	CPUShares     int64   `json:"cpu_shares"`     // CPU权重
	MaxConcurrent int     `json:"max_concurrent"` // 最大并发数
}

type Database struct {
	Address      string `json:"address"`
	Name         string `json:"name"`
	User         string `json:"user"`
	Password     string `json:"password"`
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
	MaxLifeTime  int    `json:"max_life_time"`
}

// AppConfig 配置结构体
type AppConfig struct {
	Consul   Consul   `json:"consul"`
	RabbitMQ RabbitMQ `json:"rabbitmq"`
	Server   Server   `json:"server"`
	Sandbox  Sandbox  `json:"sandbox"`
	Database Database `json:"database"`
}

// LoadConfig 加载JSON配置文件
func LoadConfig(currentDir string) (*AppConfig, error) {
	filePath := filepath.Join(currentDir, "config.json")

	// 检查配置文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 如果配置文件不存在，创建默认配置
		if err := createDefaultConfig(filePath); err != nil {
			return nil, fmt.Errorf("failed to create default config file: %v", err)
		}
		log.Println("[FeasOJ] Created default config file:", filePath)
	}

	// 读取配置文件
	configFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer configFile.Close()

	var config AppConfig
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// createDefaultConfig 创建默认配置文件
func createDefaultConfig(filePath string) error {
	defaultConfig := AppConfig{
		Consul: struct {
			Address     string `json:"address"`
			ServiceName string `json:"service_name"`
			ServiceID   string `json:"service_id"`
		}{
			Address:     "127.0.0.1:8500",
			ServiceName: "JudgeCore",
			ServiceID:   "JudgeCore-1",
		},
		RabbitMQ: struct {
			Address string `json:"address"`
		}{
			Address: "amqp://USER:PASSWORD@IP:PORT/",
		},
		Server: struct {
			Address     string `json:"address"`
			Port        int    `json:"port"`
			EnableHTTPS bool   `json:"enable_https"`
			CertPath    string `json:"cert_path"`
			KeyPath     string `json:"key_path"`
		}{
			Address:     "127.0.0.1",
			Port:        37885,
			EnableHTTPS: false,
			CertPath:    "./certificate/fullchain.pem",
			KeyPath:     "./certificate/privkey.key",
		},
		Sandbox: struct {
			Memory        int64   `json:"memory"`
			NanoCPUs      float64 `json:"nano_cpus"`
			CPUShares     int64   `json:"cpu_shares"`
			MaxConcurrent int     `json:"max_concurrent"`
		}{
			Memory:        2 * 1024 * 1024 * 1024,
			NanoCPUs:      0.5,
			CPUShares:     1024,
			MaxConcurrent: 5,
		},
		Database: struct {
			Address      string `json:"address"`
			Name         string `json:"name"`
			User         string `json:"user"`
			Password     string `json:"password"`
			MaxOpenConns int    `json:"max_open_conns"`
			MaxIdleConns int    `json:"max_idle_conns"`
			MaxLifeTime  int    `json:"max_life_time"`
		}{
			Address:      "IP:PORT",
			Name:         "DATABASE_NAME",
			User:         "DATABASE_USER",
			Password:     "DATABASE_PASSWORD",
			MaxOpenConns: 240,
			MaxIdleConns: 100,
			MaxLifeTime:  32,
		},
	}

	// 将配置写入文件
	configData, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, configData, 0644)
}

// SaveConfig 保存配置到文件
func SaveConfig(config *AppConfig, currentDir string) error {
	filePath := filepath.Join(currentDir, "config.json")
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, configData, 0644)
}
