package adapters

import (
	"context"
	"errors"
	"fmt"
	"net"

	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/redis/go-redis/v9"
)

func ParseRedisError(err error, origin string) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, redis.Nil):
		e.LogErr(origin, e.ErrNotFound, e.Info)
		return e.ErrNotFound

	case errors.As(err, new(*net.OpError)):
		e.LogErr(origin, e.ErrServiceDown, e.Warn)
		return fmt.Errorf("%w: network issue", e.ErrServiceDown)

	case errors.Is(err, context.DeadlineExceeded),
		errors.Is(err, context.Canceled):
		e.LogErr(origin, e.ErrTimeout, e.Warn)
		return fmt.Errorf("%w: timeout or cancel", e.ErrTimeout)

	default:
		e.LogErr(origin, e.ErrInternal, e.Warn)
		return fmt.Errorf("%w: %v", e.ErrInternal, err)
	}
}
