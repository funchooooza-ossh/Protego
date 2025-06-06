package helpers

import (
	"context"
)

func SafeWithContext[T any](ctx context.Context, fn func() (T, error)) (T, error) {
	type result struct {
		v   T
		err error
	}

	ch := make(chan result, 1)

	go func() {
		v, err := fn()
		ch <- result{v, err}
	}()

	select {
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	case res := <-ch:
		return res.v, res.err
	}
}
