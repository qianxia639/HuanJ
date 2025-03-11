package token

import "time"

type Token struct {
	Username string
	Duration time.Duration
}

type Maker interface {
	// 创建Token
	CreateToken(token Token) (string, error)

	// 校验Token
	VerifyToken(token string) (*Payload, error)
}

var _ Maker = (*PasetoMaker)(nil)
var _ Maker = (*PasetoMakerV2)(nil)
