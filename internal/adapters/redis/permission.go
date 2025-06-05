package adapters

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
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
		return ParseRedisError(err, origin)
	}
	return nil
}

func (r *PermissionRepository) Get(ctx context.Context, key string) (value string, err error) {
	const origin = "perm_repo.get"

	data, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", ParseRedisError(err, origin)
	}
	return data, nil
}

func (r *PermissionRepository) Delete(ctx context.Context, key string) error {
	const origin = "perm_repo.delete"

	deleted, err := r.rdb.Del(ctx, key).Result()
	if err != nil {
		return ParseRedisError(err, origin)
	}
	if deleted == 0 {
		log.Printf("redis key not found: %s", key)
	}
	return nil
}
