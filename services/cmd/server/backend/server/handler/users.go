package handler

import (
	"FeasOJ/internal/config"
	"FeasOJ/internal/global"
	"FeasOJ/internal/utils"
	"FeasOJ/internal/utils/sql"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/go-redis/redis"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 自增登录计数器
func incrLoginCounter(rdb *redis.Client, ip string, success bool) {
	var counterKey, blockKey string
	if success {
		counterKey = fmt.Sprintf("loginSuccCount:%s", ip)
		blockKey = fmt.Sprintf("loginBlock:%s", ip)
	} else {
		counterKey = fmt.Sprintf("loginFailCount:%s", ip)
		blockKey = fmt.Sprintf("loginBlock:%s", ip)
	}
	cnt, _ := rdb.Incr(counterKey).Result()
	if cnt == 1 {
		rdb.Expire(counterKey, 30*time.Minute)
	}
	if cnt >= 5 {
		rdb.Set(blockKey, 1, 30*time.Minute)
	}
}

// 注册
func Register(c *gin.Context) {
	clientIP := c.ClientIP()
	rdb := utils.ConnectRedis()
	defer rdb.Close()

	// 封禁键
	blockKey := fmt.Sprintf("regBlock:%s", clientIP)
	if ok, _ := rdb.Exists(blockKey).Result(); ok == 1 {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": GetMessage(c, "rateLimit")})
		return
	}

	var req global.RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "invalidRequest")})
		return
	}

	// 判断用户或邮箱是否存在
	if sql.IsUserExist(req.Username, req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "userAlreadyinUse")})
		return
	}
	if config.GlobalConfig.Features.ProfanityDetectorEnabled {
		if utils.DetectText(req.Username) {
			c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "profanity")})
			return
		}
	}
	vcodeStatus := utils.CompareVerifyCode(req.Vcode, req.Email)
	if !vcodeStatus {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "captchaError")})
		return
	}
	success := sql.Register(req.Username, utils.EncryptPassword(req.Password), req.Email, uuid.New().String(), 0)

	// 计数键
	var counterKey string
	if success {
		counterKey = fmt.Sprintf("regSuccCount:%s", clientIP)
	} else {
		counterKey = fmt.Sprintf("regFailCount:%s", clientIP)
	}
	// 自增计数
	cnt, _ := rdb.Incr(counterKey).Result()
	// 首次自增时，设置 30 分钟过期
	if cnt == 1 {
		rdb.Expire(counterKey, 30*time.Minute)
	}
	// 达到 5 次阈值，开启 30 分钟封禁
	if cnt >= 5 {
		rdb.Set(blockKey, 1, 30*time.Minute)
	}

	// 返回结果
	if success {
		c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
	}
}

// 登录
func Login(c *gin.Context) {
	clientIP := c.ClientIP()
	rdb := utils.ConnectRedis()
	defer rdb.Close()

	// 封禁键
	blockKey := fmt.Sprintf("loginBlock:%s", clientIP)
	if ok, _ := rdb.Exists(blockKey).Result(); ok == 1 {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": GetMessage(c, "rateLimit")})
		return
	}

	// 取参
	userParam := c.Query("username")
	password := c.Query("password")

	// 查用户名
	var username string
	if utils.IsEmail(userParam) {
		username = sql.SelectUserByEmail(userParam).Username
	} else {
		username = userParam
	}

	// 验证用户存在
	storedPwd := utils.SelectUser(username).Password
	if storedPwd == "" {
		// 统计失败
		incrLoginCounter(rdb, clientIP, false)
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}

	// 是否封号
	if sql.SelectUserInfo(username).IsBan {
		c.JSON(http.StatusForbidden, gin.H{"message": GetMessage(c, "userIsBanned")})
		return
	}

	// 验证密码
	ok := utils.VerifyPassword(password, storedPwd)
	// 统计 成功/失败
	incrLoginCounter(rdb, clientIP, ok)

	if ok {
		token, err := utils.GenerateToken(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "loginSuccess"), "token": token})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "wrongPassword")})
	}
}

// 获取验证码
// 获取验证码（IP + 邮箱 1 分钟限流）
func GetCaptcha(c *gin.Context) {
	clientIP := c.ClientIP()
	email := c.Query("email")
	isCreate := c.GetHeader("iscreate")

	rdb := utils.ConnectRedis()
	defer rdb.Close()

	// 限流键
	rateKey := fmt.Sprintf("captchaRateLimit:%s:%s", clientIP, email)
	if exists, _ := rdb.Exists(rateKey).Result(); exists == 1 {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": GetMessage(c, "rateLimit")})
		return
	}

	// 原有校验逻辑：注册/重置场景区分
	if isCreate == "false" {
		if sql.SelectUserByEmail(email).Username == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
			return
		}
	} else {
		if sql.SelectUserByEmail(email).Username != "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "userAlreadyinUse")})
			return
		}
	}

	// 发送验证码
	code := utils.GenerateVerifycode()
	sent := utils.SendVerifycode(config.GlobalConfig.Mail, email, code)
	if sent {
		// 限流键写入，1分钟过期
		rdb.Set(rateKey, 1, time.Minute)
		c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
	}
}

