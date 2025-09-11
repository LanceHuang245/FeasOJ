package handler

import (
	"FeasOJ/internal/utils/sql"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// 获取所有提交记录
func GetAllSubmitRecords(c *gin.Context) {
	submitrecords := sql.SelectAllSubmitRecords()
	c.JSON(http.StatusOK, gin.H{"submitrecords": submitrecords})
}

// 获取指定用户提交记录
func GetSubmitRecordsByUsername(c *gin.Context) {
	checker := c.GetHeader("Username")
	encodedUsername, _ := url.QueryUnescape(checker)
	username := c.Param("username")
	if encodedUsername != username {
		uid := sql.SelectUserInfo(username).Uid
		submitrecords := sql.SelectSRByUidForChecker(uid)
		c.JSON(http.StatusOK, gin.H{"submitrecords": submitrecords})
		return
	} else {
		uid := sql.SelectUserInfo(username).Uid
		submitrecords := sql.SelectSubmitRecordsByUid(uid)
		c.JSON(http.StatusOK, gin.H{"submitrecords": submitrecords})
		return
	}

}
