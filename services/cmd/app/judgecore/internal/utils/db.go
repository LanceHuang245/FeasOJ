package utils

import (
	"FeasOJ/app/judgecore/internal/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectSql 返回数据库连接对象
func ConnectSql(dbConfig config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Name)

	if dsn == "" {
		return nil, fmt.Errorf("database connection string is empty, please check config file")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	return db, nil
}
