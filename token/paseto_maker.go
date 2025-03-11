package token

import (
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto        *paseto.V2
	sysmmetricKey []byte
}

func NewPasetoMaker(sysmmetricKey string) Maker {
	if len(sysmmetricKey) != 32 {
		panic(ErrKeySize)
	}

	maker := &PasetoMaker{
		paseto:        paseto.NewV2(),
		sysmmetricKey: []byte(sysmmetricKey),
	}

	return maker
}

// 创建Token
func (maker *PasetoMaker) CreateToken(args Token) (string, error) {
	payload := NewPayload(args.Username, args.Duration)

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
