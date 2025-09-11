package judge

import (
	"FeasOJ/internal/config"
	"FeasOJ/internal/global"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/consul/api"
)

// JudgeCorePing 检查服务是否可用
func JudgeCorePing() bool {
	JudgeCoreAddress := GetJudgeCoreAddress()
	if JudgeCoreAddress == "" {
		log.Println("[FeasOJ] Unable to get JudgeCore address from Consul")
		return false
	}

	resp, err := http.Get(JudgeCoreAddress)
	if err != nil {
		log.Println("[FeasOJ] Error requesting JudgeCore:", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// 获取JudgeCore服务地址
func GetJudgeCoreAddress() string {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = config.GlobalConfig.Consul.Address
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Println("[FeasOJ] Error connecting to Consul:", err)
		return ""
	}

	services, _, err := consulClient.Catalog().Service("JudgeCore", "", nil)
	if err != nil {
		log.Println("[FeasOJ] Error querying Consul:", err)
		return ""
	}

	if len(services) == 0 {
		log.Println("[FeasOJ] JudgeCore service not found in Consul")
		return ""
	}

	service := services[0]
	global.JudgeCoreAddr = "http://" + service.ServiceAddress + ":" + fmt.Sprint(service.ServicePort)
	return "http://" + service.ServiceAddress + ":" + fmt.Sprint(service.ServicePort) + "/api/v1/judgecore/health"
}
