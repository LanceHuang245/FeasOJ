package global

import (
	"time"
)

// IP访问记录表：ip, visit_count, last_visit
type IPVisit struct {
	IP         string    `gorm:"comment:IP;primaryKey;type:varchar(45)"`
	VisitCount int64     `gorm:"comment:调用次数;not null;default:0"`
	LastVisit  time.Time `gorm:"comment:最后调用;autoUpdateTime"`
}

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
	Vcode    string `json:"vcode"`
}

// 修改密码请求体
type UpdatePasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"newpassword"`
	Vcode       string `json:"vcode"`
}

// 讨论列表请求体
type DiscussRequest struct {
	Did       string `json:"did"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	Create_at string `json:"create_at"`
}

// 讨论帖子页请求体
type DiscsInfoRequest struct {
	Did       int       `json:"did"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Uid       int       `json:"uid"`
	Create_at time.Time `json:"create_at"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
}

// 获取讨论评论请求体
type CommentRequest struct {
	Cid       int    `json:"cid"`
	Did       int    `json:"did"`
	Content   string `json:"content"`
	Uid       int    `json:"uid"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	Create_at string `json:"create_at"`
	Profanity bool   `json:"profanity"`
}

// 用户信息请求体
type UserInfoRequest struct {
	Uid      int       `json:"uid"`
	Avatar   string    `json:"avatar"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Synopsis string    `json:"synopsis"`
	Score    int       `json:"score"`
	CreateAt time.Time `json:"create_at"`
	Role     int       `json:"role"`
	IsBan    bool      `json:"isBan"`
}

// 题目信息请求体
type ProblemInfoRequest struct {
	Pid         int    `json:"pid"`
	Difficulty  string `json:"difficulty"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Timelimit   string `json:"time_limit"`
	Memorylimit string `json:"memory_limit"`
	Input       string `json:"input"`
	Output      string `json:"output"`
	ContestID   int    `json:"contest_id"`
}

// 测试样例请求体
type TestCaseRequest struct {
	InputData  string `json:"input"`
	OutputData string `json:"output"`
}

// 管理员获取题目信息请求体
type AdminProblemInfoRequest struct {
	Pid         int               `json:"pid"`
	Difficulty  string            `json:"difficulty"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	Timelimit   string            `json:"time_limit"`
	Memorylimit string            `json:"memory_limit"`
	Input       string            `json:"input"`
	Output      string            `json:"output"`
	ContestID   int               `json:"contest_id"`
	IsVisible   bool              `json:"is_visible"`
	TestCases   []TestCaseRequest `json:"test_cases"`
}

// 管理员获取竞赛信息请求体
type AdminCompetitionInfoRequest struct {
	ContestID    int       `json:"contest_id"`
	Title        string    `json:"title"`
	Subtitle     string    `json:"subtitle"`
	Difficulty   string    `json:"difficulty"`
	Password     string    `json:"password"`
	Status       int       `json:"status"`
	Scored       bool      `json:"scored"`
	HavePassword bool      `json:"have_password"`
	Announcement string    `json:"announcement"`
	IsVisible    bool      `json:"is_visible"`
	Start_at     time.Time `json:"start_at"`
	End_at       time.Time `json:"end_at"`
}

// 管理员获取竞赛情况请求体
type AdminCompetitionScoreRequest struct {
	ContestID int       `json:"contest_id"`
	Uid       int       `json:"uid"`
	Username  string    `json:"username"`
	Score     int       `json:"score"`
	JoinDate  time.Time `json:"join_date"`
}

// 用户获取竞赛请求体
type CompetitionRequest struct {
	ContestID    int       `json:"contest_id"`
	Title        string    `json:"title"`
	Subtitle     string    `json:"subtitle"`
	Difficulty   string    `json:"difficulty"`
	HavePassword bool      `json:"have_password"`
	Announcement string    `json:"announcement"`
	Status       int       `json:"status"`
	Start_at     time.Time `json:"start_at"`
	End_at       time.Time `json:"end_at"`
}

