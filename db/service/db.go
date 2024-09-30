package db

import (
	"github.com/jmoiron/sqlx"
)

type Queries struct {
	db *sqlx.DB
}

func NewQueries(db *sqlx.DB) *Queries {
	return &Queries{db: db}
}
