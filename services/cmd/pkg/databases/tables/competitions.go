package tables

import "time"

// 竞赛表
type Competitions struct {
	Id           int       `gorm:"comment:比赛ID;primaryKey;autoIncrement"`
	Title        string    `gorm:"comment:标题;not null"`
	Subtitle     string    `gorm:"comment:副标题;not null"`
	Difficulty   int       `gorm:"comment:难度(0：简单，1:中等，2:困难);not null"`
	Password     string    `gorm:"comment:密码;"`
	Scored       bool      `gorm:"comment:是否已计分;not null"`
	Encrypted    bool      `gorm:"comment:是否加密;not null"`
	IsVisible    bool      `gorm:"comment:是否可见;not null"`
	Status       int       `gorm:"comment:竞赛状态(0:未开始，1:正在进行中，2:已结束);"`
	Announcement string    `gorm:"comment:公告;"`
	StartAt      time.Time `gorm:"comment:开始时间;not null"`
	EndAt        time.Time `gorm:"comment:结束时间;not null"`
}

// 用户竞赛关联表
type UserCompetitions struct {
	Id            int       `gorm:"comment:ID;primaryKey;autoIncrement"`
	CompetitionId int       `gorm:"comment:比赛ID;not null"`
	UserId        string    `gorm:"comment:用户ID;not null"`
	Username      string    `gorm:"comment:用户名;not null"`
	JoinDate      time.Time `gorm:"comment:加入时间;not null"`
	Score         int       `gorm:"comment:用户分数;"`
}