// 获取参赛人员请求体
type CompetitionUserRequest struct {
	ContestID int       `json:"contest_id"`
	Uid       int       `json:"uid"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	JoinDate  time.Time `json:"join_date"`
}

// 用户表：uid, avatar, username, password, email, score, synopsis, create_at, role, token_secret, is_ban
type User struct {
	Uid         int       `gorm:"comment:用户ID;primaryKey;autoIncrement"`
	Avatar      string    `gorm:"comment:头像存放路径"`
	Username    string    `gorm:"comment:用户名;not null;unique"`
	Password    string    `gorm:"comment:密码;not null"`
	Email       string    `gorm:"comment:电子邮件;not null"`
	Synopsis    string    `gorm:"comment:简介"`
	Score       int       `gorm:"comment:分数"`
	Create_at   time.Time `gorm:"comment:创建时间;not null"`
	Role        int       `gorm:"comment:角色;not null"` // 0: 普通用户, 1: 管理员
	TokenSecret string    `gorm:"comment:token密钥;not null"`
	IsBan       bool      `gorm:"comment:是否被封禁;not null"`
}

// 题目表：pid, difficulty, title, content, time_limit, memory_limit, input, output, contest_id, is_visible
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

// 提交记录表：sid, pid, uid, username, result, time, language, code
type SubmitRecord struct {
	Sid      int       `gorm:"comment:提交ID;primaryKey;autoIncrement"`
	Pid      int       `gorm:"comment:题目ID"`
	Uid      int       `gorm:"comment:用户ID;not null"`
	Username string    `gorm:"comment:用户名;not null"`
	Result   string    `gorm:"comment:结果;"`
	Time     time.Time `gorm:"comment:提交时间;not null"`
	Language string    `gorm:"comment:语言;not null"`
	Code     string    `gorm:"comment:代码;not null"`
}

// 讨论表：did, title, content, uid, create_at
type Discussion struct {
	Did       int       `gorm:"comment:讨论ID;primaryKey;autoIncrement"`
	Title     string    `gorm:"comment:标题;not null"`
	Content   string    `gorm:"comment:内容;not null"`
	Uid       int       `gorm:"comment:用户;not null"`
	Create_at time.Time `gorm:"comment:创建时间;not null"`
}

// 评论表：cid, did, content, uid, create_at, profanity
type Comment struct {
	Cid       int       `gorm:"comment:评论ID;primaryKey;autoIncrement"`
	Did       int       `gorm:"comment:帖子ID;not null"`
	Content   string    `gorm:"comment:内容;not null"`
	Uid       int       `gorm:"comment:用户;not null"`
	Create_at time.Time `gorm:"comment:创建时间;not null"`
	Profanity bool      `gorm:"comment:适合展示;not null"`
}

// 测试样例表：tid, pid, input_data, output_data
type TestCase struct {
	Tid        int    `gorm:"comment:测试样例ID;primaryKey;autoIncrement"`
	Pid        int    `gorm:"comment:题目ID;not null"`
	InputData  string `gorm:"comment:输入数据;not null"`
	OutputData string `gorm:"comment:输出数据;not null"`
}

// 竞赛表：contest_id, title, subtitle, difficulty, password, scored, have_password, is_visible, status, announcement, start_at, end_at
type Competition struct {
	ContestID    int       `gorm:"comment:比赛ID;primaryKey;autoIncrement"`
	Title        string    `gorm:"comment:标题;not null"`
	Subtitle     string    `gorm:"comment:副标题;not null"`
	Difficulty   string    `gorm:"comment:难度;not null"`
	Password     string    `gorm:"comment:密码;"`
	Scored       bool      `gorm:"comment:是否已计分;not null"`
	HavePassword bool      `gorm:"comment:是否有密码;not null"`
	IsVisible    bool      `gorm:"comment:是否可见;not null"`
	Status       int       `gorm:"comment:竞赛状态，0=未开始，1=正在进行中，2=已结束;"`
	Announcement string    `gorm:"comment:公告;"`
	Start_at     time.Time `gorm:"comment:开始时间;not null"`
	End_at       time.Time `gorm:"comment:结束时间;not null"`
}

// 用户竞赛关联表：contest_id, uid, username, join_date, score
type UserCompetitions struct {
	ContestID int       `gorm:"comment:比赛ID;not null"`
	Uid       int       `gorm:"comment:用户ID;not null"`
	Username  string    `gorm:"comment:用户名;not null"`
	Join_date time.Time `gorm:"comment:加入时间;not null"`
	Score     int       `gorm:"comment:用户在该竞赛获取的分数;"`
}
