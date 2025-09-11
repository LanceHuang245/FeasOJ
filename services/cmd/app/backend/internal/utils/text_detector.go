package utils

import (
	"FeasOJ/internal/config"
	"FeasOJ/internal/global"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"fmt"

	"github.com/hashicorp/consul/api"
)

// PredictionResponse 定义检测服务返回的 JSON 结构
type PredictionResponse struct {
	Text       string `json:"text"`
	Prediction string `json:"prediction"`
}

// 获取Profanity Detector服务地址
func GetProfanityDetectorAddress() string {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = config.GlobalConfig.Consul.Address
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Println("[FeasOJ] Error connecting to Consul:", err)
		return ""
	}

	services, _, err := consulClient.Catalog().Service("ProfanityDetector", "", nil)
	if err != nil {
		log.Println("[FeasOJ] Error querying Consul:", err)
		return ""
	}

	if len(services) == 0 {
		log.Println("[FeasOJ] ProfanityDetector service not found in Consul")
		return ""
	}

	service := services[0]
	global.ProfanityDetectorAddr = "http://" + service.ServiceAddress + ":" + fmt.Sprint(service.ServicePort)
	return "http://" + service.ServiceAddress + ":" + fmt.Sprint(service.ServicePort) + "/api/v1/text/predict"
}

// DetectText 检测文本是否包含敏感词汇
func DetectText(text string) bool {
	// 获取服务地址
	serviceAddress := GetProfanityDetectorAddress()
	if serviceAddress == "" {
		log.Println("[FeasOJ] Unable to get ProfanityDetector address from Consul")
		return false
	}

	// 准备请求数据
	requestData := map[string]string{
		"text": text,
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Println("[FeasOJ] Error marshaling request data:", err)
		return false
	}

	// 发送请求
	resp, err := http.Post(serviceAddress, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("[FeasOJ] Error requesting ProfanityDetector:", err)
		return false
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[FeasOJ] Error reading response body:", err)
		return false
	}

	// 解析响应
	var result PredictionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("[FeasOJ] Error unmarshaling response:", err)
		return false
	}

	// 返回 true 表示预测结果为 "neutral"
	return result.Prediction == "neutral"
}

// ProfanityDetectorPing 检查服务是否可用
func ProfanityDetectorPing() bool {
	serviceAddress := GetProfanityDetectorAddress()
	if serviceAddress == "" {
		log.Println("[FeasOJ] Unable to get ProfanityDetector address from Consul")
		return false
	}

	// 发送ping请求
	resp, err := http.Get(global.ProfanityDetectorAddr + "/api/v1/text/predict")
	if err != nil {
		log.Println("[FeasOJ] Error requesting ProfanityDetector:", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
