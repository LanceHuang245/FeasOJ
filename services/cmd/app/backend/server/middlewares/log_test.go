package middlewares

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoggerMiddleware(t *testing.T) {
	// 捕获日志输出
	var buf bytes.Buffer
	origOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(origOutput)

	// 创建一个新的 Gin 引擎，并使用 Logger 中间件
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Logger())

	router.GET("/test/path", func(c *gin.Context) {
		// 模拟业务处理
		c.String(http.StatusOK, "OK")
	})

	req, err := http.NewRequest("GET", "/test/path", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.RemoteAddr = "192.168.0.1:12345"

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Error(http.StatusOK, recorder.Code)
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, "GET") {
		t.Error(logOutput)
	}
	if !strings.Contains(logOutput, "/test/path") {
		t.Error(logOutput)
	}
	if !strings.Contains(logOutput, "192.168.0.1") {
		t.Error(logOutput)
	}
	if !strings.Contains(logOutput, "200") {
		t.Error(logOutput)
	}
}
