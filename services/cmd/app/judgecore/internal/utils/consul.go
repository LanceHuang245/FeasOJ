package utils

import (
	"JudgeCore/internal/config"
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

// RegService 服务注册
func RegService(client *api.Client, consulConfig config.Consul, serverConfig config.Server) error {
	agent := client.Agent()

	protocol := "http"
	if serverConfig.EnableHTTPS {
		protocol = "https"
	}

	registration := &api.AgentServiceRegistration{
		ID:   consulConfig.ServiceID,   // 服务唯一ID
		Name: consulConfig.ServiceName, // 服务名称
		Port: serverConfig.Port,        // 服务端口
		Tags: []string{"gin", "judge"},
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("%s://%s:%d/api/v1/judgecore/health", protocol, serverConfig.Address, serverConfig.Port), // 健康检查地址
			Interval: "60s",
			Timeout:  "6s",
		},
	}

	err := agent.ServiceRegister(registration)
	if err != nil {
		log.Println("[FeasOJ] JudgeCore service registration failed:", err)
		return err
	}
	log.Println("[FeasOJ] JudgeCore service registered successfully")
	return nil
}
