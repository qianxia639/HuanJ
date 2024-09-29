package service

import (
	"github.com/jmoiron/sqlx"
)

type Server struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Server {
	return &Server{db: db}
}
