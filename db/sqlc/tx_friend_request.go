package db

import "context"

type FriendRequestTxParams struct {
	Status     int8   `json:"status"`
	FromUserId int32  `json:"from_user_id"`
	ToUserId   int32  `json:"to_user_id"`
	FromNote   string `json:"from_note"`
	ToNote     string `json:"to_note"`
}

func (store *SQLStore) FriendRequestTx(ctx context.Context, args *FriendRequestTxParams) error {

	err := store.execTx(ctx, func(q *Queries) error {
		// 更新好友请求状态
		err := q.UpdateFriendRequest(ctx, &UpdateFriendRequestParams{
			FromUserID: args.FromUserId,
			ToUserID:   args.ToUserId,
			Status:     args.Status,
		})
		if err != nil {
			return err
		}

		// 确保fromUserId < toUserId来避免重复
		if args.FromUserId > args.ToUserId {
			args.FromUserId, args.ToUserId = args.ToUserId, args.FromUserId
		}

		// 创建双向好友关系(批量插入)
		return q.createMutualFriendships(ctx, createMutualFriendshipParams{
			FromUserId: args.FromUserId,
			ToUserId:   args.ToUserId,
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
	FromUserId int32  `json:"from_user_id"`
	ToUserId   int32  `json:"to_user_id"`
	FromNote   string `json:"from_note"`
	ToNote     string `json:"to_note"`
}

// 创建双向好友关系(使用批量操作)
func (q *Queries) createMutualFriendships(ctx context.Context, args createMutualFriendshipParams) error {
	_, err := q.db.Exec(ctx, createMutilFriendships, args.FromUserId, args.ToUserId, args.FromNote, args.ToNote)
	return err
}
