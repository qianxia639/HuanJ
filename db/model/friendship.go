package model

import (
	"encoding/json"
	"time"
)

type Friendship struct {
	UserId    int32     `db:"user_id" json:"user_id,omitempty"`        // 用户ID
	FriendId  int32     `db:"friend_id" json:"friend_id,omitempty"`    // 好友的用户ID
	CreatedAt time.Time `db:"created_at" json:"accepted_at,omitempty"` // 创建时间
}

func (f *Friendship) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

func (f *Friendship) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, f)
}
