package handler

import (
	"FeasOJ/app/backend/internal/global"
	"FeasOJ/pkg/databases/repository"
	"FeasOJ/pkg/databases/tables"
	"FeasOJ/pkg/structs"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// 管理员获取所有题目
func GetAllProblemsAdmin(c *gin.Context) {
	problems := repository.SelectAllProblemsAdmin(global.Db)
	c.JSON(http.StatusOK, gin.H{"problems": problems})
}

// 管理员获取指定题目所有信息
func GetProblemAllInfo(c *gin.Context) {
	problemInfo := repository.SelectProblemTestCases(global.Db, c.Param("pid"))
	c.JSON(http.StatusOK, gin.H{"problemInfo": problemInfo})
}

// 更新题目信息
func UpdateProblemInfo(c *gin.Context) {
	var req structs.AdminProblemInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "invalidrequest")})
		return
	}

	// 更新题目信息
	if err := repository.UpdateProblem(global.Db, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 删除题目及其输入输出样例
func DeleteProblem(c *gin.Context) {
	pid := c.Param("pid")
	pidInt, _ := strconv.Atoi(pid)
	if !repository.DeleteProblemAllInfo(global.Db, pidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 管理员获取所有用户信息
func GetAllUsersInfo(c *gin.Context) {
	usersInfo := repository.SelectAllUsersInfo(global.Db)
	c.JSON(http.StatusOK, gin.H{"usersInfo": usersInfo})
}

// 晋升用户
func PromoteUser(c *gin.Context) {
	uid := c.Query("uid")
	uidInt, _ := strconv.Atoi(uid)

	if !repository.PromoteToAdmin(global.Db, uidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 降级用户
func DemoteUser(c *gin.Context) {
	uid := c.Query("uid")
	uidInt, _ := strconv.Atoi(uid)

	if repository.SelectAdminCount(global.Db) <= 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}

	if !repository.DemoteToUser(global.Db, uidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 封禁用户
func BanUser(c *gin.Context) {
	uid := c.Query("uid")
	uidInt, _ := strconv.Atoi(uid)

	if !repository.BanUser(global.Db, uidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 解封用户
func UnbanUser(c *gin.Context) {
	uid := c.Query("uid")
	uidInt, _ := strconv.Atoi(uid)

	if !repository.UnbanUser(global.Db, uidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 管理员获取竞赛列表
func GetCompetitionListAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"contests": repository.SelectCompetitionInfoAdmin(global.Db)})
}

// 管理员获取指定竞赛ID信息
func GetCompetitionInfoAdmin(c *gin.Context) {
	cid := c.Param("cid")
	cidInt, _ := strconv.Atoi(cid)
	c.JSON(http.StatusOK, gin.H{"contest": repository.SelectCompetitionInfoAdminByCid(global.Db, cidInt)})
}

// 删除指定ID竞赛
func DeleteCompetition(c *gin.Context) {
	cid := c.Param("cid")
	cidInt, _ := strconv.Atoi(cid)

	if !repository.DeleteCompetition(global.Db, cidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 更新/添加竞赛信息
func UpdateCompetitionInfo(c *gin.Context) {
	var req structs.AdminCompetitionInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "invalidrequest")})
		return
	}

	// 更新竞赛信息
	if err := repository.UpdateCompetition(global.Db, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 计算分数
func CalculateScore(c *gin.Context) {
	competitionId, _ := strconv.Atoi(c.Param("cid"))

	// 查询竞赛信息
	var competition tables.Competition
	global.Db.First(&competition, competitionId)
	if competition.Scored {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "competition_scored")})
		return
	}

	// 查询竞赛参与用户
	users := repository.SelectUsersCompetition(global.Db, competitionId)
	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
		return
	}

	// 遍历所有参与竞赛的用户
	for _, user := range users {
		var submissions []tables.SubmitRecord
		global.Db.
			Where("uid = ? AND result = ? AND time BETWEEN ? AND ?",
				user.UserId,
				"Success",
				competition.StartAt,
				competition.EndAt).
			Find(&submissions)

		// 计算分数
		score := 0
		for _, submission := range submissions {
			var difficulty string
			err := global.Db.
				Table("problems").
				Select("difficulty").
				Where("contest_id = ? AND pid = ?", competitionId, submission.ProblemId).
				Row().
				Scan(&difficulty)
			if err != nil {
				continue
			}

			switch difficulty {
			case "简单":
				score += 1
			case "中等":
				score += 3
			case "困难":
				score += 5
			}
		}

		// 更新用户分数
		if score > 0 {
			global.Db.Model(&tables.User{}).Where("uid = ?", user.UserId).Update("score", gorm.Expr("score + ?", score))
		}

		global.Db.Model(&tables.UserCompetitions{}).Where("uid = ?", user.UserId).Update("score", score)
	}

	global.Db.Model(&tables.Competition{}).Where("contest_id = ?", competitionId).Update("scored", true)

	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 查询指定竞赛中，参与人员的得分情况
func GetScoreBoard(c *gin.Context) {
	competitionId, _ := strconv.Atoi(c.Param("cid"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	itemsPerPage, _ := strconv.Atoi(c.DefaultQuery("itemsPerPage", "10"))

	users, total := repository.GetScores(global.Db, competitionId, page, itemsPerPage)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
	})
}

// 获取IP访问统计信息
func GetIPStatistics(c *gin.Context) {
	ipStatistics := repository.SelectIPStatistics(global.Db)
	c.JSON(http.StatusOK, gin.H{"ipStatistics": ipStatistics})
}
