package model

import "time"

type InvitationCode struct {
	ID        int       `json:"id"`
	Code      string    `json:"code"`
	UserId    uint32    `json:"user_id"`
	Status    string    `json:"status"`
	UsedAt    time.Time `json:"used_at"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}
