package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

func (q *UQueries) HasAccess(ctx context.Context, roleID pgtype.UUID, action string, resourceCode string) (bool, error) {
	row := q.db.QueryRow(ctx, hasAccess, roleID, action, resourceCode)
	var has_access bool
	err := row.Scan(&has_access)
	return has_access, err
}
