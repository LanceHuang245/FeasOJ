package repository

import (
	"FeasOJ/pkg/databases/tables"
	"FeasOJ/pkg/structs"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 用户注册
func Register(db *gorm.DB, username, password, email, tokensecret string, role int) bool {
	now := time.Now()
	err := db.Create(&tables.Users{Username: username, Password: password, Email: email, CreatedAt: now, Role: role, TokenSecret: tokensecret, IsBan: false}).Error
	return err == nil
}

// 管理员更新用户信息
func UpdateUser(db *gorm.DB, Uid int, field string, value interface{}) bool {
	return db.Table("users").Where("uid = ?", Uid).Update(field, value).Error == nil
}

// 封禁用户
func BanUser(db *gorm.DB, Uid int) bool {
	return UpdateUser(db, Uid, "is_ban", true)
}

// 解除封禁
func UnbanUser(db *gorm.DB, Uid int) bool {
	return UpdateUser(db, Uid, "is_ban", false)
}

// 晋升为管理员
func PromoteToAdmin(db *gorm.DB, Uid int) bool {
	return UpdateUser(db, Uid, "role", 1)
}

// 降级为普通用户
func DemoteToUser(db *gorm.DB, Uid int) bool {
	return UpdateUser(db, Uid, "role", 0)
}

// 管理员获取所有用户信息
func SelectAllUsersInfo(db *gorm.DB) []structs.UserInfoRequest {
	var usersInfo []structs.UserInfoRequest
	db.Table("users").Find(&usersInfo)
	return usersInfo
}

// 更新用户的头像路径
func UpdateAvatar(db *gorm.DB, username, avatarpath string) bool {
	err := db.Model(&tables.Users{}).
		Where("username = ?", username).Update("avatar", avatarpath).Error
	return err == nil
}

// 更新个人简介
func UpdateSynopsis(db *gorm.DB, username, synopsis string) bool {
	err := db.Model(&tables.Users{}).
		Where("username = ?", username).Update("synopsis", synopsis).Error
	return err == nil
}

// 根据email与username判断是否该用户已存在
func IsUserExist(db *gorm.DB, username, email string) bool {
	if db.Where("username = ?", username).
		First(&tables.Users{}).Error == nil || db.Where("email = ?", email).
		First(&tables.Users{}).Error == nil {
		return true
	}
	return false
}

// 根据邮箱获取用户信息
func SelectUserByEmail(db *gorm.DB, email string) tables.Users {
	var user tables.Users
	db.Where("email = ?", email).First(&user)
	return user
}

// 根据email修改密码
func UpdatePassword(db *gorm.DB, email, newpassword string) bool {
	tokensecret := uuid.New().String()
	err := db.Model(&tables.Users{}).
		Where("email = ?", email).Update("password", newpassword).
		Update("token_secret", tokensecret).Error
	return err == nil
}

// 根据username查询指定用户的除了password和tokensecret之外的所有信息
func SelectUserInfo(db *gorm.DB, username string) structs.UserInfoRequest {
	var user structs.UserInfoRequest
	db.Table("users").Where("username = ?", username).
		First(&user)
	return user
}

// 获取是否管理员用户
func GetAdminUser(db *gorm.DB, role int) bool {
	// role = 1表示管理员
	var user tables.Users
	err := db.Where("role = ?", role).First(&user).Error
	return err == nil
}

// 获取管理员数量
func SelectAdminCount(db *gorm.DB) int64 {
	var count int64
	db.Table("users").Where("role = ?", 1).Count(&count)
	return count
}

// 从高到低按照score排序获取前100名用户
func SelectRank100Users(db *gorm.DB) []structs.UserInfoRequest {
	var usersInfo []structs.UserInfoRequest
	db.Table("users").Order("score desc").Limit(100).Find(&usersInfo)
	return usersInfo
}

// 获取所有IP统计信息
func SelectIPStatistics(db *gorm.DB) []tables.IPVisit {
	var ipStatistics []tables.IPVisit
	db.Table("ip_visits").Order("visit_count desc").Find(&ipStatistics)
	return ipStatistics
}
