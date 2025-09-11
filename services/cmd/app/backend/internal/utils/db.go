package utils

import (
	"FeasOJ/internal/config"
	"FeasOJ/internal/global"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 创建管理员
func InitAdminAccount() (string, string, string, string, int) {
	var adminUsername string
	var adminPassword string
	var adminEmail string
	log.Println("[FeasOJ] Please input the administrator account configuration: ")
	fmt.Print("[FeasOJ] Username: ")
	fmt.Scanln(&adminUsername)
	fmt.Print("[FeasOJ] Password: ")
	fmt.Scanln(&adminPassword)
	fmt.Print("[FeasOJ] Email: ")
	fmt.Scanln(&adminEmail)

	return adminUsername, EncryptPassword(adminPassword), adminEmail, uuid.New().String(), 1
}

// 创建表
func InitTable() bool {
	err := global.DB.AutoMigrate(
		&global.User{},
		&global.Problem{},
		&global.SubmitRecord{},
		&global.Discussion{},
		&global.Comment{},
		&global.TestCase{},
		&global.Competition{},
		&global.UserCompetitions{},
		&global.IPVisit{},
	)
	return err == nil
}

// 返回数据库连接对象
func ConnectSql() *gorm.DB {
	dsn := config.GetMySQLDSN()
	if dsn == "" {
		log.Println("[FeasOJ] Database connection failed, please check config.json configuration.")
		return nil
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("[FeasOJ] Database connection failed, please check config.json configuration.")
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("[FeasOJ] Failed to get generic database object.")
		return nil
	}

	sqlDB.SetMaxIdleConns(config.GlobalConfig.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.GlobalConfig.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.GlobalConfig.MySQL.MaxLifeTime) * time.Second)
	return db
}

// 根据用户名获取用户信息
func SelectUser(username string) global.User {
	var user global.User
	global.DB.Where("username = ?", username).First(&user)
	return user
}
