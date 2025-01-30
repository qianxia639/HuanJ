package db

import (
	"Ice/db/model"
	"Ice/internal/logs"
	"context"
)

func (q *Queries) AddFriendRequest(ctx context.Context, fromUserId, toUserId uint32, requestDesc string) error {

	sql := `INSERT INTO friend_requests (from_user_id, to_user_id, request_desc) VALUES ($1, $2, $3)`

	return q.db.QueryRowContext(ctx, sql, fromUserId, toUserId, requestDesc).Err()

}

func (q *Queries) ExistsFriendRecord(ctx context.Context, fromUserId, toUserId uint32) int8 {

	sql := `SELECT COUNT(*) FROM friend_requests WHERE (from_user_id = $1 AND to_user_id = $2 AND status = 1)`

	var count int8
	err := q.db.GetContext(ctx, &count, sql, fromUserId, toUserId)
	if err != nil {
		logs.Error(err)
	}

	return count
}

func (q *Queries) ExistsFriendRequest(ctx context.Context, requestId uint32) int8 {

	sql := `SELECT COUNT(*) FROM friend_requests WHERE (id = $1 AND status = 1)`

	var count int8
	err := q.db.GetContext(ctx, &count, sql, requestId)
	if err != nil {
		logs.Error(err)
	}

	return count
}

func (q *Queries) GetFriendRequest(ctx context.Context, requestId, status uint32) (*model.FriendRequest, error) {
	sql := `SELECT * FROM friend_requests WHERE id = $1 AND status = $2 LIMIT 1`

	var fr model.FriendRequest
	err := q.db.GetContext(ctx, &fr, sql, requestId, status)
	if err != nil {
		return nil, err
	}

	return &fr, nil
}

func (q *Queries) InsertAcceptFriendRequestTx(ctx context.Context, requestId, fromUserId, toUserId uint32) error {

	sql1 := `UPDATE friend_requests SET status = 2, changed_at = now() WHERE id = $1 AND status = 1`
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

	if _, err := tx.ExecContext(ctx, sql2, fromUserId, toUserId); err != nil {
		logs.Errorf("accept friend: sql2 error: %v", err.Error())
		return tx.Rollback()
	}

	if _, err := tx.ExecContext(ctx, sql3, toUserId, fromUserId); err != nil {
		logs.Errorf("accept friend: sql3 error: %v", err.Error())
		return tx.Rollback()
	}

	return tx.Commit()

}

func (q *Queries) UpdateFriendRequest(ctx context.Context, requestId uint32) error {

	sql := `UPDATE friend_requests SET status = 4, changed_at = now() WHERE id = $1 AND status = 1`

	_, err := q.db.ExecContext(ctx, sql, requestId)
	if err != nil {
		logs.Error(err)
		return err
	}

	return nil
}
