package tables

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 用户表
type Users struct {
	Id        string    `gorm:"type:varchar(36);primaryKey;not null;unique;"`
	Avatar    string    `gorm:"comment:头像存放路径"`
	Username  string    `gorm:"comment:用户名;not null;unique"`
	Password  string    `gorm:"comment:密码;not null"`
	Email     string    `gorm:"comment:电子邮件;not null;unique"`
	Synopsis  string    `gorm:"comment:简介"`
	Score     int       `gorm:"comment:分数"`
	CreatedAt time.Time `gorm:"comment:创建时间;not null"`
	Role      int       `gorm:"comment:角色(0:普通用户，1:管理员);not null"`
	IsBan     bool      `gorm:"comment:是否被封禁;not null"`
}

// BeforeCreate 在创建记录前生成UUID(PostgreSQL专属)
func (u *Users) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.New().String()
	return
}
