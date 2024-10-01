package token

import (
	"time"

	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto        *paseto.V2
	sysmmetricKey []byte
}

func NewPasetoMaker(sysmmetricKey string) (Maker, error) {
	if len(sysmmetricKey) != 32 {
		return nil, ErrKeySize
	}

	maker := &PasetoMaker{
		paseto:        paseto.NewV2(),
		sysmmetricKey: []byte(sysmmetricKey),
	}

	return maker, nil
}

// 创建Token
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	token, err := maker.paseto.Encrypt(maker.sysmmetricKey, payload, nil)

	return token, err
}

// 校验Token
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.sysmmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
