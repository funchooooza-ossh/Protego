package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/funchooooza-ossh/protego/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap/zapcore"
)

type RedisCacheAccesRepository struct {
	rdb *redis.Client
	ttl time.Duration // seconds
}

func NewCacheAccesRepository(rdb *redis.Client, ttl time.Duration) *RedisCacheAccesRepository {
	return &RedisCacheAccesRepository{
		rdb: rdb,
		ttl: ttl,
	}
}

func (r *RedisCacheAccesRepository) Set(ctx context.Context, key string, value string) error {
	const origin = "access_cache.set"

	if err := r.rdb.Set(ctx, key, value, r.ttl).Err(); err != nil {
		return ParseRedisError(ctx, err, origin)
	}
	return nil
}

func (r *RedisCacheAccesRepository) Get(ctx context.Context, key string) (string, error) {
	const origin = "access_cache.get"

	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", ParseRedisError(ctx, err, origin)
	}
	return val, nil
}

func (r *RedisCacheAccesRepository) Delete(ctx context.Context, key string) error {
	deleted, err := r.rdb.Del(ctx, key).Result()
	if err != nil {
		return ParseRedisError(ctx, err, "CounterRepository.Delete")
	}
	if deleted == 0 {
		logger.Log(ctx, zapcore.InfoLevel, fmt.Sprintf("redis key not found: %s", key))
	}
	return nil

}
