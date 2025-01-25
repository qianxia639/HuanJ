package model

import (
	"encoding/json"
	"time"
)

type User struct {
	ID                uint32    `db:"id" json:"id,omitempty"`                                   // 用户Id
	Username          string    `db:"username" json:"username,omitempty"`                       // 用户名
	Nickname          string    `db:"nickname" json:"nickname,omitempty"`                       // 用户昵称
	Password          string    `db:"password" json:"-"`                                        // 用户密码
	Email             string    `db:"email" json:"email,omitempty"`                             // 用户邮箱
	Gender            int8      `db:"gender" json:"gender,omitempty"`                           // 用户性别, 1: 男, 2: 女, 3: 未知
	Avatar            string    `db:"avatar" json:"avatar,omitempty"`                           // 用户头像
	PasswordChangedAt time.Time `db:"password_changed_at" json:"password_changed_at,omitempty"` // 上次密码更新时间
	CreatedAt         time.Time `db:"created_at" json:"created_at,omitempty"`                   // 创建时间
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at,omitempty"`                   // 更新时间
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
