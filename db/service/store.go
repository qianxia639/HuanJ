package service

import "github.com/jmoiron/sqlx"

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
}

func NewStore(db *sqlx.DB) Store {
	return &SQLStore{
		Queries: New(db),
	}
}
