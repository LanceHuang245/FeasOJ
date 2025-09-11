package middlewares

import (
	"FeasOJ/app/backend/internal/global"
	"FeasOJ/app/backend/internal/utils"
	"FeasOJ/app/backend/server/handler"
	"FeasOJ/pkg/databases/repository"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func HeaderVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		var User string
		encodedUsername := c.GetHeader("Username")
		username, err := url.QueryUnescape(encodedUsername)
		token := c.GetHeader("Authorization")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": handler.GetMessage(c, "userNotFound")})
			c.Abort()
			return
		}
		if utils.IsEmail(username) {
			User = repository.SelectUserByEmail(global.Db, username).Username
		} else {
			User = username
		}
		if !utils.VerifyToken(User, token) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": handler.GetMessage(c, "unauthorized")})
			c.Abort()
			return
		}
		c.Next()
	}
}
