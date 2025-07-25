// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: friendship.sql

package db

import (
	"context"
)

const createFriendship = `-- name: CreateFriendship :one
INSERT INTO friendships (
    user_id, friend_id, remark
) VALUES (
    $1, $2, $3
)
RETURNING user_id, friend_id, remark, created_at, updated_at
`

type CreateFriendshipParams struct {
	UserID   int32  `json:"user_id"`
	FriendID int32  `json:"friend_id"`
	Remark   string `json:"remark"`
}

func (q *Queries) CreateFriendship(ctx context.Context, arg *CreateFriendshipParams) (Friendship, error) {
	row := q.db.QueryRow(ctx, createFriendship, arg.UserID, arg.FriendID, arg.Remark)
	var i Friendship
	err := row.Scan(
		&i.UserID,
		&i.FriendID,
		&i.Remark,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteFriend = `-- name: DeleteFriend :exec
DELETE FROM friendships 
WHERE (user_id, friend_id) 
IN (($1, $2), ($2, $1))
`

type DeleteFriendParams struct {
	UserID   int32 `json:"user_id"`
	UserID_2 int32 `json:"user_id_2"`
}

// DELETE FROM friendships
// WHERE (user_id = $1 AND friend_id = $2)
//
//	OR (user_id = $2 AND friend_id = $1);
func (q *Queries) DeleteFriend(ctx context.Context, arg *DeleteFriendParams) error {
	_, err := q.db.Exec(ctx, deleteFriend, arg.UserID, arg.UserID_2)
	return err
}

const existsFriendship = `-- name: ExistsFriendship :one
SELECT EXISTS(
    SELECT 1 FROM friendships
    WHERE (user_id = $1 AND friend_id = $2)
)
`

type ExistsFriendshipParams struct {
	UserID   int32 `json:"user_id"`
	FriendID int32 `json:"friend_id"`
}

// 是否已是好友关系
func (q *Queries) ExistsFriendship(ctx context.Context, arg *ExistsFriendshipParams) (bool, error) {
	row := q.db.QueryRow(ctx, existsFriendship, arg.UserID, arg.FriendID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getFriendList = `-- name: GetFriendList :many
SELECT u.id, u.username, u.nickname, u.password, u.email, u.gender, u.brithday, u.avatar_url, u.signature, u.password_changed_at, u.created_at, u.updated_at
FROM friendships f
JOIN users u ON f.friend_id = u.id
WHERE f.user_id = $1
`

func (q *Queries) GetFriendList(ctx context.Context, userID int32) ([]User, error) {
	rows, err := q.db.Query(ctx, getFriendList, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Nickname,
			&i.Password,
			&i.Email,
			&i.Gender,
			&i.Brithday,
			&i.AvatarUrl,
			&i.Signature,
			&i.PasswordChangedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
