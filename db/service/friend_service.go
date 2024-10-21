package db

import (
	"Dandelion/db/model"
	"context"
)

func (q *Queries) AddFriend(ctx context.Context, userId, friendId uint32) error {

	sql := `INSERT INTO friends (user_id, friend_id) VALUES ($1, $2)`

	row := q.db.QueryRowContext(ctx, sql, userId, friendId)

	return row.Err()

}

func (q *Queries) GetFriend(ctx context.Context, userId, friendId uint32) (*model.Friend, error) {

	sql := `SELECT * FROM friends WHERE user_id = $1 AND friend_id = $2 AND status != 3`

	var friend model.Friend
	err := q.db.GetContext(ctx, &friend, sql, userId, friendId)

	return &friend, err

}

func (q *Queries) GetFriendAll(ctx context.Context, userId uint32) ([]model.Friend, error) {

	sql := `SELECT * FROM friends WHERE user_id = $1 AND status = 1`

	friends := []model.Friend{}
	err := q.db.SelectContext(ctx, &friends, sql, userId)

	return friends, err
}
