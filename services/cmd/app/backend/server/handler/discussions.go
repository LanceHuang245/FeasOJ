package handler

import (
	"FeasOJ/app/backend/internal/config"
	"FeasOJ/app/backend/internal/global"
	"FeasOJ/app/backend/internal/utils"
	"FeasOJ/pkg/databases/repository"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取所有讨论列表
func GetAllDiscussions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	itemsPerPage, _ := strconv.Atoi(c.DefaultQuery("itemsPerPage", "12"))

	discussions, total := repository.SelectDiscussList(global.Db, page, itemsPerPage)
	c.JSON(http.StatusOK, gin.H{
		"data":  discussions,
		"total": total,
	})
}

// 获取指定id讨论信息
func GetDiscussionByDid(c *gin.Context) {
	did, _ := strconv.Atoi(c.Param("id"))
	discussion := repository.SelectDiscussionByDid(global.Db, did)
	c.JSON(http.StatusOK, gin.H{"data": discussion})
}

// 创建讨论
func CreateDiscussion(c *gin.Context) {
	encodedUsername := c.GetHeader("Username")
	username, _ := url.QueryUnescape(encodedUsername)
	title := c.PostForm("title")
	content := c.PostForm("content")

	// 检测文本是否包含敏感词汇
	if config.GlobalConfig.Features.ProfanityDetectorEnabled {
		if utils.DetectText(title) {
			c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "profanity")})
			return
		}
		if utils.DetectText(content) {
			c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "profanity")})
			return
		}
	}

	// 获取用户ID
	userInfo := repository.SelectUserInfo(global.Db, username)

	rdb := utils.ConnectRedis()
	defer rdb.Close()

	// 设置频率限制键
	userRateLimitKey := fmt.Sprintf("discussionRateLimit:%d", userInfo.Id)
	exists, _ := rdb.Exists(userRateLimitKey).Result()
	if exists == 1 {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": GetMessage(c, "rateLimit")})
		return
	}

	// 设置限流键
	rdb.Set(userRateLimitKey, 1, 15*time.Second)

	// 创建讨论
	if !repository.AddDiscussion(global.Db, title, content, userInfo.Id) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 删除讨论
func DeleteDiscussion(c *gin.Context) {
	did, _ := strconv.Atoi(c.Param("id"))
	if repository.DelDiscussion(global.Db, did) {
		c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
	}
}

// 获取指定讨论的评论
func GetComment(c *gin.Context) {
	did, _ := strconv.Atoi(c.Param("id"))
	comments := repository.SelectCommentsByDid(global.Db, did)
	c.JSON(http.StatusOK, gin.H{"data": comments})
}

// 删除指定Cid的评论
func DelComment(c *gin.Context) {
	cid, _ := strconv.Atoi(c.Param("id"))
	if !repository.DeleteCommentByCid(global.Db, cid) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 添加评论
func AddComment(c *gin.Context) {
	encodedUsername := c.GetHeader("Username")
	username, _ := url.QueryUnescape(encodedUsername)
	content := c.PostForm("content")
	did, _ := strconv.Atoi(c.Param("id"))
	// 获取用户ID
	userInfo := repository.SelectUserInfo(global.Db, username)

	rdb := utils.ConnectRedis()
	defer rdb.Close()

	// 设置频率限制键
	userRateLimitKey := fmt.Sprintf("commentRateLimit:%d", userInfo.Id)
	exists, _ := rdb.Exists(userRateLimitKey).Result()
	if exists == 1 {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": GetMessage(c, "rateLimit")})
		return
	}

	// 设置限流键
	rdb.Set(userRateLimitKey, 1, 10*time.Second)

	// 检测文本是否包含敏感词汇
	if config.GlobalConfig.Features.ProfanityDetectorEnabled {
		if utils.DetectText(content) {
			c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "profanity")})
			return
		}
	}

	if !repository.AddComment(global.Db, content, did, userInfo.Id, false) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}
