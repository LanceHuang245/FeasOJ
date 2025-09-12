package handler

import (
	"FeasOJ/app/backend/internal/global"
	"FeasOJ/pkg/databases/repository"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 用户获取竞赛列表
func GetCompetitionsList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": repository.SelectCompetitionsInfo(global.Db)})
}

// 用户获取指定竞赛ID信息
func GetCompetitionInfoByID(c *gin.Context) {
	encodedUsername := c.GetHeader("Username")
	username, _ := url.QueryUnescape(encodedUsername)
	competitionId, _ := strconv.Atoi(c.Param("id"))
	uid := repository.SelectUserInfo(global.Db, username).Id
	if repository.SelectUserCompetition(global.Db, uid, competitionId) {
		c.JSON(http.StatusOK, gin.H{"data": repository.SelectCompetitionInfoByCid(global.Db, competitionId)})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
	}
}

// 用户加入有密码的竞赛
func JoinCompetitionWithPassword(c *gin.Context) {
	encodedUsername := c.GetHeader("Username")
	username, _ := url.QueryUnescape(encodedUsername)
	competitionId, _ := strconv.Atoi(c.Param("id"))
	competitionPwd := c.Query("password")
	uid := repository.SelectUserInfo(global.Db, username).Id
	if repository.SelectUserCompetition(global.Db, uid, competitionId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "userAlreadyExists")})
		return
	}
	if repository.SelectCompetitionInfoAdminByCid(global.Db, competitionId).Password == competitionPwd {
		if repository.AddUserCompetition(global.Db, uid, competitionId) == nil {
			c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
}

// 用户加入竞赛
func JoinCompetition(c *gin.Context) {
	encodedUsername := c.GetHeader("Username")
	username, _ := url.QueryUnescape(encodedUsername)
	competitionId, _ := strconv.Atoi(c.Param("id"))
	uid := repository.SelectUserInfo(global.Db, username).Id
	if repository.SelectUserCompetition(global.Db, uid, competitionId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "userAlreadyExists")})
		return
	}
	if repository.AddUserCompetition(global.Db, uid, competitionId) == nil {
		c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
}

// 查询用户是否在竞赛中
func IsInCompetition(c *gin.Context) {
	encodedUsername := c.GetHeader("Username")
	username, _ := url.QueryUnescape(encodedUsername)
	competitionId, _ := strconv.Atoi(c.Param("id"))
	uid := repository.SelectUserInfo(global.Db, username).Id

	in := repository.SelectUserCompetition(global.Db, uid, competitionId)
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success"), "isIn": in})
}

// 查询指定竞赛中的所有参与用户
func GetCompetitionUsers(c *gin.Context) {
	competitionId, _ := strconv.Atoi(c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"data": repository.SelectUsersCompetition(global.Db, competitionId)})
}

// 用户退出竞赛
func QuitCompetition(c *gin.Context) {
	encodedUsername := c.GetHeader("Username")
	username, _ := url.QueryUnescape(encodedUsername)
	competitionId, _ := strconv.Atoi(c.Param("id"))
	uid := repository.SelectUserInfo(global.Db, username).Id
	if repository.DeleteUserCompetition(global.Db, uid, competitionId) == nil {
		c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
}

// 获取指定竞赛的所有题目
func GetProblemsByCompetitionID(c *gin.Context) {
	competitionId, _ := strconv.Atoi(c.Param("id"))
	if repository.SelectCompetitionInfoByCid(global.Db, competitionId).Status == 0 {
		c.JSON(http.StatusForbidden, gin.H{"message": GetMessage(c, "forbidden")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": repository.SelectProblemsByCompID(global.Db, competitionId)})
}
