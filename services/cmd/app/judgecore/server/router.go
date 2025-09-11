package server

import (
	"JudgeCore/internal/judge"
	"JudgeCore/server/handler"
	"JudgeCore/server/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoadRouter(r *gin.Engine, db *gorm.DB, pool *judge.JudgePool, codeDir string) {
	r.Use(middlewares.Logger())

	// Create a handler instance with its dependencies
	h := handler.NewHandler(codeDir)

	apiV1 := r.Group("/api/v1/judgecore")
	{
		apiV1.GET("/health", h.Health)
		apiV1.POST("/judge", h.Judge)
	}
}