package adapters

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"reflect"

	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/jackc/pgx/v5"
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

	// Попробуем reflection-based распаковку pgconn.PgError
	if code, ok := extractSQLStateCode(err); ok {
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

	e.LogErr(origin, err, e.Warn)
	return e.ErrInternal
}

func extractSQLStateCode(err error) (string, bool) { // необходимый костыль, к сожалению иного решения я так и не нашел
	// Type must be named *pgconn.PgError
	t := reflect.TypeOf(err)
	if t == nil || t.Kind() != reflect.Ptr {
		return "", false
	}

	if t.String() != "*pgconn.PgError" {
		return "", false
	}

	// Пытаемся достать поле Code через reflect
	v := reflect.ValueOf(err).Elem()
	codeField := v.FieldByName("Code")
	if !codeField.IsValid() || codeField.Kind() != reflect.String {
		return "", false
	}

	return codeField.String(), true
}
