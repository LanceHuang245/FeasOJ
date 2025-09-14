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
	// 更新题目表. GORM的Save方法会自动处理Id=0时创建新纪录，Id!=0时更新记录。
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
	if err := db.Where("problem_id = ?", problem.Id).Find(&existingTestCases).Error; err != nil {
		return err
	}

	// 找出需要删除的测试样例
	existingTestCaseMap := make(map[string]tables.TestCases)
	for _, testCase := range existingTestCases {
		existingTestCaseMap[testCase.InputData] = testCase
	}

	for _, testCase := range req.TestCases {
		delete(existingTestCaseMap, testCase.InputData)
	}

	for _, testCaseToDelete := range existingTestCaseMap {
		if err := db.Delete(&testCaseToDelete).Error; err != nil {
			return err
		}
	}

	// 更新或添加新的测试样例
	for _, testCase := range req.TestCases {
		var existingTestCase tables.TestCases
		if err := db.Where("problem_id = ? AND input_data = ?", problem.Id, testCase.InputData).First(&existingTestCase).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果测试样例不存在，则创建新的样例
				newTestCase := tables.TestCases{
					ProblemId:  problem.Id,
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
	return db.Table("problems").Where("id = ? AND is_visible = ?", id, true).First(&tables.Problems{}).Error == nil
}

// 题目状态更新
func UpdateProblemVisibility(db *gorm.DB) error {
	// 更新状态为正在进行中的题目：is_visible 为 true
	if err := db.Table("problems").
		Where("competition_id IN (SELECT competition_id FROM competitions WHERE status = ?)", 1).
		Update("is_visible", true).Error; err != nil {
		return err
	}

	// 更新状态为已结束的题目：is_visible 为 true
	if err := db.Table("problems").
		Where("competition_id IN (SELECT competition_id FROM competitions WHERE status = ?)", 2).
		Update("is_visible", true).Error; err != nil {
		return err
	}

	// 更新状态为未开始的题目：is_visible 为 false
	if err := db.Table("problems").
		Where("competition_id IN (SELECT competition_id FROM competitions WHERE status = ?)", 0).
		Update("is_visible", false).Error; err != nil {
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
	result := db.Table("problems").Where("id = ?", pid).First(&problem)
	if result.Error != nil {
		return nil, result.Error
	}
	return &problem, nil
}

// SelectRandomProblem 获取一个随机的、公开的、非竞赛的题目
func SelectRandomProblem(db *gorm.DB, dbType string) (tables.Problems, error) {
	var problem tables.Problems
	var orderClause string

	if dbType == "mysql" {
		orderClause = "RAND()"
	} else {
		orderClause = "RANDOM()"
	}

	err := db.Where("is_visible = ?", true).Order(orderClause).First(&problem).Error
	if err != nil {
		return tables.Problems{}, err
	}
	return problem, nil
}
