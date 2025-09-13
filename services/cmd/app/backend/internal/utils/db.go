package utils

import (
	"FeasOJ/app/backend/internal/config"
	"FeasOJ/app/backend/internal/global"
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

	return adminUsername, EncryptPassword(adminPassword), adminEmail, 1
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

	// 检查数据库类型，只有MySQL才需要创建触发器
	if config.GlobalConfig.Database.Type == "mysql" {
		// 检查触发器是否已存在
		var count int64
		global.Db.Raw("SELECT COUNT(*) FROM information_schema.triggers WHERE trigger_name = 'users_uuid_trigger' AND trigger_schema = DATABASE()").Scan(&count)

		// 如果触发器不存在，则创建触发器
		if count == 0 {
			// 为users表创建UUID生成触发器
			triggerSQL := `
			CREATE TRIGGER users_uuid_trigger 
			BEFORE INSERT ON users 
			FOR EACH ROW 
			BEGIN 
				IF NEW.id IS NULL OR NEW.id = '' THEN 
					SET NEW.id = UUID(); 
				END IF; 
			END;
			`

			// 执行创建触发器的SQL语句
			err = global.Db.Exec(triggerSQL).Error
			if err != nil {
				log.Printf("[FeasOJ] Failed to create UUID trigger for users table: %v", err)
				return false
			}
			log.Println("[FeasOJ] UUID trigger for users table created successfully")
		} else {
			log.Println("[FeasOJ] UUID trigger for users table already exists")
		}
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
