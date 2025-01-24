package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
	ErrKeySize      = errors.New("invalid key size: must be 32 byte")
)

type Payload struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"` // 创建时间
	ExpiredAt time.Time `json:"expire_at"` // 过期时间
}

func NewPayload(username string, duration time.Duration) *Payload {
	tokenId := uuid.New().String()

	payload := &Payload{
		ID:        tokenId,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
