package utils

import (
	"FeasOJ/internal/config"
	"FeasOJ/internal/global"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/nfnt/resize"
)

// 压缩图像为256*256
func CompressImage(inputPath, outputPath string) error {
	// 打开图像文件
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}
	if format != "png" {
		return errors.New("unsupported image format")
	}

	newImage := resize.Resize(256, 256, img, resize.Lanczos3)

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	err = png.Encode(out, newImage)
	return err
}

// ImageGuardPing 检查服务是否可用
func ImageGuardPing() bool {
	imageGuardPingAddress := GetImageGuardAddress()
	if imageGuardPingAddress == "" {
		log.Println("[FeasOJ] Unable to get ImageGuard address from Consul")
		return false
	}
	resp, err := http.Get(imageGuardPingAddress)
	if err != nil {
		log.Println("[FeasOJ] Error requesting ImageGuard:", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// 获取ImageGuard服务
func GetImageGuardAddress() string {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = config.GlobalConfig.Consul.Address
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Println("[FeasOJ] Error connecting to Consul:", err)
		return ""
	}

	services, _, err := consulClient.Catalog().Service("ImageGuard", "", nil)
	if err != nil {
		log.Println("[FeasOJ] Error querying Consul:", err)
		return ""
	}

	if len(services) == 0 {
		log.Println("[FeasOJ] ImageGuard service not found in Consul")
		return ""
	}

	service := services[0]
	global.ImageGuardAddr = "http://" + service.ServiceAddress + ":" + fmt.Sprint(service.ServicePort)
	return "http://" + service.ServiceAddress + ":" + fmt.Sprint(service.ServicePort) + "/api/v1/image/predict"
}

// PredictImage 判断图片是否违规
func PredictImage(imagePath string) bool {
	// 打开图片文件
	file, err := os.Open(imagePath)
	if err != nil {
		return false
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", imagePath)
	if err != nil {
		return false
	}
	if _, err = io.Copy(part, file); err != nil {
		return false
	}
	if err = writer.Close(); err != nil {
		return false
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/image/predict", global.ImageGuardAddr), &buf)
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Prediction string `json:"prediction"`
		Error      string `json:"error,omitempty"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false
	}

	// 返回 true 表示预测结果为 "neutral"
	return result.Prediction == "neutral"
}
