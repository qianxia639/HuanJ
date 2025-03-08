package token

import (
	"crypto/ed25519"
	"time"

	"github.com/o1egl/paseto"
)

type PasetoMakerV2 struct {
	paseto     *paseto.V2
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

func NewPasetoMakerV2(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) Maker {
	maker := &PasetoMakerV2{
		paseto:     paseto.NewV2(),
		privateKey: privateKey,
		publicKey:  publicKey,
	}

	return maker
}

// 创建Token
func (maker *PasetoMakerV2) CreateToken(username string, duration time.Duration) (string, error) {
	payload := NewPayload(username, duration)

	token, err := maker.paseto.Sign(maker.privateKey, payload, nil)

	return token, err
}

// 校验Token
func (maker *PasetoMakerV2) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Verify(token, maker.publicKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
