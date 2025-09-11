package tables

import "time"

// 提交记录表
type SubmitRecord struct {
	Id        int       `gorm:"comment:提交ID;primaryKey;autoIncrement"`
	ProblemId int       `gorm:"comment:题目ID"`
	UserId    int       `gorm:"comment:用户ID;not null"`
	Username  string    `gorm:"comment:用户名;not null"`
	Result    string    `gorm:"comment:结果;"`
	Time      time.Time `gorm:"comment:提交时间;not null"`
	Language  string    `gorm:"comment:语言;not null"`
	Code      string    `gorm:"comment:代码;not null"`
}
