package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"

	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
)

var sqlstateRegex = regexp.MustCompile(`SQLSTATE (\d{5})`)

func ParseDBError(err error, origin string) error {
	if err == nil {
		return nil
	}

	fmt.Printf("Error: %v\nType: %T\n", err, err) // TODO Убрать позже

	// Not found
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) || err.Error() == "no rows in result set" {
		e.LogErr(origin, err, e.Info)
		return e.ErrNotFound
	}
	fmt.Printf("reflect.TypeOf: %v\n", reflect.TypeOf(err))
	fmt.Printf("*pgconn.PgError? %v\n", reflect.TypeOf(err) == reflect.TypeOf(&pgconn.PgError{}))

	// pgconn error codes
	if matches := sqlstateRegex.FindStringSubmatch(err.Error()); len(matches) == 2 {
		code := matches[1]
		e.LogErr(origin, err, e.Warn)

		switch code {
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

	// fallback
	e.LogErr(origin, err, e.Warn)
	return e.ErrInternal
}
