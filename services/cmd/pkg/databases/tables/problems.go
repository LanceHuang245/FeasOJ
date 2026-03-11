package tables

// Problems 表
type Problems struct {
	Id            int    `gorm:"comment:题目ID;primaryKey;autoIncrement" json:"id"`
	Difficulty    int    `gorm:"comment:难度(0:简单,1:中等,2:困难);not null" json:"difficulty"`
	Title         string `gorm:"comment:题目标题;not null" json:"title"`
	Content       string `gorm:"comment:题目详情;not null" json:"content"`
	TimeLimit     string `gorm:"comment:运行时间限制;not null" json:"time_limit"`
	MemoryLimit   string `gorm:"comment:内存大小限制;not null" json:"memory_limit"`
	Input         string `gorm:"comment:输入样例;not null" json:"input"`
	Output        string `gorm:"comment:输出样例;not null" json:"output"`
	CompetitionId int    `gorm:"comment:所属竞赛ID;not null" json:"competition_id"`
	IsVisible     bool   `gorm:"comment:是否可见;not null" json:"is_visible"`
}

// TestCases 表
type TestCases struct {
	Id         int    `gorm:"comment:测试样例ID;primaryKey;autoIncrement" json:"id"`
	ProblemId  int    `gorm:"comment:题目ID;not null" json:"problem_id"`
	InputData  string `gorm:"comment:输入数据;not null" json:"input_data"`
	OutputData string `gorm:"comment:输出数据;not null" json:"output_data"`
}
