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

var incrAndExpire = redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
  redis.call("EXPIRE", KEYS[1], ARGV[1])
end
return current
`) //use lua to avoid edge-cases with redis fails

func (r *CounterRepository) Increment(ctx context.Context, key string) (int, error) {
	const origin = "counter_repo.lua"

	val, err := incrAndExpire.Run(ctx, r.rdb, []string{key}, int(r.ttl.Seconds())).Int()
	if err != nil {
		return 0, ParseRedisError(ctx, err, origin)
	}
	return val, nil
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
