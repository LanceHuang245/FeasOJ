package tables

import "time"

// IPVisits 表
type IPVisits struct {
	Id         int64     `gorm:"comment:记录ID;primaryKey;autoIncrement;not null" json:"id"`
	IpAddress  string    `gorm:"comment:IP地址;primaryKey;type:varchar(45);unique;not null" json:"ip_address"`
	VisitCount int64     `gorm:"comment:访问次数;not null;default:0" json:"visit_count"`
	LastVisit  time.Time `gorm:"comment:最后访问时间;autoUpdateTime" json:"last_visit"`
}
