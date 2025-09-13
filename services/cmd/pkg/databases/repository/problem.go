package repository

import (
	"FeasOJ/pkg/databases/tables"
	"FeasOJ/pkg/structs"
	"errors"

	"gorm.io/gorm"
)

// 获取Problem表中的所有数据
func SelectAllProblems(db *gorm.DB) []tables.Problems {
	var problems []tables.Problems
	db.Where("is_visible = ?", true).Find(&problems)
	return problems
}

// 管理员获取Problem表中的所有数据
func SelectAllProblemsAdmin(db *gorm.DB) []tables.Problems {
	var problems []tables.Problems
	db.Find(&problems)
	return problems
}

// 获取指定PID的题目除了Input_full_path Output_full_path外的所有信息
func SelectProblemInfo(db *gorm.DB, id string) structs.ProblemInfoRequest {
	var problemall tables.Problems
	var problem structs.ProblemInfoRequest
	db.Table("problems").Where("id = ? AND is_visible = ?", id, true).First(&problemall)
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
func SelectProblemTestCases(db *gorm.DB, problemId string) structs.AdminProblemInfoRequest {
	var problem tables.Problems
	var testCases []structs.TestCaseRequest
	var result structs.AdminProblemInfoRequest

	if err := db.First(&problem, problemId).Error; err != nil {
		return result
	}

	if err := db.Table("test_cases").Where("problem_id = ?", problemId).Select("input_data,output_data").Find(&testCases).Error; err != nil {
		return result
	}

	result = structs.AdminProblemInfoRequest{
		Id:            problem.Id,
		Difficulty:    problem.Difficulty,
		Title:         problem.Title,
		Content:       problem.Content,
		TimeLimit:     problem.TimeLimit,
		MemoryLimit:   problem.MemoryLimit,
		Input:         problem.Input,
		Output:        problem.Output,
		CompetitionId: problem.CompetitionId,
		IsVisible:     problem.IsVisible,
		TestCases:     testCases,
	}

	return result
}

// 更新题目信息
func UpdateProblem(db *gorm.DB, req structs.AdminProblemInfoRequest) error {
	// 更新题目表
	problem := tables.Problems{
		Id:            req.Id,
		Difficulty:    req.Difficulty,
		Title:         req.Title,
		Content:       req.Content,
		TimeLimit:     req.TimeLimit,
		MemoryLimit:   req.MemoryLimit,
		Input:         req.Input,
		Output:        req.Output,
		CompetitionId: req.CompetitionId,
		IsVisible:     req.IsVisible,
	}
	if err := db.Save(&problem).Error; err != nil {
		return err
	}

	// 获取该题目的测试样例
	var existingTestCases []tables.TestCases
	if err := db.Where("problem_id = ?", req.Id).Find(&existingTestCases).Error; err != nil {
		return err
	}

	// 找出前端传入的测试样例中不存在的测试样例，并将其从数据库中删除
	existingTestCaseMap := make(map[string]tables.TestCases)
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
		var existingTestCase tables.TestCases
		if err := db.Where("problem_id = ? AND input_data = ?", req.Id, testCase.InputData).First(&existingTestCase).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果测试样例不存在，则创建新的样例
				newTestCase := tables.TestCases{
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
func DeleteProblemAllInfo(db *gorm.DB, problemId int) bool {
	if db.Table("problems").Where("id = ?", problemId).Delete(&tables.Problems{}).Error != nil {
		return false
	}

	if db.Table("test_cases").Where("problem_id = ?", problemId).Delete(&tables.TestCases{}).Error != nil {
		return false
	}

	return true
}

// 获取指定竞赛ID的所有题目列表
func SelectProblemsByCompID(db *gorm.DB, competitionId int) []structs.ProblemInfoRequest {
	var problems []structs.ProblemInfoRequest
	if err := db.Table("problems").Where("competition_id = ?", competitionId).Find(&problems).Error; err != nil {
		return nil
	}
	return problems
}

// 获取指定题目ID是否可用
func IsProblemVisible(db *gorm.DB, id int) bool {
	return db.Table("problems").Where("id = ? AND is_visible = ?", id, 1).First(&tables.Problems{}).Error == nil
}

// 题目状态更新
func UpdateProblemVisibility(db *gorm.DB) error {
	// 更新状态为正在进行中的题目：is_visible 为 1
	if err := db.Table("problems").
		Where("competition_id IN (SELECT competition_id FROM competitions WHERE status = ?)", 1).
		Update("is_visible", 1).Error; err != nil {
		return err
	}

	// 更新状态为已结束的题目：is_visible 为 1
	if err := db.Table("problems").
		Where("competition_id IN (SELECT competition_id FROM competitions WHERE status = ?)", 1).
		Update("is_visible", 1).Error; err != nil {
		return err
	}

	// 更新状态为未开始的题目：is_visible 为 0
	if err := db.Table("problems").
		Where("competition_id IN (SELECT competition_id FROM competitions WHERE status = ?)", 0).
		Update("is_visible", 0).Error; err != nil {
		return err
	}

	return nil
}

// SelectTestCasesByPid 获取指定题目的测试样例
func SelectTestCasesByPid(db *gorm.DB, problemId int) []*structs.TestCaseRequest {
	var testCases []*structs.TestCaseRequest
	db.Table("test_cases").Where("problem_id = ?", problemId).Select("input_data, output_data").Find(&testCases)
	return testCases
}

// SelectProblemByPid 获取指定题目信息
func SelectProblemByPid(db *gorm.DB, pid int) (*tables.Problems, error) {
	var problem tables.Problems
	result := db.Table("problems").Where("pid = ?", pid).First(&problem)
	if result.Error != nil {
		return nil, result.Error
	}
	return &problem, nil
}
