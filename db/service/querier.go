package service

import (
	"context"
)

type Querier interface {
	CreateUser(ctx context.Context, args *CreateUserParams) error
	ExistsUsername(ctx context.Context, username string) int8
	ExistsNickname(ctx context.Context, nickname string) int8
}

var _ Querier = (*Queries)(nil)
