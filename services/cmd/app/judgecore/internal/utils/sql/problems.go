package sql

import (
	"JudgeCore/internal/global"

	"gorm.io/gorm"
)

// SelectTestCasesByPid 获取指定题目的测试样例
func SelectTestCasesByPid(db *gorm.DB, pid int) []*global.TestCaseRequest {
	var testCases []*global.TestCaseRequest
	db.Table("test_cases").Where("pid = ?", pid).Select("input_data, output_data").Find(&testCases)
	return testCases
}

// ModifyJudgeStatus 修改提交记录状态
func ModifyJudgeStatus(db *gorm.DB, Uid, Pid int, Result string) error {
	// 将result为Running...的记录修改为返回状态
	result := db.Table("submit_records").Where("uid = ? AND pid = ? AND result = ?", Uid, Pid, "Running...").Update("result", Result)
	return result.Error
}

// SelectProblemByPid 获取指定题目信息
func SelectProblemByPid(db *gorm.DB, pid int) (*global.Problem, error) {
	var problem global.Problem
	result := db.Table("problems").Where("pid = ?", pid).First(&problem)
	if result.Error != nil {
		return nil, result.Error
	}
	return &problem, nil
}