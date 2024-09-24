package model

import "time"

type User struct {
	// 用户Id
	ID int32 `json:"id,omitempty"`
	// 用户名
	Username string `json:"username,omitempty"`
	// 用户昵称
	Nickname string `json:"nickname,omitempty"`
	// 密码
	Password string `json:"-"`
	// 邮箱
	Email string `json:"email,omitempty"`
	// 性别, 1 男, 2 女, 3 未知
	Gender int8 `json:"gender,omitempty"`
	// 头像图片路径或链接
	ProfilePictureUrl string `json:"profile_picture_url,omitempty"`
	// 在线状态(在线/离线)
	OnlineStatus bool `json:"online_status"`
	// 密码更新时间
	PasswordChangedAt time.Time `json:"password_changed_at,omitempty"`
	// 最后在线时间
	LastLoginAt time.Time `json:"last_login_at,omitempty"`
	// 创建时间
	CreatedAt time.Time `json:"created_at,omitempty"`
	// 更新时间
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
