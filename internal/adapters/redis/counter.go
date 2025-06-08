package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/funchooooza-ossh/protego/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap/zapcore"
)

type CounterRepository struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewCounterRepository(rdb *redis.Client, ttl time.Duration) *CounterRepository {
	return &CounterRepository{
		rdb: rdb,
		ttl: ttl,
	}
}

func (r *CounterRepository) Increment(ctx context.Context, key string) (int, error) {
	const origin = "counter_repo.set"

	val, err := r.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, ParseRedisError(ctx, err, origin)
	}

	if val == 1 {
		err = r.rdb.Expire(ctx, key, r.ttl).Err()
		if err != nil {
			return int(val), ParseRedisError(ctx, err, origin)
		}
	}

	return int(val), nil
}

func (r *CounterRepository) Delete(ctx context.Context, key string) error {
	deleted, err := r.rdb.Del(ctx, key).Result()
	if err != nil {
		return ParseRedisError(ctx, err, "CounterRepository.Delete")
	}
	if deleted == 0 {
		logger.Log(ctx, zapcore.InfoLevel, fmt.Sprintf("redis key not found: %s", key))
	}
	return nil
}
