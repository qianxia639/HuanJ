package service

import (
	"github.com/jmoiron/sqlx"
)

type Queries struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Queries {
	return &Queries{db: db}
}
