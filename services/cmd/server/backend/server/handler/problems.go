package handler

import (
	"FeasOJ/internal/global"
	"FeasOJ/internal/utils"
	"FeasOJ/internal/utils/sql"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取所有题目
func GetAllProblems(c *gin.Context) {
	// 实时性要求较高，不做数据缓存
	problems := sql.SelectAllProblems()
	c.JSON(http.StatusOK, gin.H{"problems": problems})
}

// 获取题目信息
func GetProblemInfo(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Param("id"))
	if !sql.IsProblemVisible(pid) {
		c.JSON(http.StatusForbidden, gin.H{"message": GetMessage(c, "forbidden")})
		return
	}
	// 生成缓存键
	cacheKey := "problemInfo_" + c.Param("id")
	var problemInfo global.ProblemInfoRequest

	// 从缓存中获取数据
	err := utils.GetCache(cacheKey, &problemInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}
	if problemInfo.Pid == 0 {
		// 缓存未命中
		problemInfo = sql.SelectProblemInfo(c.Param("id"))
		// 数据存入缓存，时间10分钟
		err = utils.SetCache(cacheKey, problemInfo, 10*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
			return
		}

	}
	c.JSON(http.StatusOK, gin.H{"problemInfo": problemInfo})
}

// 提交代码
func UploadCode(c *gin.Context) {
	problem := c.Param("pid")
	pidInt, _ := strconv.Atoi(problem)
	encodedUsername := c.GetHeader("Username")
	username, _ := url.QueryUnescape(encodedUsername)
	file, err := c.FormFile("code")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "invalidrequest")})
		return
	}

	// 获取用户ID
	userInfo := sql.SelectUserInfo(username)

	// 连接到 RabbitMQ
	conn, ch, err := utils.ConnectRabbitMQ()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}
	defer conn.Close()
	defer ch.Close()

	rdb := utils.ConnectRedis()
	defer rdb.Close()

	// 提交频率限制
	userRateLimitKey := fmt.Sprintf("ratelimit:%d", userInfo.Uid)

	// 检查是否在限制时间内
	exists, _ := rdb.Exists(userRateLimitKey).Result()
	if exists == 1 {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": GetMessage(c, "rateLimit")})
		return
	}

	// 设置限流键，10秒后自动失效
	rdb.Set(userRateLimitKey, 1, 10*time.Second)

	// 将文件名改为用户ID_题目ID.扩展名
	newFileName := fmt.Sprintf("%d_%s%s", userInfo.Uid, problem, path.Ext(file.Filename))

	// 保存文件到临时目录
	tempFilePath := filepath.Join(os.TempDir(), newFileName)
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save file"})
		return
	}

	// 转发文件至 JudgeCore
	if err := ForwardFile(tempFilePath, newFileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to forward file"})
		return
	}

	// 读取文件内容
	code, err := os.ReadFile(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read file"})
		return
	}

	// 移除临时文件
	os.Remove(tempFilePath)

	var language string
	switch path.Ext(file.Filename) {
	case ".cpp":
		language = "C++"
	case ".java":
		language = "Java"
	case ".py":
		language = "Python"
	case ".go":
		language = "Go"
	case ".rs":
		language = "Rust"
	case ".php":
		language = "PHP"
	case ".pas":
		language = "Pascal"
	default:
		language = "Unknown"
	}

	// 将任务发送到RabbitMQ
	err = utils.PublishTask(ch, newFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}
	sql.AddSubmitRecord(userInfo.Uid, pidInt, "Running...", language, username, string(code))
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 转发文件至 JudgeCore
func ForwardFile(filePath, fileName string) error {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建一个 buffer 来存储 multipart 数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 创建 form 文件字段
	part, err := writer.CreateFormFile("code", fileName)
	if err != nil {
		return err
	}

	// 复制文件数据到 part
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	// 关闭 writer，确保 boundary 信息正确写入
	writer.Close()

	// 发送 HTTP 请求
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/judgecore/judge", global.JudgeCoreAddr), body)
	if err != nil {
		return err
	}

	// 设置 Content-Type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to forward file, status code: %d", resp.StatusCode)
	}

	return nil
}
