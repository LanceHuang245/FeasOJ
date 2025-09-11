package repository

import (
	"FeasOJ/pkg/databases/tables"
	"FeasOJ/pkg/structs"
	"errors"

	"gorm.io/gorm"
)

// 获取Problem表中的所有数据
func SelectAllProblems(db *gorm.DB) []tables.Problem {
	var problems []tables.Problem
	db.Where("is_visible = ?", true).Find(&problems)
	return problems
}

// 管理员获取Problem表中的所有数据
func SelectAllProblemsAdmin(db *gorm.DB) []tables.Problem {
	var problems []tables.Problem
	db.Find(&problems)
	return problems
}

// 获取指定PID的题目除了Input_full_path Output_full_path外的所有信息
func SelectProblemInfo(db *gorm.DB, pid string) structs.ProblemInfoRequest {
	var problemall tables.Problem
	var problem structs.ProblemInfoRequest
	db.Table("problems").Where("pid = ? AND is_visible = ?", pid, true).First(&problemall)
	problem = structs.ProblemInfoRequest{
		Id:          problemall.Id,
		Difficulty:  problemall.Difficulty,
		Title:       problemall.Title,
		Content:     problemall.Content,
		TimeLimit:   problemall.TimeLimit,
		MemoryLimit: problemall.MemoryLimit,
		Input:       problemall.Input,
		Output:      problemall.Output,
	}
	return problem
}

// 获取指定题目所有信息
func SelectProblemTestCases(db *gorm.DB, pid string) structs.AdminProblemInfoRequest {
	var problem tables.Problem
	var testCases []structs.TestCaseRequest
	var result structs.AdminProblemInfoRequest

	if err := db.First(&problem, pid).Error; err != nil {
		return result
	}

	if err := db.Table("test_cases").Where("pid = ?", pid).Select("input_data,output_data").Find(&testCases).Error; err != nil {
		return result
	}

	result = structs.AdminProblemInfoRequest{
		Id:          problem.Id,
		Difficulty:  problem.Difficulty,
		Title:       problem.Title,
		Content:     problem.Content,
		TimeLimit:   problem.TimeLimit,
		MemoryLimit: problem.MemoryLimit,
		Input:       problem.Input,
		Output:      problem.Output,
		ContestId:   problem.ContestId,
		IsVisible:   problem.IsVisible,
		TestCases:   testCases,
	}

	return result
}

// 更新题目信息
func UpdateProblem(db *gorm.DB, req structs.AdminProblemInfoRequest) error {
	// 更新题目表
	problem := tables.Problem{
		Id:          req.Id,
		Difficulty:  req.Difficulty,
		Title:       req.Title,
		Content:     req.Content,
		TimeLimit:   req.TimeLimit,
		MemoryLimit: req.MemoryLimit,
		Input:       req.Input,
		Output:      req.Output,
		ContestId:   req.ContestId,
		IsVisible:   req.IsVisible,
	}
	if err := db.Save(&problem).Error; err != nil {
		return err
	}

	// 获取该题目的测试样例
	var existingTestCases []tables.TestCase
	if err := db.Where("pid = ?", req.Id).Find(&existingTestCases).Error; err != nil {
		return err
	}

	// 找出前端传入的测试样例中不存在的测试样例，并将其从数据库中删除
	existingTestCaseMap := make(map[string]tables.TestCase)
	for _, testCase := range existingTestCases {
		existingTestCaseMap[testCase.InputData] = testCase
	}

	for _, testCase := range req.TestCases {
		delete(existingTestCaseMap, testCase.InputData)
	}

	for _, testCase := range existingTestCaseMap {
		if err := db.Delete(&testCase).Error; err != nil {
			return err
		}
	}

	// 更新或添加新的测试样例
	for _, testCase := range req.TestCases {
		var existingTestCase tables.TestCase
		if err := db.Where("pid = ? AND input_data = ?", req.Id, testCase.InputData).First(&existingTestCase).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果测试样例不存在，则创建新的样例
				newTestCase := tables.TestCase{
					Id:         req.Id,
					InputData:  testCase.InputData,
					OutputData: testCase.OutputData,
				}
				if err := db.Create(&newTestCase).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			// 如果测试样例存在，则更新该样例
			existingTestCase.OutputData = testCase.OutputData
			if err := db.Save(&existingTestCase).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// 删除题目及其所有测试样例
func DeleteProblemAllInfo(db *gorm.DB, pid int) bool {
	if db.Table("problems").Where("pid = ?", pid).Delete(&tables.Problem{}).Error != nil {
		return false
	}

	if db.Table("test_cases").Where("pid = ?", pid).Delete(&tables.TestCase{}).Error != nil {
		return false
	}

	return true
}

// 获取指定竞赛ID的所有题目列表
func SelectProblemsByCompID(db *gorm.DB, competitionID int) []structs.ProblemInfoRequest {
	var problems []structs.ProblemInfoRequest
	if err := db.Table("problems").Where("contest_id = ?", competitionID).Find(&problems).Error; err != nil {
		return nil
	}
	return problems
}

// 获取指定题目ID是否可用
func IsProblemVisible(db *gorm.DB, problemID int) bool {
	return db.Table("problems").Where("pid = ? AND is_visible = ?", problemID, 1).First(&tables.Problem{}).Error == nil
}

// 题目状态更新
func UpdateProblemVisibility(db *gorm.DB) error {
	// 更新状态为正在进行中的题目：is_visible 为 1
	if err := db.Table("problems").
		Where("contest_id IN (SELECT contest_id FROM competitions WHERE status = ?)", 1).
		Update("is_visible", 1).Error; err != nil {
		return err
	}

	// 更新状态为已结束的题目：is_visible 为 1
	if err := db.Table("problems").
		Where("contest_id IN (SELECT contest_id FROM competitions WHERE status = ?)", 1).
		Update("is_visible", 1).Error; err != nil {
		return err
	}

	// 更新状态为未开始的题目：is_visible 为 0
	if err := db.Table("problems").
		Where("contest_id IN (SELECT contest_id FROM competitions WHERE status = ?)", 0).
		Update("is_visible", 0).Error; err != nil {
		return err
	}

	return nil
}

// SelectTestCasesByPid 获取指定题目的测试样例
func SelectTestCasesByPid(db *gorm.DB, pid int) []*structs.TestCaseRequest {
	var testCases []*structs.TestCaseRequest
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
func SelectProblemByPid(db *gorm.DB, pid int) (*tables.Problem, error) {
	var problem tables.Problem
	result := db.Table("problems").Where("pid = ?", pid).First(&problem)
	if result.Error != nil {
		return nil, result.Error
	}
	return &problem, nil
}
