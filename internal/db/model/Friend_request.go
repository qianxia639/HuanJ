package model

import (
	"encoding/json"
	"time"
)

type FriendRequest struct {
	Id          uint32    `db:"id" json:"id,omitempty"`                     // 好友请求Id
	SenderId    uint32    `db:"sender_id" json:"sender_id,omitempty"`       // 请求者Id
	ReceiverId  uint32    `db:"receiver_id" json:"receiver_id,omitempty"`   // 接收者Id
	Status      bool      `db:"status" json:"status"`                       // 请求状态, 1: 待处理, 2: 已添加, 3: 已过期
	RequestDesc string    `db:"request_desc" json:"request_desc,omitempty"` // 请求信息
	RequestedAt time.Time `db:"requested_at" json:"requested_at,omitempty"` // 请求时间
	ChangedAt   time.Time `db:"changed_at" json:"changed_at,omitempty"`     // 变更时间
}

func (fr *FriendRequest) MarshalBinary() ([]byte, error) {
	return json.Marshal(fr)
}

func (fr *FriendRequest) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, fr)
}
