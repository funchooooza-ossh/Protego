package adapters

import (
	"context"
	"database/sql"
	"errors"
	"net"

	e "github.com/funchooooza-ossh/protego/internal/errors"
	m "github.com/funchooooza-ossh/protego/internal/metrics/lifespan"
	"github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
)

func ParseDBError(ctx context.Context, err error, origin string) error {
	const component = "database"

	if err == nil {
		return nil
	}

	// Not found
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) || err.Error() == "no rows in result set" {
		e.LogErr(ctx, origin, err, e.Info)
		return e.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) { //TODO DO NOT REPEAT log return
		switch pgErr.Code {
		case "23505":
			e.LogErr(ctx, origin, err, e.Info)
			return e.ErrAlreadyExists
		case "23503":
			e.LogErr(ctx, origin, err, e.Info)
			return e.ErrConflict
		case "23502":
			e.LogErr(ctx, origin, err, e.Info)
			return e.ErrInvalidInput
		case "23514":
			e.LogErr(ctx, origin, err, e.Info)
			return e.ErrValidationFailed
		default:
			e.LogErr(ctx, origin, err, e.Error)
			return e.ErrInternal
		}
	}

	// network
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		m.Inc(component, "unavailable")
		e.LogErr(ctx, origin, netErr, e.Error)
		return e.ErrServiceDown
	}

	// context
	switch {
	case errors.Is(err, context.DeadlineExceeded),
		errors.Is(err, context.Canceled):
		m.Inc(component, "timeout")
		e.LogErr(ctx, origin, err, e.Error)
		return e.ErrTimeout
	}

	e.LogErr(ctx, origin, err, e.Error)
	m.Inc(component, "unhandled")
	return e.ErrInternal
}
