package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserParams struct {
	ID       pgtype.UUID `json:"id"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
	RoleID   pgtype.UUID `json:"role_id"`
	Blocked  bool        `json:"blocked"`
}

func (q *UQueries) CreateUser(ctx context.Context, arg UserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.Password,
		arg.RoleID,
		arg.Blocked,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.RoleID,
		&i.Blocked,
	)
	return i, err
}

func (q *UQueries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.RoleID,
		&i.Blocked,
	)
	if err != nil && err.Error() == "no rows in result set" {
		return i, fmt.Errorf("%w", sql.ErrNoRows)
	}
	return i, err
}

func (q *UQueries) GetUserByID(ctx context.Context, id pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.RoleID,
		&i.Blocked,
	)
	if err != nil && err.Error() == "no rows in result set" {
		return i, fmt.Errorf("%w", sql.ErrNoRows)
	}
	return i, err
}

func (q *UQueries) UpdateUser(ctx context.Context, arg UserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.ID,
		arg.Email,
		arg.Password,
		arg.RoleID,
		arg.Blocked,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.RoleID,
		&i.Blocked,
	)
	return i, err
}

func (q *UQueries) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}
