package db

import (
	"Ice/db/model"
	"Ice/internal/logs"
	"context"
)

func (q *Queries) GetCode(ctx context.Context, code string) *model.InvitationCode {

	sql := `SELECT * FROM invitation_codes WHERE code = $1 LIMIT 1`
	var ic model.InvitationCode
	err := q.db.GetContext(ctx, &ic, sql, code)
	if err != nil {
		logs.Errorf("Get code error: %s\n", err.Error())
	}

	return &ic

}
