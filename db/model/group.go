package model

import (
	"encoding/json"
	"time"
)

type Group struct {
	ID             int32     `db:"id" json:"id,omitempty"`                             // 群组ID
	GroupName      string    `db:"group_name" json:"group_name,omitempty"`             // 群组名称
	CreatorId      int32     `db:"creator_id" json:"creator_id,omitempty"`             // 群组创建者ID
	GroupAvatarUrl string    `db:"group_avatar_url" json:"group_avatar_url,omitempty"` // 群组头像URL
	Description    string    `db:"description" json:"description,omitempty"`           // 群组描述
	MaxMember      uint32    `db:"max_member" json:"max_member,omitempty"`             // 最大成员数
	CreatedAt      time.Time `db:"created_at" json:"created_at,omitempty"`             // 创建时间
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at,omitempty"`             // 群组信息更新时间
}

func (g *Group) MarshalBinary() ([]byte, error) {
	return json.Marshal(g)
}

func (g *Group) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, g)
}
