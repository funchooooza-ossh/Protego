package contracts

import "context"

type CacheRepositoryInterface interface { // All of Redis repos works with already completed keys
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (value string, err error)
	Delete(ctx context.Context, key string) error
}

type CounterRepositoryInterface interface {
	Increment(ctx context.Context, key string) (int, error)
	Delete(ctx context.Context, key string) error
}
