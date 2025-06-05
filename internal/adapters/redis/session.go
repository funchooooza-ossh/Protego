package adapters

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewSessionRepository(rdb *redis.Client, ttl time.Duration) *SessionRepository {
	return &SessionRepository{
		rdb: rdb,
		ttl: ttl,
	}
}

func (r *SessionRepository) Set(ctx context.Context, key, value string) error {
	const origin = "session_repo.set"
	if err := r.rdb.Set(ctx, key, value, r.ttl).Err(); err != nil {
		return ParseRedisError(err, origin)
	}
	return nil
}

func (r *SessionRepository) Get(ctx context.Context, key string) (value string, err error) {
	const origin = "session_repo.get"

	data, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", ParseRedisError(err, origin)
	}
	return data, nil
}

func (r *SessionRepository) Delete(ctx context.Context, key string) error {
	const origin = "session_repo.delete"

	deleted, err := r.rdb.Del(ctx, key).Result()
	if err != nil {
		return ParseRedisError(err, origin)
	}
	if deleted == 0 {
		log.Printf("redis key not found: %s", key)
	}
	return nil
}
