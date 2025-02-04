package model

import "time"

type InvitationCode struct {
	ID        int       `json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	UserId    int32     `json:"user_id" db:"user_id"`
	Status    int       `json:"status" db:"status"`
	UsedAt    time.Time `json:"used_at" db:"used_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ExpiredAt time.Time `json:"expired_at" db:"expired_at"`
}
