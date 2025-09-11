package utils

import (
	"FeasOJ/internal/config"
	"crypto/tls"
	"fmt"
	"math/rand"
	"time"

	"gopkg.in/gomail.v2"
)

// 随机生成4位数字验证码
func GenerateVerifycode() string {
	return fmt.Sprintf("%04d", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
}

// 发送验证码
func SendVerifycode(mailConfig config.MailConfig, to string, verifycode string) bool {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", mailConfig.User, "FeasOJ")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "FeasOJ验证码")
	m.SetBody("text/html", "<div style='text-align: center;'><h1><b>FeasOJ</b></h1><p>您正在进行敏感操作，<span style='color:red;'>5分钟后</span>失效</p><h1 style='letter-spacing: 10px;'><b>"+verifycode+"</b></h1></div>")
	d := gomail.NewDialer(mailConfig.Host, mailConfig.Port, mailConfig.User, mailConfig.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return false
	}
	// 将验证码同时存进Redis中等待校验
	rdb := ConnectRedis()
	err := rdb.Set(to, verifycode, 5*time.Minute).Err()
	return err == nil
}

// 测试发送
func TestSend(mailConfig config.MailConfig) bool {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", mailConfig.User, "FeasOJ")
	m.SetHeader("To", mailConfig.User)
	m.SetHeader("Subject", "FeasOJ邮件服务测试")
	m.SetBody("text/html", "<div style='text-align: center;'><h1><b>FeasOJ</b></h1><p>这是一封测试邮件，若看到该邮件意味着您的FeasOJ Email服务运行正常。</p></div>")
	d := gomail.NewDialer(mailConfig.Host, mailConfig.Port, mailConfig.User, mailConfig.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return false
	}
	return true
}

// 检验Redis中验证码与前端返回的是否相同
func CompareVerifyCode(frontendCode, to string) bool {
	// 通过邮箱来获取Redis中的验证码
	rdb := ConnectRedis()
	verifyCode, err := rdb.Get(to).Result()
	if err != nil {
		return false
	}
	if verifyCode == frontendCode {
		// 移除Redis中的验证码
		rdb.Del(to)
		return true
	}
	return false
}
