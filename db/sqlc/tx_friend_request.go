package db

import "context"

type FriendRequestTxParams struct {
	Status       int8   `json:"status"`
	UserId       int32  `json:"user_id"`
	FriendId     int32  `json:"friend_id"`
	FromNickname string `json:"from_nickname"`
	ToNickname   string `json:"to_nickname"`
}

func (store *SQLStore) FriendRequestTx(ctx context.Context, args FriendRequestTxParams) error {

	err := store.execTx(ctx, func(q *Queries) error {
		err := q.UpdateFriendRequest(ctx, &UpdateFriendRequestParams{
			UserID:   args.UserId,
			FriendID: args.FriendId,
			Status:   args.Status,
		})
		if err != nil {
			return err
		}

		if args.UserId < args.FriendId {
			err = addFriendship(ctx, q, addFriendshipParams{
				UserId:      args.UserId,
				FriendId:    args.FriendId,
				FromComment: args.ToNickname,
				ToComment:   args.FromNickname,
			})
		} else {
			err = addFriendship(ctx, q, addFriendshipParams{
				UserId:      args.FriendId,
				FriendId:    args.UserId,
				FromComment: args.FromNickname,
				ToComment:   args.ToNickname,
			})
		}

		return err
	})
	return err
}

type addFriendshipParams struct {
	UserId      int32
	FriendId    int32
	FromComment string
	ToComment   string
}

func addFriendship(ctx context.Context, q *Queries, args addFriendshipParams) error {
	_, err := q.CreateFriendship(ctx, &CreateFriendshipParams{
		UserID:   args.UserId,
		FriendID: args.FriendId,
		Comment:  args.FromComment,
	})
	if err != nil {
		return err
	}

	_, err = q.CreateFriendship(ctx, &CreateFriendshipParams{
		UserID:   args.FriendId,
		FriendID: args.UserId,
		Comment:  args.ToComment,
	})

	return err
}
