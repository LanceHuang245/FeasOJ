package tables

// 题目表
type Problem struct {
	Id          int    `gorm:"comment:题目ID;primaryKey;autoIncrement"`
	Difficulty  int    `gorm:"comment:难度(0：简单，1:中等，2:困难);not null"`
	Title       string `gorm:"comment:题目标题;not null"`
	Content     string `gorm:"comment:题目详细;not null"`
	TimeLimit   string `gorm:"comment:运行时间限制;not null"`
	MemoryLimit string `gorm:"comment:内存大小限制;not null"`
	Input       string `gorm:"comment:输入样例;not null"`
	Output      string `gorm:"comment:输出样例;not null"`
	ContestId   int    `gorm:"comment:所属竞赛ID;not null"`
	IsVisible   bool   `gorm:"comment:是否可见;not null"`
}

// 测试样例表
type TestCase struct {
	Id         int    `gorm:"comment:测试样例ID;primaryKey;autoIncrement"`
	ProblemId  int    `gorm:"comment:题目ID;not null"`
	InputData  string `gorm:"comment:输入数据;not null"`
	OutputData string `gorm:"comment:输出数据;not null"`
}
