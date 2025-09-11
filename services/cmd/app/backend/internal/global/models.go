package global

import (
	"time"
)

// 判题结果消息结构体
type JudgeResultMessage struct {
	UserID    int    `json:"user_id"`
	ProblemID int    `json:"problem_id"`
	Status    string `json:"status"`
}

// 注册请求体
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}

// 修改密码请求体
type UpdatePasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
	Captcha     string `json:"Captcha"`
}

// 讨论列表请求体
type DiscussRequest struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

// 讨论帖子页请求体
type DiscsInfoRequest struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
}

// 获取讨论评论请求体
type CommentRequest struct {
	Id           int    `json:"id"`
	DiscussionId int    `json:"discussion_id"`
	Content      string `json:"content"`
	UserId       int    `json:"user_id"`
	Username     string `json:"username"`
	Avatar       string `json:"avatar"`
	CreatedAt    string `json:"created_at"`
	Profanity    bool   `json:"profanity"`
}

// 用户信息请求体
type UserInfoRequest struct {
	Id        int       `json:"id"`
	Avatar    string    `json:"avatar"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Synopsis  string    `json:"synopsis"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	Role      int       `json:"role"`
	IsBan     bool      `json:"is_Ban"`
}

// 题目信息请求体
type ProblemInfoRequest struct {
	Id          int    `json:"id"`
	Difficulty  int    `json:"difficulty"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	TimeLimit   string `json:"time_limit"`
	MemoryLimit string `json:"memory_limit"`
	Input       string `json:"input"`
	Output      string `json:"output"`
	ContestId   int    `json:"contest_id"`
}

// 测试样例请求体
type TestCaseRequest struct {
	InputData  string `json:"input_data"`
	OutputData string `json:"output_data"`
}

// 管理员获取题目信息请求体
type AdminProblemInfoRequest struct {
	Id          int               `json:"id"`
	Difficulty  int               `json:"difficulty"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	TimeLimit   string            `json:"time_limit"`
	MemoryLimit string            `json:"memory_limit"`
	Input       string            `json:"input"`
	Output      string            `json:"output"`
	ContestId   int               `json:"contest_id"`
	IsVisible   bool              `json:"is_visible"`
	TestCases   []TestCaseRequest `json:"test_cases"`
}

// 管理员获取竞赛信息请求体
type AdminCompetitionInfoRequest struct {
	Id           int64     `json:"id"`
	Title        string    `json:"title"`
	Subtitle     string    `json:"subtitle"`
	Difficulty   int       `json:"difficulty"`
	Password     string    `json:"password"`
	Status       int       `json:"status"`
	Scored       bool      `json:"scored"`
	Crypto       bool      `json:"crypto"`
	Announcement string    `json:"announcement"`
	IsVisible    bool      `json:"is_visible"`
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
}

// 管理员获取竞赛情况请求体
type AdminCompetitionScoreRequest struct {
	Id       int       `json:"id"`
	UserId   int       `json:"user_id"`
	Username string    `json:"username"`
	Score    int       `json:"score"`
	JoinDate time.Time `json:"join_date"`
}

// 用户获取竞赛请求体
type CompetitionRequest struct {
	Id           int       `json:"id"`
	Title        string    `json:"title"`
	Subtitle     string    `json:"subtitle"`
	Difficulty   int       `json:"difficulty"`
	Crypto       bool      `json:"crypto"`
	Announcement string    `json:"announcement"`
	Status       int       `json:"status"`
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
}

