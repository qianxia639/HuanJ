package service

import (
	"Ice/db/model"
	"Ice/internal/logs"
	"context"

	"github.com/jmoiron/sqlx"
)

func (q *Queries) GetGroup(ctx context.Context, groupName string) (*model.Group, error) {

	sql := `SELECT * FROM groups WHERE group_name = $1 LIMIT 1`

	var group model.Group
	err := q.db.GetContext(ctx, &group, sql, groupName)
	if err != nil {
		logs.Errorf("GetGroup Error: %v\n", err.Error())
		// return nil, err
	}

	return &group, nil
}

type CreateGroupParams struct {
	Role        int8   `json:"role"`
	UserId      int32  `json:"user_id"`
	CreatorId   int32  `json:"creator_id"`
	GroupName   string `json:"group_name"`
	Description string `json:"description"`
}

func (q *Queries) CreateGroup(ctx context.Context, args *CreateGroupParams) error {

	sql := `INSERT INTO groups (
		group_name, creator_id, description
	) VALUES (
		$1, $2, $3
	) RETURNING id`

	sql2 := `INSERT INTO group_members (
		group_id, user_id, role, waiting
	) VALUES (
		$1, $2, $3, $4
	)`

	err := q.execTx(ctx, func(tx *sqlx.Tx) error {
		var group model.Group
		err := tx.QueryRowxContext(ctx, sql,
			args.GroupName,
			args.CreatorId,
			args.Description,
		).StructScan(&group)
		if err != nil {
			return err
		}

		err = tx.QueryRowxContext(ctx, sql2,
			group.ID,
			args.UserId,
			args.Role,
			false,
		).Err()
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
