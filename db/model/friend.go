package model

import (
	"encoding/json"
	"time"
)

type Friend struct {
	ID        uint32    `db:"id" json:"id,omitempty"`                 // 好友关系标识
	UserId    uint32    `db:"user_id" json:"user_id,omitempty"`       // 用户ID
	FriendId  uint32    `db:"friend_id" json:"friend_id,omitempty"`   // 好友的用户ID
	Status    uint8     `db:"status" json:"status,omitempty"`         // 好友关系状态, 1: 待确认, 2: 已确认
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"` // 关系创建时间
	UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"` // 更新时间
}

func (f *Friend) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

func (f *Friend) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, f)
}
