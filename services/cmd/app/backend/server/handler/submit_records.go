package handler

import (
	"FeasOJ/app/backend/internal/global"
	"FeasOJ/pkg/databases/repository"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// 获取所有提交记录
func GetAllSubmitRecords(c *gin.Context) {
	submitrecords := repository.SelectAllSubmitRecords(global.Db)
	c.JSON(http.StatusOK, gin.H{"submitrecords": submitrecords})
}

// 获取指定用户提交记录
func GetSubmitRecordsByUsername(c *gin.Context) {
	checker := c.GetHeader("Username")
	encodedUsername, _ := url.QueryUnescape(checker)
	username := c.Param("username")
	if encodedUsername != username {
		uid := repository.SelectUserInfo(global.Db, username).Id
		submitrecords := repository.SelectSRByUidForChecker(global.Db, uid)
		c.JSON(http.StatusOK, gin.H{"submitrecords": submitrecords})
		return
	} else {
		uid := repository.SelectUserInfo(global.Db, username).Id
		submitrecords := repository.SelectSubmitRecordsByUid(global.Db, uid)
		c.JSON(http.StatusOK, gin.H{"submitrecords": submitrecords})
		return
	}

}
