package model

import (
	"encoding/json"
	"time"
)

type GroupMember struct {
	GroupId  int32     `db:"group_id" json:"group_id,omitempty"`   // 群组ID
	UserId   int32     `db:"user_id" json:"user_id,omitempty"`     // 用户ID
	Role     int8      `db:"role" json:"role,omitempty"`           // 成员角色
	Waiting  bool      `db:"waiting" json:"waiting,omitempty"`     // 等待通过
	JoinedAt time.Time `db:"joined_at" json:"joined_at,omitempty"` // 加入时间
}

func (gm *GroupMember) MarshalBinary() ([]byte, error) {
	return json.Marshal(gm)
}

func (gm *GroupMember) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, gm)
}
