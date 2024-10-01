package models

import "time"

type User struct {
	ID       uint32 `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Password string `json:"-"`
	Salt     string `json:"-"`
	Email    string `json:"email,omitempty"`
	// 性别, 1 男, 2 女, 3 未知
	Gender            int8      `json:"gender,omitempty"`
	IsOnline          bool      `json:"is_online"`
	ProfilePictureUrl string    `json:"profile_picture_url,omitempty"`
	PasswordChangedAt time.Time `json:"password_changed_at,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}
