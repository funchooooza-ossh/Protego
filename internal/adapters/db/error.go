package adapters

import (
	"context"
	"database/sql"
	"errors"
	"net"

	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
)

func ParseDBError(ctx context.Context, err error, origin string) error {
	if err == nil {
		return nil
	}

	// Not found
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) || err.Error() == "no rows in result set" {
		e.LogErr(ctx, origin, err, e.Info)
		return e.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return e.ErrAlreadyExists
		case "23503":
			return e.ErrConflict
		case "23502":
			return e.ErrInvalidInput
		case "23514":
			return e.ErrValidationFailed
		default:
			e.LogErr(ctx, origin, err, e.Error)
			return e.ErrInternal
		}
	}

	// network
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		e.LogErr(ctx, origin, netErr, e.Error)
		return e.ErrServiceDown
	}

	// context
	switch {
	case errors.Is(err, context.DeadlineExceeded),
		errors.Is(err, context.Canceled):
		e.LogErr(ctx, origin, err, e.Error)
		return e.ErrTimeout
	}

	e.LogErr(ctx, origin, err, e.Error)
	return e.ErrInternal
}
