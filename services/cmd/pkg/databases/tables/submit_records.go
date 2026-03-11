package tables

import "time"

// SubmitRecord 表
type SubmitRecord struct {
	Id        int       `gorm:"comment:提交ID;primaryKey;autoIncrement" json:"id"`
	ProblemId int       `gorm:"comment:题目ID" json:"problem_id"`
	UserId    string    `gorm:"comment:用户ID;not null" json:"user_id"`
	Username  string    `gorm:"comment:用户名;not null" json:"username"`
	Result    string    `gorm:"comment:结果;" json:"result"`
	Time      time.Time `gorm:"comment:提交时间;not null" json:"time"`
	Language  string    `gorm:"comment:语言;not null" json:"language"`
	Code      string    `gorm:"comment:代码;not null" json:"code"`
}
