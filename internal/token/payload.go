package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
	ErrKeySize      = errors.New("invalid key size: must be 32 byte")
)

type Payload struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	paseto.JSONToken
}

func NewPayload(username string, duration time.Duration) *Payload {
	tokenId := uuid.New().String()

	payload := &Payload{
		ID:       tokenId,
		Username: username,
		JSONToken: paseto.JSONToken{
			IssuedAt:   time.Now(),
			Expiration: time.Now().Add(duration),
			NotBefore:  time.Now(),
		},
	}

	return payload
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.Expiration) {
		return ErrExpiredToken
	}
	return nil
}
