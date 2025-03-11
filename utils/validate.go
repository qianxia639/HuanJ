package utils

import "regexp"

// 校验用户名
func ValidateUsername(username string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_?!@.]{4,20}$`).MatchString(username)
}

// 校验密码
func ValidatePassword(password string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_?!@]{6,20}$`).MatchString(password)
}
