package db

import (
	"Ice/db/model"
	"Ice/internal/logs"
	"context"
)

// func (q *Queries) AddFriendRecord(ctx context.Context, userId, friendId uint32) error {

// 	sql := `INSERT INTO friendships (user_id, friend_id) VALUES ($1, $2)`

// 	return q.db.QueryRowContext(ctx, sql, userId, friendId).Err()

// }

func (q *Queries) CheckFriendship(ctx context.Context, userId, friendId int32) int8 {

	sql := `SELECT * FROM friend_requests WHERE user_id = $1 AND friend_id = $2 AND status = 2`
	var count int8
	err := q.db.GetContext(ctx, &count, sql, userId, friendId)
	if err != nil {
		logs.Errorf("CheckFriendship: %v\n", err.Error())
	}

	return count
}

func (q *Queries) AddFriendTx(ctx context.Context, userId, friendId int32) error {
	tx, err := q.db.BeginTxx(ctx, nil)
	if err != nil {
		logs.Errorf("AddFriend: begin transaction error: %v", err.Error())
		return err
	}

	sql1 := `UPDATE friendships SET status = 2 WHERE user_id = $1 AND friend_id = $2`
	sql2 := `INSERT INTO friendships (user_id, friend_id, status) VALUES ($1, $2, $3)`

	if _, err := tx.ExecContext(ctx, sql1, friendId, userId); err != nil {
		logs.Errorf("AddFriend: sql1 error: %v", err.Error())
		return tx.Rollback()
	}
	if _, err := tx.ExecContext(ctx, sql2, userId, friendId, 2); err != nil {
		logs.Errorf("AddFriend: sql2 error: %v", err.Error())
		return tx.Rollback()
	}

	return tx.Commit()

}

func (q *Queries) ExistsFriend(ctx context.Context, userId, friendId int32, status int8) int8 {

	sql := `SELECT COUNT(*) FROM friendships WHERE user_id = $1 AND friend_id = $2`

	var count int8
	_ = q.db.GetContext(ctx, &count, sql, userId, friendId, status)

	return count
}

func (q *Queries) GetFriend(ctx context.Context, userId, friendId int32) (*model.Friendship, error) {

	sql := `SELECT * FROM friendships WHERE user_id = $1 AND friend_id = $2`

	var friend model.Friendship
	err := q.db.GetContext(ctx, &friend, sql, userId, friendId)

	return &friend, err

}

func (q *Queries) GetFriendAll(ctx context.Context, userId int32) ([]model.Friendship, error) {

	sql := `SELECT * FROM friendships WHERE user_id = $1`

	friends := []model.Friendship{}
	err := q.db.SelectContext(ctx, &friends, sql, userId)

	return friends, err
}

func (q *Queries) DeleteFriend(ctx context.Context, userId, friendId int32) error {

	sql := `DELETE FROM friendships WHERE (user_id = $1 AND friend_id = $2) 
								OR (user_id = $2 AND friend_id = $1)`

	_, err := q.db.ExecContext(ctx, sql, userId, friendId)

	return err
}
