package global

// 测试样例请求体
type TestCaseRequest struct {
	InputData  string `json:"input"`
	OutputData string `json:"output"`
}

// 判题结果信息结构体
type JudgeResultMessage struct {
	UserID    int    `json:"user_id"`
	ProblemID int    `json:"problem_id"`
	Status    string `json:"status"`
}

// 题目表: pid, difficulty, title, content, time_limit, memory_limit, input, output, contestid, is_visible
type Problem struct {
	Pid         int    `gorm:"comment:题目ID;primaryKey;autoIncrement"`
	Difficulty  string `gorm:"comment:难度;not null"`
	Title       string `gorm:"comment:题目标题;not null"`
	Content     string `gorm:"comment:题目详细;not null"`
	Timelimit   string `gorm:"comment:运行时间限制;not null"`
	Memorylimit string `gorm:"comment:内存大小限制;not null"`
	Input       string `gorm:"comment:输入样例;not null"`
	Output      string `gorm:"comment:输出样例;not null"`
	ContestID   int    `gorm:"comment:所属竞赛ID;not null"`
	IsVisible   bool   `gorm:"comment:是否可见;not null"`
}
