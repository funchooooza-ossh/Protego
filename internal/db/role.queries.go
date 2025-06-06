package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

func (q *UQueries) GetRoleByID(ctx context.Context, id pgtype.UUID) (Role, error) {
	row := q.db.QueryRow(ctx, getRoleByID, id)
	var i Role
	err := row.Scan(&i.ID, &i.Code)
	return i, err
}

func (q *UQueries) GetRoleByCode(ctx context.Context, code string) (Role, error) {
	row := q.db.QueryRow(ctx, getRoleByCode, code)
	var i Role
	err := row.Scan(&i.ID, &i.Code)
	return i, err
}

type RoleParams struct {
	ID   pgtype.UUID `json:"id"`
	Code string      `json:"code"`
}

func (q *UQueries) CreateRole(ctx context.Context, arg RoleParams) (Role, error) {
	row := q.db.QueryRow(ctx, createRole, arg.ID, arg.Code)
	var i Role
	err := row.Scan(&i.ID, &i.Code)
	return i, err
}