// 验证用户信息
func VerifyUserInfo(c *gin.Context) {
	var Username string
	user := c.GetHeader("Username")
	unescapeUsername, err := url.QueryUnescape(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}
	if utils.IsEmail(unescapeUsername) {
		Username = sql.SelectUserByEmail(unescapeUsername).Username
	} else {
		Username = unescapeUsername
	}
	// 查询对应的用户信息
	userInfo := sql.SelectUserInfo(Username)
	if userInfo.Username == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
	} else {
		c.JSON(http.StatusOK, gin.H{"info": userInfo})
	}
}

// 获取用户信息
func GetUserInfo(c *gin.Context) {
	// 获取用户名
	username := c.Param("username")
	// 查询对应的用户信息
	userInfo := sql.SelectUserInfo(username)
	if userInfo.Username == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
	} else {
		c.JSON(http.StatusOK, gin.H{"info": userInfo})
	}
}

// 更新密码
func UpdatePassword(c *gin.Context) {
	var req global.UpdatePasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "invalidRequest")})
		return
	}

	// 检查邮箱格式
	if !utils.IsEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "invalidEmail")})
		return
	}

	// 检查邮箱是否存在
	if sql.SelectUserByEmail(req.Email).Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "emailNotExist")})
		return
	}

	// 验证验证码
	vcodeStatus := utils.CompareVerifyCode(req.Vcode, req.Email)
	if !vcodeStatus {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "captchaError")})
		return
	}

	// 更新密码
	success := sql.UpdatePassword(req.Email, utils.EncryptPassword(req.NewPassword))
	if success {
		c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
	}
}

// 更新个人简介
func UpdateSynopsis(c *gin.Context) {
	synopsis := c.PostForm("synopsis")
	encodedUsername := c.GetHeader("Username")
	username, _ := url.QueryUnescape(encodedUsername)
	// 更新简介
	if config.GlobalConfig.Features.ImageGuardEnabled {
		if utils.DetectText(synopsis) {
			c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "profanity")})
			return
		}
	}
	if sql.UpdateSynopsis(username, synopsis) {
		c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
	}
}

// 上传头像
func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "invalidrequest")})
		return
	}
	encodedUsername := c.GetHeader("Username")
	username, _ := url.QueryUnescape(encodedUsername)
	// 获取用户信息
	userInfo := sql.SelectUserInfo(username)
	if userInfo.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "invalidrequest")})
		return
	}
	newFileName := fmt.Sprintf("%d%s", userInfo.Uid, path.Ext(file.Filename))
	tempFilePath := filepath.Join(global.AvatarsDir, newFileName)
	compressedFilePath := filepath.Join(global.AvatarsDir, fmt.Sprintf("%d%s", userInfo.Uid, path.Ext(file.Filename)))
	if _, err := os.Stat(compressedFilePath); err == nil {
		os.Remove(compressedFilePath)
	}
	// 保存临时文件
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}
	// 检测图像，若违规则删除
	if config.GlobalConfig.Features.ImageGuardEnabled {
		if !utils.PredictImage(tempFilePath) {
			c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "illegalImage")})

			// 删除图像
			if err := os.Remove(tempFilePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
				return
			}
			return
		}
	}

	// 压缩图像
	if err := utils.CompressImage(tempFilePath, compressedFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}

	// 上传压缩后的头像路径至数据库
	if !sql.UpdateAvatar(username, newFileName) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 获取排行榜
func GetRanking(c *gin.Context) {
	cacheKey := "ranking"
	// 从缓存中获取数据
	var ranking []global.UserInfoRequest
	err := utils.GetCache(cacheKey, &ranking)
	if err != nil || len(ranking) == 0 {
		ranking = sql.SelectRank100Users()
		// 缓存5分钟
		err := utils.SetCache(cacheKey, ranking, 5*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"ranking": ranking})
}

// 发送验证码
func SendVerifycode(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "invalidRequest")})
		return
	}

	// 检查邮箱格式
	if !utils.IsEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "invalidEmail")})
		return
	}

	// 检查邮箱是否已被注册
	if sql.SelectUserByEmail(email).Username != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "emailAlreadyinUse")})
		return
	}

	// 生成验证码
	code := utils.GenerateVerifycode()

	// 发送验证码
	sent := utils.SendVerifycode(config.GlobalConfig.Mail, email, code)
	if sent {
		c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
	}
}
