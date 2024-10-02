package db

import (
	"Dandelion/db/models"
	"context"
	"time"
)

type CreateUserParams struct {
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Password  string    `json:"password"`
	Salt      string    `json:"salt"`
	Email     string    `json:"email"`
	Gender    int8      `json:"gender"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) CreateUser(ctx context.Context, args *CreateUserParams) error {

	sql := `
	INSERT INTO users (
		username, nickname, password, salt, email, gender, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4,$5, $6, $7, $8
	)`

	row := q.db.QueryRowContext(ctx, sql,
		args.Username,
		args.Nickname,
		args.Password,
		args.Salt,
		args.Email,
		args.Gender,
		args.CreatedAt,
		args.UpdatedAt,
	)
	return row.Err()

}

func (q *Queries) ExistsUser(ctx context.Context, username, email string) int8 {
	sql := `SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2`

	var count int8
	_ = q.db.GetContext(ctx, &count, sql, username, email)

	return count
}

func (q *Queries) ExistsNickname(ctx context.Context, nickname string) int8 {
	sql := `SELECT COUNT(*) FROM users WHERE nickname = $1`

	var count int8
	_ = q.db.GetContext(ctx, &count, sql, nickname)

	return count
}

func (q *Queries) GetUser(ctx context.Context, username string) (u models.User, err error) {
	// sql := `SELECT
	// 	id, username, nickname, password, salt, email, gender, is_online, profile_picture_url, password_changed_at, created_at, updated_at
	//  FROM users WHERE username = $1 LIMIT 1`

	sql := `SELECT * FROM users WHERE username = $1 LIMIT 1`
	// TODO 问题尚未解决
	err = q.db.GetContext(ctx, &u, sql, username)
	// row := q.db.QueryRowContext(ctx, sql, username)

	// err = row.Scan(
	// 	&u.ID,
	// 	&u.Username,
	// 	&u.Nickname,
	// 	&u.Password,
	// 	&u.Salt,
	// 	&u.Email,
	// 	&u.Gender,
	// 	&u.IsOnline,
	// 	&u.ProfilePictureUrl,
	// 	&u.PasswordChangedAt,
	// 	&u.CreatedAt,
	// 	&u.UpdatedAt,
	// )

	return
}
