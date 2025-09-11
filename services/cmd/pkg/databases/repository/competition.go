package repository

import (
	"FeasOJ/pkg/databases/tables"
	"FeasOJ/pkg/structs"
	"time"

	"gorm.io/gorm"
)

// 用户获取竞赛列表信息
func SelectCompetitionsInfo(db *gorm.DB) []structs.CompetitionRequest {
	var competitions []structs.CompetitionRequest
	db.Table("competitions").Where("is_visible = ?", true).Order("start_at DESC").Find(&competitions)
	return competitions
}

// 用户获取指定竞赛ID信息
func SelectCompetitionInfoByCid(db *gorm.DB, Cid int) structs.CompetitionRequest {
	var competition structs.CompetitionRequest
	db.Table("competitions").Where("contest_id = ?", Cid).Find(&competition)
	return competition
}

// 管理员获取竞赛信息
func SelectCompetitionInfoAdmin(db *gorm.DB) []structs.AdminCompetitionInfoRequest {
	var competitions []structs.AdminCompetitionInfoRequest
	db.Table("competitions").Find(&competitions)
	return competitions
}

// 管理员获取指定竞赛ID信息
func SelectCompetitionInfoAdminByCid(db *gorm.DB, Cid int) structs.AdminCompetitionInfoRequest {
	var competition structs.AdminCompetitionInfoRequest
	db.Table("competitions").Where("contest_id = ?", Cid).Find(&competition)
	return competition
}

// 管理员删除竞赛
func DeleteCompetition(db *gorm.DB, Cid int) bool {
	result := db.Table("competitions").Where("contest_id = ?", Cid).Delete(&structs.CompetitionRequest{})
	return result.RowsAffected > 0
}

// 管理员更新/添加竞赛
func UpdateCompetition(db *gorm.DB, req structs.AdminCompetitionInfoRequest) error {
	if err := db.Table("competitions").Where("contest_id = ?", req.Id).Save(&req).Error; err != nil {
		return err
	}
	return nil
}

// 将用户添加至用户-竞赛表
func AddUserCompetition(db *gorm.DB, userId int, competitionId int) error {
	var userInfo tables.User
	db.Table("users").Where("uid = ?", userId).Find(&userInfo)
	// 当前时间
	nowDateTime := time.Now()
	if err := db.Table("user_competitions").Create(
		&tables.UserCompetitions{Id: competitionId, UserId: userId, Username: userInfo.Username, JoinDate: nowDateTime}).Error; err != nil {
		return err
	}
	return nil
}

// 查询指定竞赛参加的所有用户
func SelectUsersCompetition(db *gorm.DB, competitionId int) []structs.CompetitionUserRequest {
	var users []structs.CompetitionUserRequest

	db.Table("user_competitions").
		Select("user_competitions.contest_id, user_competitions.uid, user_competitions.username, user_competitions.join_date, users.avatar").
		Joins("JOIN users ON user_competitions.uid = users.uid").
		Where("user_competitions.contest_id = ?", competitionId).
		Find(&users)

	return users
}

// 查询用户是否在指定竞赛中
func SelectUserCompetition(db *gorm.DB, userId int, competitionId int) bool {
	return db.Table("user_competitions").Where("uid = ? AND contest_id = ?", userId, competitionId).Find(&tables.UserCompetitions{}).RowsAffected > 0
}

// 将用户从用户-竞赛表删除
func DeleteUserCompetition(db *gorm.DB, userId int, competitionId int) error {
	if err := db.Table("user_competitions").Where("uid = ? AND contest_id = ?", userId, competitionId).Delete(&tables.UserCompetitions{}).Error; err != nil {
		return err
	}
	return nil
}

// 竞赛状态更新
func UpdateCompetitionStatus(db *gorm.DB) error {
	now := time.Now()

	// 状态为 1：正在进行中
	if err := db.Table("competitions").
		Where("start_at <= ? AND end_at >= ?", now, now).
		Update("status", 1).Error; err != nil {
		return err
	}

	// 状态为 2：已结束
	if err := db.Table("competitions").
		Where("end_at < ?", now).
		Update("status", 2).Error; err != nil {
		return err
	}

	// 状态为 0：未开始
	if err := db.Table("competitions").
		Where("start_at > ?", now).
		Update("status", 0).Error; err != nil {
		return err
	}

	return nil
}

// 获取未开始的竞赛
func GetUpcomingCompetitions(db *gorm.DB) []tables.Competition {
	var competitions []tables.Competition
	err := db.Where("start_at > ?", time.Now()).Find(&competitions).Error
	if err != nil {
		return nil
	}
	return competitions
}

// 获取竞赛分数情况
func GetScores(db *gorm.DB, competitionId, page, itemsPerPage int) ([]tables.UserCompetitions, int64) {
	var users []tables.UserCompetitions
	var total int64

	db.Where("contest_id = ?", competitionId).Model(&tables.UserCompetitions{}).Count(&total).Offset((page - 1) * itemsPerPage).Limit(itemsPerPage).Find(&users)

	return users, total
}