// 获取参赛人员请求体
type CompetitionUserRequest struct {
	Id       int       `json:"id"`
	UserId   int       `json:"user_id"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	JoinDate time.Time `json:"join_date"`
}

// 用户表
type User struct {
	Id          int       `gorm:"comment:用户ID;primaryKey;autoIncrement"`
	Avatar      string    `gorm:"comment:头像存放路径"`
	Username    string    `gorm:"comment:用户名;not null;unique"`
	Password    string    `gorm:"comment:密码;not null"`
	Email       string    `gorm:"comment:电子邮件;not null"`
	Synopsis    string    `gorm:"comment:简介"`
	Score       int       `gorm:"comment:分数"`
	CreatedAt   time.Time `gorm:"comment:创建时间;not null"`
	Role        int       `gorm:"comment:角色(0:普通用户，1:管理员);not null"`
	TokenSecret string    `gorm:"comment:token密钥;not null"`
	IsBan       bool      `gorm:"comment:是否被封禁;not null"`
}

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

// 提交记录表
type SubmitRecord struct {
	Id        int       `gorm:"comment:提交ID;primaryKey;autoIncrement"`
	ProblemId int       `gorm:"comment:题目ID"`
	UserId    int       `gorm:"comment:用户ID;not null"`
	Username  string    `gorm:"comment:用户名;not null"`
	Result    string    `gorm:"comment:结果;"`
	Time      time.Time `gorm:"comment:提交时间;not null"`
	Language  string    `gorm:"comment:语言;not null"`
	Code      string    `gorm:"comment:代码;not null"`
}

// 讨论表
type Discussion struct {
	Id        int       `gorm:"comment:讨论ID;primaryKey;autoIncrement"`
	Title     string    `gorm:"comment:标题;not null"`
	Content   string    `gorm:"comment:内容;not null"`
	UserId    int       `gorm:"comment:用户;not null"`
	CreatedAt time.Time `gorm:"comment:创建时间;not null"`
}

// 评论表
type Comment struct {
	Id           int       `gorm:"comment:评论ID;primaryKey;autoIncrement"`
	DiscussionId int       `gorm:"comment:帖子ID;not null"`
	Content      string    `gorm:"comment:内容;not null"`
	UserId       int       `gorm:"comment:用户;not null"`
	CreatedAt    time.Time `gorm:"comment:创建时间;not null"`
	Profanity    bool      `gorm:"comment:适合展示;not null"`
}

// 测试样例表
type TestCase struct {
	Id         int    `gorm:"comment:测试样例ID;primaryKey;autoIncrement"`
	ProblemId  int    `gorm:"comment:题目ID;not null"`
	InputData  string `gorm:"comment:输入数据;not null"`
	OutputData string `gorm:"comment:输出数据;not null"`
}

// 竞赛表
type Competition struct {
	Id           int       `gorm:"comment:比赛ID;primaryKey;autoIncrement"`
	Title        string    `gorm:"comment:标题;not null"`
	Subtitle     string    `gorm:"comment:副标题;not null"`
	Difficulty   int       `gorm:"comment:难度(0：简单，1:中等，2:困难);not null"`
	Password     string    `gorm:"comment:密码;"`
	Scored       bool      `gorm:"comment:是否已计分;not null"`
	Crypto       bool      `gorm:"comment:是否加密;not null"`
	IsVisible    bool      `gorm:"comment:是否可见;not null"`
	Status       int       `gorm:"comment:竞赛状态(0:未开始，1:正在进行中，2:已结束);"`
	Announcement string    `gorm:"comment:公告;"`
	StartAt      time.Time `gorm:"comment:开始时间;not null"`
	EndAt        time.Time `gorm:"comment:结束时间;not null"`
}

// 用户竞赛关联表
type UserCompetitions struct {
	Id       int       `gorm:"comment:比赛ID;not null"`
	UserId   int       `gorm:"comment:用户ID;not null"`
	Username string    `gorm:"comment:用户名;not null"`
	JoinDate time.Time `gorm:"comment:加入时间;not null"`
	Score    int       `gorm:"comment:用户分数;"`
}

// IP访问记录表
type IPVisit struct {
	Id         int64     `gorm:"comment:记录ID;primaryKey;autoIncrement"`
	IpAddress  string    `gorm:"comment:IP地址;primaryKey;type:varchar(45)"`
	VisitCount int64     `gorm:"comment:访问次数;not null;default:0"`
	LastVisit  time.Time `gorm:"comment:最后调用;autoUpdateTime"`
}
