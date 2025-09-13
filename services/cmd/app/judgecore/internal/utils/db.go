package utils

import (
	"FeasOJ/app/judgecore/internal/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectSql 返回数据库连接对象
func ConnectSql(dbConfig config.Database) (*gorm.DB, error) {
	// 尝试连接PostgreSQL
	if dbConfig.Type == "postgresql" {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.SSLMode)

		if dsn == "" {
			return nil, fmt.Errorf("database connection string is empty, please check config file")
		}

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("database connection failed: %w", err)
		}

		return db, nil
	}

	// 默认连接MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)

	if dsn == "" {
		return nil, fmt.Errorf("database connection string is empty, please check config file")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	return db, nil
}
