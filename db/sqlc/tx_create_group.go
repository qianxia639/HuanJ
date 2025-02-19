package db

import (
	"context"
)

type CreateGroupTxParams struct {
	CreateGroupParams
	AfterCreate func(group Group) error
	UserId      int32
	Role        int16
	Agreed      bool
}

type CreateGroupTxResult struct {
	Group       Group       `json:"group"`
	GroupMember GroupMember `json:"group_member"`
}

func (store *SQLStore) CreateGroupTx(ctx context.Context, args CreateGroupTxParams) (CreateGroupTxResult, error) {
	var result CreateGroupTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Group, err = q.CreateGroup(ctx, &args.CreateGroupParams)
		if err != nil {
			return err
		}

		result.GroupMember, err = q.CreateGroupMember(ctx, &CreateGroupMemberParams{
			GroupID: result.Group.ID,
			UserID:  args.UserId,
			Role:    args.Role,
			Agreed:  args.Agreed,
		})

		return err
	})

	return result, err
}
