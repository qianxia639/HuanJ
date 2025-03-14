package db

import "context"

type FriendRequestTxParams struct {
	Status     int8   `json:"status"`
	SenderId   int32  `json:"sender_id"`
	ReceiverId int32  `json:"receiver_id"`
	FromNote   string `json:"from_note"`
	ToNote     string `json:"to_note"`
}

func (store *SQLStore) FriendRequestTx(ctx context.Context, args FriendRequestTxParams) error {

	err := store.execTx(ctx, func(q *Queries) error {
		// 更新好友请求状态
		err := q.UpdateFriendRequest(ctx, &UpdateFriendRequestParams{
			SenderID:   args.SenderId,
			ReceiverID: args.ReceiverId,
			Status:     args.Status,
		})
		if err != nil {
			return err
		}

		// 创建双向好友关系(批量插入)
		return createMutualFriendships(ctx, q, createMutualFriendshipParams{
			SenderId:   args.SenderId,
			ReceiverId: args.ReceiverId,
			FromNote:   args.FromNote,
			ToNote:     args.ToNote,
		})
	})
	return err
}

const createMutilFriendships = `
	INSERT INTO friendships (
		user_id, friend_id, note
	) VALUES 
		($1, $2, $3),	-- 发送发 -> 接收方
		($2, $1, $4)	-- 接收方 -> 发送方
	ON CONFLICT(user_id, friend_id) DO NOTHING	-- ON CONFLICT 处理约束冲突
`

type createMutualFriendshipParams struct {
	SenderId   int32  `json:"sender_id"`
	ReceiverId int32  `json:"receiver_id"`
	FromNote   string `json:"from_note"`
	ToNote     string `json:"to_note"`
}

// 创建双向好友关系(使用批量操作)
func createMutualFriendships(ctx context.Context, q *Queries, args createMutualFriendshipParams) error {
	_, err := q.db.Exec(ctx, createMutilFriendships, args.SenderId, args.ReceiverId, args.FromNote, args.ToNote)
	return err
}
