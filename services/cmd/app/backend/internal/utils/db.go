package utils

import (
	"FeasOJ/app/backend/internal/config"
	"FeasOJ/app/backend/internal/global"
	"FeasOJ/pkg/auth"
	"FeasOJ/pkg/databases/tables"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 创建管理员
func InitAdminAccount() (string, string, string, int) {
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

	return adminUsername, auth.EncryptPassword(adminPassword), adminEmail, 1
}

// 创建表
func InitTable() bool {
	err := global.Db.AutoMigrate(
		&tables.Users{},
		&tables.Problems{},
		&tables.SubmitRecord{},
		&tables.Discussions{},
		&tables.Comments{},
		&tables.TestCases{},
		&tables.Competitions{},
		&tables.UserCompetitions{},
		&tables.IPVisits{},
	)

	if err != nil {
		return false
	}
	return true
}

// 返回数据库连接对象
func ConnectSql() *gorm.DB {
	dsn := config.GetDatabaseDSN()
	if dsn == "" {
		log.Println("[FeasOJ] Database connection failed, please check config.toml configuration.")
		return nil
	}

	var db *gorm.DB
	var err error

	// 根据数据库类型连接相应的数据库
	switch config.GlobalConfig.Database.Type {
	case "postgresql":
		log.Println("[FeasOJ] Trying to connect to PostgreSQL...")
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("[FeasOJ] Failed to connect to PostgreSQL: %v", err)
			log.Println("[FeasOJ] Database connection failed, please check config.toml configuration.")
			return nil
		}
		log.Println("[FeasOJ] Connected to PostgreSQL successfully")
	default: // 默认为MySQL
		log.Println("[FeasOJ] Trying to connect to MySQL...")
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("[FeasOJ] Failed to connect to MySQL: %v", err)
			log.Println("[FeasOJ] Database connection failed, please check config.toml configuration.")
			return nil
		}
		log.Println("[FeasOJ] Connected to MySQL successfully")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Println("[FeasOJ] Failed to get generic database object.")
		return nil
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(config.GlobalConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.GlobalConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.GlobalConfig.Database.MaxLifeTime) * time.Second)

	return db
}
