package tables

import "time"

// Discussions 表
type Discussions struct {
	Id        int       `gorm:"comment:讨论ID;primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"comment:标题;not null" json:"title"`
	Content   string    `gorm:"comment:内容;not null" json:"content"`
	UserId    string    `gorm:"comment:用户ID;not null" json:"user_id"`
	CreatedAt time.Time `gorm:"comment:创建时间;not null" json:"created_at"`
}

// Comments 表
type Comments struct {
	Id           int       `gorm:"comment:评论ID;primaryKey;autoIncrement" json:"id"`
	DiscussionId int       `gorm:"comment:帖子ID;not null" json:"discussion_id"`
	Content      string    `gorm:"comment:内容;not null" json:"content"`
	UserId       string    `gorm:"comment:用户ID;not null" json:"user_id"`
	CreatedAt    time.Time `gorm:"comment:创建时间;not null" json:"created_at"`
	Profanity    bool      `gorm:"comment:是否命中敏感词;not null" json:"profanity"`
}
