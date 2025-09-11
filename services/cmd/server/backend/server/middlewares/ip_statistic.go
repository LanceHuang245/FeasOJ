package middlewares

import (
	"FeasOJ/internal/global"
	"time"

	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IP 访问统计
func IPStatistic() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		visit := global.IPVisit{
			IP:         ip,
			VisitCount: 1,
			LastVisit:  now,
		}

		global.DB.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "ip"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"visit_count": gorm.Expr("visit_count + ?", 1),
				"last_visit":  now,
			}),
		}).Create(&visit)

		c.Next()
	}
}
