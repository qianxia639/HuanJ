package utils

import (
	"strings"
)

// 手机号脱敏
func MaskPhone(phone string) string {
	if len(phone) < 8 {
		return phone
	}

	return phone[0:3] + "****" + phone[len(phone)-4:]

}

// 邮箱脱敏
func MaskEmail(email string) string {
	index := strings.LastIndex(email, "@")
	if index <= 1 {
		return email
	}

	username := email[:index]
	lenght := len(username)

	switch lenght {
	case 2:
		username = username[:1] + "****" + username[lenght:]
	case 3:
		username = username[:2] + "****" + username[lenght:]
	default:
		username = username[:3] + "****" + username[lenght:]
	}

	return username + email[index:]
}
