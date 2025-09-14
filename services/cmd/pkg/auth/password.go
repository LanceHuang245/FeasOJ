package auth

import "golang.org/x/crypto/bcrypt"

// 密码加密
func EncryptPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// 密码验证
func VerifyPassword(password, encryptedPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPwd), []byte(password))
	return err == nil
}
