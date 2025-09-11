package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// InitializeLogger 初始化日志模块
func InitializeLogger(logDir string) (*os.File, error) {
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return nil, err
	}

	currentTime := time.Now().Format("2006-01-02_15-04-05")
	logFileName := filepath.Join(logDir, currentTime+".txt")
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	return logFile, nil
}

// CloseLogger 关闭日志模块
func CloseLogger(logFile *os.File) {
	if logFile != nil {
		logFile.Close()
	}
}
