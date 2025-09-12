package tables

import "time"

// IP访问记录表
type IPVisits struct {
	Id         int64     `gorm:"comment:记录ID;primaryKey;autoIncrement"`
	IpAddress  string    `gorm:"comment:IP地址;primaryKey;type:varchar(45)"`
	VisitCount int64     `gorm:"comment:访问次数;not null;default:0"`
	LastVisit  time.Time `gorm:"comment:最后调用;autoUpdateTime"`
}
