package model

import (
	"encoding/json"
	"time"
)

type Friendship struct {
	UserId     uint32    `db:"user_id" json:"user_id,omitempty"`         // 用户ID
	FriendId   uint32    `db:"friend_id" json:"friend_id,omitempty"`     // 好友的用户ID
	AcceptedAt time.Time `db:"accepted_at" json:"accepted_at,omitempty"` // 更新时间
}

func (f *Friendship) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

func (f *Friendship) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, f)
}
