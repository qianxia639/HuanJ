package db

import (
	"Ice/db/model"
	"Ice/internal/logs"
	"context"
)

type CreateUserParams struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Gender   int8   `json:"gender"`
}

func (q *Queries) CreateUser(ctx context.Context, args *CreateUserParams) error {

	sql := `
	INSERT INTO users (
		username, nickname, password, email, gender
	) VALUES (
		$1, $2, $3, $4, $5
	)`

	row := q.db.QueryRowContext(ctx, sql,
		args.Username,
		args.Nickname,
		args.Password,
		args.Email,
		args.Gender,
	)
	return row.Err()
}

func (q *Queries) ExistsUser(ctx context.Context, username string) int8 {
	sql := `SELECT COUNT(*) FROM users WHERE username = $1`

	var count int8
	_ = q.db.GetContext(ctx, &count, sql, username)

	return count
}

func (q *Queries) ExistsEmail(ctx context.Context, email string) int8 {
	sql := `SELECT COUNT(*) FROM users WHERE email = $2`

	var count int8
	_ = q.db.GetContext(ctx, &count, sql, email)

	return count
}

func (q *Queries) ExistsNickname(ctx context.Context, nickname string) int8 {
	sql := `SELECT COUNT(*) FROM users WHERE nickname = $1`

	var count int8
	_ = q.db.GetContext(ctx, &count, sql, nickname)

	return count
}

func (q *Queries) GetUser(ctx context.Context, username string) (u model.User, err error) {

	sql := `SELECT * FROM users WHERE username = $1 LIMIT 1`
	err = q.db.GetContext(ctx, &u, sql, username)
	if err != nil {
		logs.Errorf("Get user: %v\n", err.Error())
	}

	return
}

func (q *Queries) GetUserById(ctx context.Context, id int32) (u model.User, err error) {

	sql := `SELECT * FROM users WHERE id = $1 LIMIT 1`
	err = q.db.GetContext(ctx, &u, sql, id)

	return
}

// IS DISTINCT FROM：此语法用于比较两个值是否不同，即使其中一个值为 NULL，也能正确处理
// 通过添加 AND 子句，避免在值未更改时执行更新操作，从而减少不必要的数据库写入
func (q *Queries) UpdateUser(ctx context.Context, user model.User) error {

	sql := `UPDATE users 
			SET 
				gender = $1, 
				nickname = $2, 
				updated_at = now()
			WHERE id = $3
			AND (gender IS DISTINCT FROM $1 OR nickname IS DISTINCT FROM $2)`
	_, err := q.db.ExecContext(ctx, sql,
		user.Gender,
		user.Nickname,
		user.ID,
	)

	return err
}
