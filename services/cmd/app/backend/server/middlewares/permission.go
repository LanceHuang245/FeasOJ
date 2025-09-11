package middlewares

import (
	"FeasOJ/app/backend/internal/global"
	"FeasOJ/app/backend/server/handler"
	"FeasOJ/pkg/databases/repository"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func PermissionChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		encodedUsername := c.GetHeader("Username")
		username, _ := url.QueryUnescape(encodedUsername)
		if repository.SelectUserInfo(global.Db, username).Role != 1 {
			c.JSON(http.StatusForbidden, gin.H{"message": handler.GetMessage(c, "forbidden")})
			c.Abort()
			return
		}
		c.Next()
	}
}
