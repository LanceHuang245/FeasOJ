package sql

import (
	"FeasOJ/internal/global"
	"strings"
	"time"
)

// 倒序查询指定用户ID的30天内的提交题目记录
func SelectSubmitRecordsByUid(uid int) []global.SubmitRecord {
	var records []global.SubmitRecord
	global.DB.Where("uid = ?", uid).
		Where("time > ?", time.Now().Add(-30*24*time.Hour)).Order("time desc").Find(&records)
	return records
}

// 倒序查询指定用户ID的30天内的提交题目记录
// 如果题目所属竞赛正在进行中，则返回的 Code 字段为空字符串
func SelectSRByUidForChecker(uid int) []global.SubmitRecord {
	var records []global.SubmitRecord
	selectFields := []string{
		"submit_records.sid",
		"submit_records.pid",
		"submit_records.uid",
		"submit_records.username",
		"submit_records.result",
		"submit_records.time",
		"submit_records.language",
		// 如果竞赛正在进行(status=1)，则返回空字符串，否则返回原 code
		"CASE WHEN competitions.status = 1 THEN '' ELSE submit_records.code END AS code",
	}

	global.DB.
		Table("submit_records").
		Select(strings.Join(selectFields, ", ")).
		Joins("JOIN problems ON problems.pid = submit_records.pid").
		Joins("JOIN competitions ON competitions.contest_id = problems.contest_id").
		Where("submit_records.uid = ?", uid).
		Where("submit_records.time > ?", time.Now().Add(-30*24*time.Hour)).
		Order("submit_records.time DESC").
		Find(&records)

	return records
}

// 返回 SubmitRecord 表中 30 天内的记录
func SelectAllSubmitRecords() []global.SubmitRecord {
	var records []global.SubmitRecord
	selectFields := []string{
		"sr.sid",
		"sr.pid",
		"sr.uid",
		"sr.username",
		"sr.result",
		"sr.time",
		"sr.language",
		"CASE WHEN c.status = 1 THEN '' ELSE sr.code END AS code",
	}

	global.DB.
		Table("submit_records AS sr").
		Select(strings.Join(selectFields, ", ")).
		Joins("JOIN problems AS p ON p.pid = sr.pid").
		Joins("LEFT JOIN competitions AS c ON c.contest_id = p.contest_id").
		Where("sr.time > ?", time.Now().Add(-30*24*time.Hour)).
		Where("c.status IS NULL OR c.status != 0 OR p.is_visible = ?", true).
		Order("sr.time DESC").
		Find(&records)

	return records
}

// 添加提交记录
func AddSubmitRecord(Uid, Pid int, Result, Language, Username, Code string) bool {
	err := global.DB.Table("submit_records").Create(&global.SubmitRecord{Uid: Uid, Pid: Pid, Username: Username, Result: Result, Time: time.Now(), Language: Language, Code: Code})
	return err == nil
}
