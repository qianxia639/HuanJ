package model

import "encoding/json"

type LoginUserInfo struct {
	User
	UserAgent string `json:"user_agent"` // 用户登录标识
	LoginIp   string `json:"login_ip"`   // 用户登录IP
}

func (l *LoginUserInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(l)
}

func (l *LoginUserInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, l)
}
