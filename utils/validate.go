package utils

import (
	"regexp"
	"unicode"
)

// 校验用户名
func ValidateUsername(username string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_?!@.]{4,30}$`).MatchString(username)
}

// 校验密码
func ValidatePassword(password string) bool {
	// 长度为 6-20
	if len(password) < 6 || len(password) > 30 {
		return false
	}

	var (
		hasLetter  = false
		hasDigit   = false
		hasSpecial = false
		invalid    = false
	)

	specialChars := map[rune]bool{
		'_': true,
		'?': true,
		'!': true,
		'@': true,
		'.': true,
	}

	for _, c := range password {
		switch {
		case unicode.IsLetter(c): // 字母
			hasLetter = true
		case unicode.IsDigit(c): // 数字
			hasDigit = true
		case specialChars[c]: // 特殊字符
			hasSpecial = true
		default:
			invalid = true
		}
	}

	if invalid {
		return false
	}

	// 是否包含其中两种及以上
	return (boolToInt(hasLetter)+boolToInt(hasDigit)+boolToInt(hasSpecial) >= 2)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
