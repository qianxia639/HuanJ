package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Queries struct {
	db *sqlx.DB
}

func NewQueries(db *sqlx.DB) *Queries {
	return &Queries{db: db}
}

func (q *Queries) execTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := q.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("事务回滚失败: %v, 原始错误: %v", rbErr, err)
		}
		return err
	}
	return tx.Commit()
}
