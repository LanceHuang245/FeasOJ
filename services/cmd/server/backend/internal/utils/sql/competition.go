package sql

import (
	"FeasOJ/internal/global"
	"time"
)

// 用户获取竞赛列表信息
func SelectCompetitionsInfo() []global.CompetitionRequest {
	var competitions []global.CompetitionRequest
	global.DB.Table("competitions").Where("is_visible = ?", true).Order("start_at DESC").Find(&competitions)
	return competitions
}

// 用户获取指定竞赛ID信息
func SelectCompetitionInfoByCid(Cid int) global.CompetitionRequest {
	var competition global.CompetitionRequest
	global.DB.Table("competitions").Where("contest_id = ?", Cid).Find(&competition)
	return competition
}

// 管理员获取竞赛信息
func SelectCompetitionInfoAdmin() []global.AdminCompetitionInfoRequest {
	var competitions []global.AdminCompetitionInfoRequest
	global.DB.Table("competitions").Find(&competitions)
	return competitions
}

// 管理员获取指定竞赛ID信息
func SelectCompetitionInfoAdminByCid(Cid int) global.AdminCompetitionInfoRequest {
	var competition global.AdminCompetitionInfoRequest
	global.DB.Table("competitions").Where("contest_id = ?", Cid).Find(&competition)
	return competition
}

// 管理员删除竞赛
func DeleteCompetition(Cid int) bool {
	result := global.DB.Table("competitions").Where("contest_id = ?", Cid).Delete(&global.CompetitionRequest{})
	return result.RowsAffected > 0
}

// 管理员更新/添加竞赛
func UpdateCompetition(req global.AdminCompetitionInfoRequest) error {
	if err := global.DB.Table("competitions").Where("contest_id = ?", req.ContestID).Save(&req).Error; err != nil {
		return err
	}
	return nil
}

// 将用户添加至用户-竞赛表
func AddUserCompetition(userId int, competitionId int) error {
	var userInfo global.User
	global.DB.Table("users").Where("uid = ?", userId).Find(&userInfo)
	// 当前时间
	nowDateTime := time.Now()
	if err := global.DB.Table("user_competitions").Create(
		&global.UserCompetitions{ContestID: competitionId, Uid: userId, Username: userInfo.Username, Join_date: nowDateTime}).Error; err != nil {
		return err
	}
	return nil
}

// 查询指定竞赛参加的所有用户
func SelectUsersCompetition(competitionId int) []global.CompetitionUserRequest {
	var users []global.CompetitionUserRequest

	global.DB.Table("user_competitions").
		Select("user_competitions.contest_id, user_competitions.uid, user_competitions.username, user_competitions.join_date, users.avatar").
		Joins("JOIN users ON user_competitions.uid = users.uid").
		Where("user_competitions.contest_id = ?", competitionId).
		Find(&users)

	return users
}

// 查询用户是否在指定竞赛中
func SelectUserCompetition(userId int, competitionId int) bool {
	return global.DB.Table("user_competitions").Where("uid = ? AND contest_id = ?", userId, competitionId).Find(&global.UserCompetitions{}).RowsAffected > 0
}

// 将用户从用户-竞赛表删除
func DeleteUserCompetition(userId int, competitionId int) error {
	if err := global.DB.Table("user_competitions").Where("uid = ? AND contest_id = ?", userId, competitionId).Delete(&global.UserCompetitions{}).Error; err != nil {
		return err
	}
	return nil
}

// 竞赛状态更新
func UpdateCompetitionStatus() error {
	now := time.Now()

	// 状态为 1：正在进行中
	if err := global.DB.Table("competitions").
		Where("start_at <= ? AND end_at >= ?", now, now).
		Update("status", 1).Error; err != nil {
		return err
	}

	// 状态为 2：已结束
	if err := global.DB.Table("competitions").
		Where("end_at < ?", now).
		Update("status", 2).Error; err != nil {
		return err
	}

	// 状态为 0：未开始
	if err := global.DB.Table("competitions").
		Where("start_at > ?", now).
		Update("status", 0).Error; err != nil {
		return err
	}

	return nil
}

// 获取未开始的竞赛
func GetUpcomingCompetitions() []global.Competition {
	var competitions []global.Competition
	err := global.DB.Where("start_at > ?", time.Now()).Find(&competitions).Error
	if err != nil {
		return nil
	}
	return competitions
}

// 获取竞赛分数情况
func GetScores(competitionId, page, itemsPerPage int) ([]global.UserCompetitions, int64) {
	var users []global.UserCompetitions
	var total int64

	db := global.DB.Where("contest_id = ?", competitionId)
	db.Model(&global.UserCompetitions{}).Count(&total)
	db.Offset((page - 1) * itemsPerPage).Limit(itemsPerPage).Find(&users)

	return users, total
}
