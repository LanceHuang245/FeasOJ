package tables

import "time"

// 用户表
type Users struct {
	Id        int       `gorm:"comment:用户ID;primaryKey;autoIncrement"`
	Avatar    string    `gorm:"comment:头像存放路径"`
	Username  string    `gorm:"comment:用户名;not null;unique"`
	Password  string    `gorm:"comment:密码;not null"`
	Email     string    `gorm:"comment:电子邮件;not null"`
	Synopsis  string    `gorm:"comment:简介"`
	Score     int       `gorm:"comment:分数"`
	CreatedAt time.Time `gorm:"comment:创建时间;not null"`
	Role      int       `gorm:"comment:角色(0:普通用户，1:管理员);not null"`
	// TokenSecret string    `gorm:"comment:token密钥;not null"`
	IsBan bool `gorm:"comment:是否被封禁;not null"`
}
