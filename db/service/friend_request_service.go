package service

import (
	"Ice/internal/logs"
	"context"
)

func (q *Queries) AddFriendRequest(ctx context.Context, fromUserId, toUserId int32, requestDesc string) error {

	sql := `INSERT INTO friend_requests (user_id, friend_id, request_desc) VALUES ($1, $2, $3)`

	return q.db.QueryRowContext(ctx, sql, fromUserId, toUserId, requestDesc).Err()

}

func (q *Queries) ExistsFriendRequest(ctx context.Context, fromUserId, toUserId int32) int8 {

	sql := `SELECT COUNT(*) FROM friend_requests WHERE 
		((user_id = $1 AND friend_id = $2) OR 
		(user_id = $2 AND friend_id = $1)) AND status = 1`

	var count int8
	err := q.db.GetContext(ctx, &count, sql, fromUserId, toUserId)
	if err != nil {
		logs.Error(err)
	}

	return count
}

func (q *Queries) AcceptFriendRequest(ctx context.Context, requestId, userId int32) error {

	sql1 := `UPDATE friend_requests SET status = 2, updated_at = now() WHERE friend_id = $1 AND status = 1`
	sql2 := `INSERT INTO friendships (user_id, friend_id) VALUES ($1, $2)`
	sql3 := `INSERT INTO friendships (user_id, friend_id) VALUES ($1, $2)`

	tx, err := q.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, sql1, requestId); err != nil {
		logs.Errorf("accept friend: sql1 error: %v", err.Error())
		return tx.Rollback()
	}

	if _, err := tx.ExecContext(ctx, sql2, userId, requestId); err != nil {
		logs.Errorf("accept friend: sql2 error: %v", err.Error())
		return tx.Rollback()
	}

	if _, err := tx.ExecContext(ctx, sql3, requestId, userId); err != nil {
		logs.Errorf("accept friend: sql3 error: %v", err.Error())
		return tx.Rollback()
	}

	// TODO: 成功后异步发送通知

	return tx.Commit()

}

func (q *Queries) RejectFriendRequest(ctx context.Context, requestId, userId int32) error {

	sql := `
		UPDATE friend_requests
		SET
			status  = 3,
			updated_at = now()
		WHERE
			user_id = $1 AND friend_id = $2 AND status = 1
	`
	_, err := q.db.ExecContext(ctx, sql, userId, requestId)
	if err != nil {
		logs.Error(err)
		return err
	}

	return nil
}
