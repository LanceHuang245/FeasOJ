package tables

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Users 表
type Users struct {
	Id        string    `gorm:"type:varchar(36);primaryKey;not null;unique;" json:"id"`
	Avatar    string    `gorm:"comment:头像存放路径" json:"avatar"`
	Username  string    `gorm:"comment:用户名;not null;unique" json:"username"`
	Password  string    `gorm:"comment:密码;not null" json:"password"`
	Email     string    `gorm:"comment:电子邮件;not null;unique" json:"email"`
	Synopsis  string    `gorm:"comment:简介" json:"synopsis"`
	Score     int       `gorm:"comment:分数" json:"score"`
	CreatedAt time.Time `gorm:"comment:创建时间;not null" json:"created_at"`
	Role      int       `gorm:"comment:角色(0:普通用户,1:管理员);not null" json:"role"`
	IsBan     bool      `gorm:"comment:是否被封禁;not null" json:"is_banned"`
}

// BeforeCreate 在创建记录前生成 UUID。
func (u *Users) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.New().String()
	return
}
