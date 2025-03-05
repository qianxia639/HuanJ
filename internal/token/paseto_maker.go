package token

import (
	"crypto/ed25519"
	"time"

	"github.com/o1egl/paseto"
)

const keySize = 32

type PasetoMaker struct {
	paseto        *paseto.V2
	sysmmetricKey []byte
	privateKey    ed25519.PrivateKey
	publicKey     ed25519.PublicKey
}

func NewPasetoMakerV2(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) Maker {
	maker := &PasetoMaker{
		paseto:     paseto.NewV2(),
		privateKey: privateKey,
		publicKey:  publicKey,
	}

	return maker
}

func NewPasetoMaker(sysmmetricKey string) Maker {
	if len(sysmmetricKey) != keySize {
		panic(ErrKeySize)
	}

	maker := &PasetoMaker{
		paseto:        paseto.NewV2(),
		sysmmetricKey: []byte(sysmmetricKey),
	}

	return maker
}

// 创建Token
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload := NewPayload(username, duration)

	// token, err := maker.paseto.Encrypt(maker.sysmmetricKey, payload, nil)
	token, err := maker.paseto.Encrypt(maker.privateKey, payload, nil)

	return token, err
}

// 校验Token
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	// err := maker.paseto.Decrypt(token, maker.sysmmetricKey, payload, nil)
	err := maker.paseto.Verify(token, maker.publicKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// err = payload.Valid()
	// if err != nil {
	// 	return nil, err
	// }
	return payload, nil
}
