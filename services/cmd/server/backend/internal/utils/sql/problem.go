package sql

import (
	"FeasOJ/internal/global"
	"errors"

	"gorm.io/gorm"
)

// 获取Problem表中的所有数据
func SelectAllProblems() []global.Problem {
	var problems []global.Problem
	global.DB.Where("is_visible = ?", true).Find(&problems)
	return problems
}

// 管理员获取Problem表中的所有数据
func SelectAllProblemsAdmin() []global.Problem {
	var problems []global.Problem
	global.DB.Find(&problems)
	return problems
}

// 获取指定PID的题目除了Input_full_path Output_full_path外的所有信息
func SelectProblemInfo(pid string) global.ProblemInfoRequest {
	var problemall global.Problem
	var problem global.ProblemInfoRequest
	global.DB.Table("problems").Where("pid = ? AND is_visible = ?", pid, true).First(&problemall)
	problem = global.ProblemInfoRequest{
		Pid:         problemall.Pid,
		Difficulty:  problemall.Difficulty,
		Title:       problemall.Title,
		Content:     problemall.Content,
		Timelimit:   problemall.Timelimit,
		Memorylimit: problemall.Memorylimit,
		Input:       problemall.Input,
		Output:      problemall.Output,
	}
	return problem
}

// 获取指定题目所有信息
func SelectProblemTestCases(pid string) global.AdminProblemInfoRequest {
	var problem global.Problem
	var testCases []global.TestCaseRequest
	var result global.AdminProblemInfoRequest

	if err := global.DB.First(&problem, pid).Error; err != nil {
		return result
	}

	if err := global.DB.Table("test_cases").Where("pid = ?", pid).Select("input_data,output_data").Find(&testCases).Error; err != nil {
		return result
	}

	result = global.AdminProblemInfoRequest{
		Pid:         problem.Pid,
		Difficulty:  problem.Difficulty,
		Title:       problem.Title,
		Content:     problem.Content,
		Timelimit:   problem.Timelimit,
		Memorylimit: problem.Memorylimit,
		Input:       problem.Input,
		Output:      problem.Output,
		ContestID:   problem.ContestID,
		IsVisible:   problem.IsVisible,
		TestCases:   testCases,
	}

	return result
}

// 更新题目信息
func UpdateProblem(req global.AdminProblemInfoRequest) error {
	// 更新题目表
	problem := global.Problem{
		Pid:         req.Pid,
		Difficulty:  req.Difficulty,
		Title:       req.Title,
		Content:     req.Content,
		Timelimit:   req.Timelimit,
		Memorylimit: req.Memorylimit,
		Input:       req.Input,
		Output:      req.Output,
		ContestID:   req.ContestID,
		IsVisible:   req.IsVisible,
	}
	if err := global.DB.Save(&problem).Error; err != nil {
		return err
	}

	// 获取该题目的测试样例
	var existingTestCases []global.TestCase
	if err := global.DB.Where("pid = ?", req.Pid).Find(&existingTestCases).Error; err != nil {
		return err
	}

	// 找出前端不存在的测试样例，并将其从数据库中删除
	existingTestCaseMap := make(map[string]global.TestCase)
	for _, testCase := range existingTestCases {
		existingTestCaseMap[testCase.InputData] = testCase
	}

	for _, testCase := range req.TestCases {
		delete(existingTestCaseMap, testCase.InputData)
	}

	for _, testCase := range existingTestCaseMap {
		if err := global.DB.Delete(&testCase).Error; err != nil {
			return err
		}
	}

	// 更新或添加新的测试样例
	for _, testCase := range req.TestCases {
		var existingTestCase global.TestCase
		if err := global.DB.Where("pid = ? AND input_data = ?", req.Pid, testCase.InputData).First(&existingTestCase).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果测试样例不存在，则创建新的样例
				newTestCase := global.TestCase{
					Pid:        req.Pid,
					InputData:  testCase.InputData,
					OutputData: testCase.OutputData,
				}
				if err := global.DB.Create(&newTestCase).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			// 如果测试样例存在，则更新该样例
			existingTestCase.OutputData = testCase.OutputData
			if err := global.DB.Save(&existingTestCase).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// 删除题目及其所有测试样例
func DeleteProblemAllInfo(pid int) bool {
	if global.DB.Table("problems").Where("pid = ?", pid).Delete(&global.Problem{}).Error != nil {
		return false
	}

	if global.DB.Table("test_cases").Where("pid = ?", pid).Delete(&global.TestCase{}).Error != nil {
		return false
	}

	return true
}

// 获取指定竞赛ID的所有题目列表
func SelectProblemsByCompID(competitionID int) []global.ProblemInfoRequest {
	var problems []global.ProblemInfoRequest
	if err := global.DB.Table("problems").Where("contest_id = ?", competitionID).Find(&problems).Error; err != nil {
		return nil
	}
	return problems
}

// 获取指定题目ID是否可用
func IsProblemVisible(problemID int) bool {
	return global.DB.Table("problems").Where("pid = ? AND is_visible = ?", problemID, 1).First(&global.Problem{}).Error == nil
}

// 题目状态更新
func UpdateProblemVisibility() error {
	// 更新状态为正在进行中的题目：is_visible 为 1
	if err := global.DB.Table("problems").
		Where("contest_id IN (SELECT contest_id FROM competitions WHERE status = ?)", 1).
		Update("is_visible", 1).Error; err != nil {
		return err
	}

	// 更新状态为已结束的题目：is_visible 为 1
	if err := global.DB.Table("problems").
		Where("contest_id IN (SELECT contest_id FROM competitions WHERE status = ?)", 1).
		Update("is_visible", 1).Error; err != nil {
		return err
	}

	// 更新状态为未开始的题目：is_visible 为 0
	if err := global.DB.Table("problems").
		Where("contest_id IN (SELECT contest_id FROM competitions WHERE status = ?)", 0).
		Update("is_visible", 0).Error; err != nil {
		return err
	}

	return nil
}
