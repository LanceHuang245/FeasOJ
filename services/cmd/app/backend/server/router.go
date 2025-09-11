package server

import (
	"FeasOJ/server/handler"
	"FeasOJ/server/middlewares"

	"github.com/gin-gonic/gin"
)

func LoadRouter(r *gin.Engine) *gin.RouterGroup {
	r.Use(middlewares.Logger(), middlewares.IPStatistic())
	// 设置路由组
	router1 := r.Group("/api/v1")
	{
		// 注册
		router1.POST("/register", handler.Register)

		// 登录
		router1.GET("/login", handler.Login)

		// 获取验证码
		router1.GET("/captcha", handler.GetCaptcha)

		// 获取用户信息
		router1.GET("/users/:username", handler.GetUserInfo)

		// 密码修改
		router1.PUT("/users/password", handler.UpdatePassword)

		// 通知
		router1.GET("/notification/:uid", handler.SSEHandler)
	}

	authGroup := router1.Group("")
	authGroup.Use(middlewares.HeaderVerify())
	{
		// 验证用户信息
		authGroup.GET("/verify", handler.VerifyUserInfo)

		// 获取指定用户的提交记录
		router1.GET("/users/:username/submitrecords", handler.GetSubmitRecordsByUsername)

		// 获取指定帖子的评论
		authGroup.GET("/discussions/comments/:did", handler.GetComment)

		// 获取竞赛列表
		authGroup.GET("/competitions", handler.GetCompetitionsList)

		// 获取指定竞赛ID信息
		authGroup.GET("/competitions/info/:cid", handler.GetCompetitionInfoByID)

		// 获取竞赛参与的用户列表
		authGroup.GET("/competitions/info/:cid/users", handler.GetCompetitionUsers)

		// 获取指定竞赛的所有题目
		authGroup.GET("/competitions/info/:cid/problems", handler.GetProblemsByCompetitionID)

		// 获取用户是否在竞赛中
		authGroup.GET("/competitions/:cid/in", handler.IsInCompetition)

		// 获取所有题目
		authGroup.GET("/problems", handler.GetAllProblems)

		// 获取所有帖子
		authGroup.GET("/discussions", handler.GetAllDiscussions)

		// 获取排行榜
		authGroup.GET("/ranking", handler.GetRanking)

		// 获取指定题目ID的所有信息
		authGroup.GET("/problems/:id", handler.GetProblemInfo)

		// 获取总提交记录
		authGroup.GET("/submitrecords", handler.GetAllSubmitRecords)

		// 获取指定帖子
		authGroup.GET("/discussions/:did", handler.GetDiscussionByDid)

		// 上传代码
		authGroup.POST("/problems/:pid/code", handler.UploadCode)

		// 创建讨论
		authGroup.POST("/discussions", handler.CreateDiscussion)

		// 添加评论
		authGroup.POST("/discussions/comments/:did", handler.AddComment)

		// 加入有密码的竞赛
		authGroup.POST("/competitions/join/pwd/:cid", handler.JoinCompetitionWithPassword)

		// 加入竞赛
		authGroup.POST("/competitions/join/:cid", handler.JoinCompetition)

		// 退出竞赛
		authGroup.POST("/competitions/quit/:cid", handler.QuitCompetition)

		// 用户上传头像
		authGroup.PUT("/users/avatar", handler.UploadAvatar)

		// 简介更新
		authGroup.PUT("/users/synopsis", handler.UpdateSynopsis)

		// 删除讨论
		authGroup.DELETE("/discussions/:did", handler.DeleteDiscussion)

		// 删除评论
		authGroup.DELETE("/discussions/comments/:cid", handler.DelComment)

	}

	adminGroup := authGroup.Group("/admin")
	// 管理员权限检查
	adminGroup.Use(middlewares.PermissionChecker())
	{
		// 管理员晋升用户
		adminGroup.PUT("/users/promote", handler.PromoteUser)

		// 管理员降级用户
		adminGroup.PUT("/users/demote", handler.DemoteUser)

		// 管理员封禁用户
		adminGroup.PUT("/users/ban", handler.BanUser)

		// 管理员解封用户
		adminGroup.PUT("/users/unban", handler.UnbanUser)

		// 管理员获取竞赛列表
		adminGroup.GET("/competitions", handler.GetCompetitionListAdmin)

		// 管理员获取所有题目
		adminGroup.GET("/problems", handler.GetAllProblemsAdmin)

		// 管理员获取指定竞赛ID信息
		adminGroup.GET("/competitions/:cid", handler.GetCompetitionInfoAdmin)

		// 管理员获取指定题目的所有信息
		adminGroup.GET("/problems/:pid", handler.GetProblemAllInfo)

		// 管理员获取所有用户信息
		adminGroup.GET("/users", handler.GetAllUsersInfo)

		// 管理员新增/更新题目信息
		adminGroup.POST("/problems", handler.UpdateProblemInfo)

		// 管理员新增/更新竞赛信息
		adminGroup.POST("/competitions", handler.UpdateCompetitionInfo)

		// 管理员删除题目
		adminGroup.DELETE("/problems/:pid", handler.DeleteProblem)

		// 管理员删除竞赛
		adminGroup.DELETE("/competitions/:cid", handler.DeleteCompetition)

		// 管理员启用竞赛计分
		adminGroup.GET("/competitions/:cid/score", handler.CalculateScore)

		// 管理员查看竞赛得分情况
		adminGroup.GET("/competitions/:cid/scoreboard", handler.GetScoreBoard)

		// 管理员获取IP访问统计信息
		adminGroup.GET("/ipstats", handler.GetIPStatistics)
	}
	return router1
}
