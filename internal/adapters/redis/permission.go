package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/funchooooza-ossh/protego/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap/zapcore"
)

type PermissionRepository struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewPermissionRepository(rdb *redis.Client, ttl time.Duration) *PermissionRepository {
	return &PermissionRepository{
		rdb: rdb,
		ttl: ttl,
	}
}

func (r *PermissionRepository) Set(ctx context.Context, key, value string) error {
	const origin = "perm_repo.set"

	if err := r.rdb.Set(ctx, key, value, r.ttl).Err(); err != nil {
		return ParseRedisError(ctx, err, origin)
	}
	return nil
}

func (r *PermissionRepository) Get(ctx context.Context, key string) (value string, err error) {
	const origin = "perm_repo.get"

	data, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", ParseRedisError(ctx, err, origin)
	}
	return data, nil
}

func (r *PermissionRepository) Delete(ctx context.Context, key string) error {
	const origin = "perm_repo.delete"

	deleted, err := r.rdb.Del(ctx, key).Result()
	if err != nil {
		return ParseRedisError(ctx, err, origin)
	}
	if deleted == 0 {
		logger.Log(ctx, zapcore.InfoLevel, fmt.Sprintf("redis key not found: %s", key))
	}
	return nil
}
