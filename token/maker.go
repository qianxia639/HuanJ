package token

import "time"

type Maker interface {
	// 创建Token
	CreateToken(username string, duration time.Duration) (string, error)

	// 校验Token
	VerifyToken(token string) (*Payload, error)
}
