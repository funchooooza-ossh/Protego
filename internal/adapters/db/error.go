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

func ParseDBError(err error, origin string) error {
	if err == nil {
		return nil
	}

	// Not found
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) || err.Error() == "no rows in result set" {
		e.LogErr(origin, err, e.Info)
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
			return e.ErrInternal
		}
	}

	// network
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		e.LogErr(origin, netErr, e.Warn)
		return e.ErrServiceDown
	}

	// context
	switch {
	case errors.Is(err, context.DeadlineExceeded),
		errors.Is(err, context.Canceled):
		e.LogErr(origin, err, e.Warn)
		return e.ErrTimeout
	}

	e.LogErr(origin, err, e.Warn)
	return e.ErrInternal
}
