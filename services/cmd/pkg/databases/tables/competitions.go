package tables

import "time"

// Competitions 表
type Competitions struct {
	Id           int       `gorm:"comment:比赛ID;primaryKey;autoIncrement" json:"id"`
	Title        string    `gorm:"comment:标题;not null" json:"title"`
	Subtitle     string    `gorm:"comment:副标题;not null" json:"subtitle"`
	Difficulty   int       `gorm:"comment:难度(0:简单,1:中等,2:困难);not null" json:"difficulty"`
	Password     string    `gorm:"comment:密码;" json:"password"`
	Scored       bool      `gorm:"comment:是否已计分;not null" json:"scored"`
	Encrypted    bool      `gorm:"comment:是否加密;not null" json:"encrypted"`
	IsVisible    bool      `gorm:"comment:是否可见;not null" json:"is_visible"`
	Status       int       `gorm:"comment:比赛状态(0:未开始,1:进行中,2:已结束);" json:"status"`
	Announcement string    `gorm:"comment:公告;" json:"announcement"`
	StartAt      time.Time `gorm:"comment:开始时间;not null" json:"start_at"`
	EndAt        time.Time `gorm:"comment:结束时间;not null" json:"end_at"`
}

// UserCompetitions 表
type UserCompetitions struct {
	Id            int       `gorm:"comment:ID;primaryKey;autoIncrement" json:"id"`
	CompetitionId int       `gorm:"comment:比赛ID;not null" json:"competition_id"`
	UserId        string    `gorm:"comment:用户ID;not null" json:"user_id"`
	Username      string    `gorm:"comment:用户名;not null" json:"username"`
	JoinDate      time.Time `gorm:"comment:加入时间;not null" json:"join_date"`
	Score         int       `gorm:"comment:用户分数;" json:"score"`
}
