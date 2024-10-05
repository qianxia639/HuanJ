package models

import "time"

type User struct {
	ID       uint32 `json:"id,omitempty" db:"id"`
	Username string `json:"username,omitempty" db:"username"`
	Nickname string `json:"nickname,omitempty" db:"nickname"`
	Password string `json:"-" db:"password"`
	Salt     string `json:"-" db:"salt"`
	Email    string `json:"email,omitempty" db:"email"`
	// 性别, 1 男, 2 女, 3 未知
	Gender            int8      `json:"gender,omitempty" db:"gender"`
	IsOnline          bool      `json:"is_online"  db:"is_online"`
	Avatar            string    `json:"avatar,omitempty" db:"avatar"`
	PasswordChangedAt time.Time `json:"password_changed_at,omitempty" db:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
