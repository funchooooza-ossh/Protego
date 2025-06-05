package adapters

import (
	"errors"
	"fmt"
	"net"

	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/jackc/pgconn"
	"golang.org/x/net/context"
)

func ParseDBError(err error, origin string) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			e.LogErr(origin, err, e.Warn)
			return fmt.Errorf("%w", e.ErrAlreadyExists)
		case "23503": // foreign_key_violation
			e.LogErr(origin, err, e.Warn)
			return fmt.Errorf("%w", e.ErrConflict)
		case "23502": // not_null_violation
			e.LogErr(origin, err, e.Warn)
			return fmt.Errorf("%w", e.ErrInvalidInput)
		case "23514": // check_violation
			e.LogErr(origin, err, e.Warn)
			return fmt.Errorf("%w", e.ErrValidationFailed)
		default:
			e.LogErr(origin, pgErr, e.Warn)
			return fmt.Errorf("%w", e.ErrInternal)
		}
	}

	// network errors
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		e.LogErr(origin, netErr, e.Warn)
		return fmt.Errorf("%w", e.ErrServiceDown)
	}

	// context errors
	switch {
	case errors.Is(err, context.DeadlineExceeded),
		errors.Is(err, context.Canceled):
		e.LogErr(origin, err, e.Warn)
		return fmt.Errorf("%w", e.ErrTimeout)
	}

	// fallback
	e.LogErr(origin, err, e.Warn)
	return fmt.Errorf("%w", e.ErrInternal)
}
