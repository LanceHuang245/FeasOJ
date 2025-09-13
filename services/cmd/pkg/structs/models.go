package structs

import (
	"time"
)

// 判题结果消息结构体
type JudgeResultMessage struct {
	UserID    string `json:"user_id"`
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
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
}

// 获取讨论评论请求体
type CommentRequest struct {
	Id           int    `json:"id"`
	DiscussionId int    `json:"discussion_id"`
	Content      string `json:"content"`
	UserId       string `json:"user_id"`
	Username     string `json:"username"`
	Avatar       string `json:"avatar"`
	CreatedAt    string `json:"created_at"`
	Profanity    bool   `json:"profanity"`
}

// 用户信息请求体
type UserInfoRequest struct {
	Id        string    `json:"id"`
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
	Id            int               `json:"id"`
	Difficulty    int               `json:"difficulty"`
	Title         string            `json:"title"`
	Content       string            `json:"content"`
	TimeLimit     string            `json:"time_limit"`
	MemoryLimit   string            `json:"memory_limit"`
	Input         string            `json:"input"`
	Output        string            `json:"output"`
	CompetitionId int               `json:"competition_id"`
	IsVisible     bool              `json:"is_visible"`
	TestCases     []TestCaseRequest `json:"test_cases"`
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
	UserId   string    `json:"user_id"`
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
	UserId   string    `json:"user_id"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	JoinDate time.Time `json:"join_date"`
}
