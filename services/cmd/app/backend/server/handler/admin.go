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
	c.JSON(http.StatusOK, gin.H{"data": problems})
}

// 管理员获取指定题目所有信息
func GetProblemAllInfo(c *gin.Context) {
	problemInfo := repository.SelectProblemTestCases(global.Db, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"data": problemInfo})
}

// 新增/更新题目信息
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
	problemId, _ := strconv.Atoi(c.Param("id"))
	if !repository.DeleteProblemAllInfo(global.Db, problemId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 管理员获取所有用户信息
func GetAllUsersInfo(c *gin.Context) {
	usersInfo := repository.SelectAllUsersInfo(global.Db)
	c.JSON(http.StatusOK, gin.H{"data": usersInfo})
}

// 晋升/降级用户
func ChangeUserPrivilege(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))
	action, _ := strconv.Atoi(c.Query("action"))
	if !repository.ChangePrivilege(global.Db, userId, action) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 封禁/解禁用户
func ChangeUserStatus(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))
	status, _ := strconv.ParseBool(c.Query("status"))
	if !repository.ChangeUserStatus(global.Db, userId, status) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 管理员获取竞赛列表
func GetCompetitionListAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": repository.SelectCompetitionInfoAdmin(global.Db)})
}

// 管理员获取指定竞赛ID信息
func GetCompetitionInfoAdmin(c *gin.Context) {
	competitionId, _ := strconv.Atoi(c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"data": repository.SelectCompetitionInfoAdminByCid(global.Db, competitionId)})
}

// 删除指定ID竞赛
func DeleteCompetition(c *gin.Context) {
	competitionId, _ := strconv.Atoi(c.Param("id"))

	if !repository.DeleteCompetition(global.Db, competitionId) {
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
	competitionId, _ := strconv.Atoi(c.Param("id"))

	// 查询竞赛信息
	var competition tables.Competitions
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
			Where("user_id = ? AND result = ? AND time BETWEEN ? AND ?",
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
				Where("competition_id = ? AND problem_id = ?", competitionId, submission.ProblemId).
				Row().
				Scan(&difficulty)
			if err != nil {
				continue
			}

			switch difficulty {
			case "0":
				score += 1
			case "1":
				score += 3
			case "2":
				score += 5
			}
		}

		// 更新用户分数
		if score > 0 {
			global.Db.Model(&tables.Users{}).Where("user_id = ?", user.UserId).Update("score", gorm.Expr("score + ?", score))
		}

		global.Db.Model(&tables.UserCompetitions{}).Where("user_id = ?", user.UserId).Update("score", score)
	}

	global.Db.Model(&tables.Competitions{}).Where("competition_id = ?", competitionId).Update("scored", true)

	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 查询指定竞赛中，参与人员的得分情况
func GetScoreBoard(c *gin.Context) {
	competitionId, _ := strconv.Atoi(c.Param("id"))
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
	c.JSON(http.StatusOK, gin.H{"data": ipStatistics})
}
