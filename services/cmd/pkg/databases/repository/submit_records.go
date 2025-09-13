package repository

import (
	"FeasOJ/pkg/databases/tables"
	"strings"
	"time"

	"gorm.io/gorm"
)

// 倒序查询指定用户ID的30天内的提交题目记录
func SelectSubmitRecordsByUid(db *gorm.DB, uid string) []tables.SubmitRecord {
	var records []tables.SubmitRecord
	db.Where("user_id = ?", uid).
		Where("time > ?", time.Now().Add(-30*24*time.Hour)).Order("time desc").Find(&records)
	return records
}

// 倒序查询指定用户ID的30天内的提交题目记录
// 如果题目所属竞赛正在进行中，则返回的 Code 字段为空字符串
func SelectSRByUidForChecker(db *gorm.DB, uid string) []tables.SubmitRecord {
	var records []tables.SubmitRecord
	selectFields := []string{
		"submit_records.id",
		"submit_records.problem_id",
		"submit_records.user_id",
		"submit_records.username",
		"submit_records.result",
		"submit_records.time",
		"submit_records.language",
		// 如果竞赛正在进行(status=1)，则返回空字符串，否则返回原 code
		"CASE WHEN competitions.status = 1 THEN '' ELSE submit_records.code END AS code",
	}

	db.
		Table("submit_records").
		Select(strings.Join(selectFields, ", ")).
		Joins("JOIN problems ON problems.id = submit_records.problem_id").
		Joins("JOIN competitions ON competitions.id = problems.competition_id").
		Where("submit_records.user_id = ?", uid).
		Where("submit_records.time > ?", time.Now().Add(-30*24*time.Hour)).
		Order("submit_records.time DESC").
		Find(&records)

	return records
}

// 返回 SubmitRecord 表中 30 天内的记录
func SelectAllSubmitRecords(db *gorm.DB) []tables.SubmitRecord {
	var records []tables.SubmitRecord
	selectFields := []string{
		"sr.id",
		"sr.problem_id",
		"sr.user_id",
		"sr.username",
		"sr.result",
		"sr.time",
		"sr.language",
		"CASE WHEN c.status = 1 THEN '' ELSE sr.code END AS code",
	}

	db.
		Table("submit_records AS sr").
		Select(strings.Join(selectFields, ", ")).
		Joins("JOIN problems AS p ON p.id = sr.problem_id").
		Joins("LEFT JOIN competitions AS c ON c.id = p.competition_id").
		Where("sr.time > ?", time.Now().Add(-30*24*time.Hour)).
		Where("c.status IS NULL OR c.status != 0 OR p.is_visible = ?", true).
		Order("sr.time DESC").
		Find(&records)

	return records
}

// 添加提交记录
func AddSubmitRecord(db *gorm.DB, Pid int, Uid, Result, Language, Username, Code string) bool {
	err := db.Table("submit_records").Create(&tables.SubmitRecord{UserId: Uid, ProblemId: Pid, Username: Username, Result: Result, Time: time.Now(), Language: Language, Code: Code})
	return err == nil
}

// ModifyJudgeStatus 修改提交记录状态
func ModifyJudgeStatus(db *gorm.DB, Pid int, Uid, Result string) error {
	// 将result为Running...的记录修改为返回状态
	result := db.Table("submit_records").Where("user_id = ? AND problem_id = ? AND result = ?", Uid, Pid, "Running...").Update("result", Result)
	return result.Error
}
