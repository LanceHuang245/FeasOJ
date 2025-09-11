package tables

import "time"

// 讨论表
type Discussion struct {
	Id        int       `gorm:"comment:讨论ID;primaryKey;autoIncrement"`
	Title     string    `gorm:"comment:标题;not null"`
	Content   string    `gorm:"comment:内容;not null"`
	UserId    int       `gorm:"comment:用户;not null"`
	CreatedAt time.Time `gorm:"comment:创建时间;not null"`
}

// 评论表
type Comment struct {
	Id           int       `gorm:"comment:评论ID;primaryKey;autoIncrement"`
	DiscussionId int       `gorm:"comment:帖子ID;not null"`
	Content      string    `gorm:"comment:内容;not null"`
	UserId       int       `gorm:"comment:用户;not null"`
	CreatedAt    time.Time `gorm:"comment:创建时间;not null"`
	Profanity    bool      `gorm:"comment:适合展示;not null"`
}
