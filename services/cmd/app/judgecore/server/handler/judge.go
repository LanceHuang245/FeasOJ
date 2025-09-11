package handler

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	CodeDir string
}

func NewHandler(codeDir string) *Handler {
	return &Handler{CodeDir: codeDir}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (h *Handler) Judge(c *gin.Context) {
	file, err := c.FormFile("code")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get form file"})
		return
	}

	savePath := filepath.Join(h.CodeDir, file.Filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File received"})
}
