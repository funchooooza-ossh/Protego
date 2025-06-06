package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/funchooooza-ossh/protego/internal/adapters"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type RedisAsideAccess struct {
	rdb      *redis.Client
	delegate adapters.AccessRepositoryInterface

	ttl   time.Duration
	group singleflight.Group
}

func NewRedisAside(rdb *redis.Client, delegate adapters.AccessRepositoryInterface, ttl time.Duration) *RedisAsideAccess {
	return &RedisAsideAccess{
		rdb:      rdb,
		delegate: delegate,
		ttl:      ttl,
	}
}

func (r *RedisAsideAccess) HasAccess(ctx context.Context, roleID, action, resourceCode string) (bool, error) {
	const origin = "cache_aside.has_access"

	key := r.getKey(roleID, resourceCode, action)

	cached, err := r.get(ctx, key)
	if err != nil {
		e.LogErr(origin, err, e.Info)
	}
	if cached != nil {
		return *cached, nil
	}

	val, err, _ := r.group.Do(key, func() (any, error) {
		allowed, err := r.delegate.HasAccess(ctx, roleID, action, resourceCode)
		if err != nil {
			return nil, e.ReturnErr(origin, err, e.Warn)
		}
		r.cacheAccess(key, origin, allowed)
		return allowed, nil

	})

	if err != nil {
		return false, err
	}

	return val.(bool), nil

}

func (r *RedisAsideAccess) get(ctx context.Context, key string) (*bool, error) {
	const origin = "cache_aside.get"

	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return nil, nil
		default:
			return nil, ParseRedisError(err, origin)
		}
	}

	result := (val == "1")
	return &result, nil
}

func (r *RedisAsideAccess) cacheAccess(key, origin string, allowed bool) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				e.LogErr(origin, fmt.Errorf("panic in cacheAccess: %v", r), e.Info)
			}
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()

		if err := r.set(ctx, key, allowed); err != nil {
			e.LogErr(origin, err, e.Info)
		}
	}()
}

func (r *RedisAsideAccess) set(ctx context.Context, key string, value bool) error {
	const origin = "cache_aside.set"

	val := "0"
	if value {
		val = "1"
	}

	if err := r.rdb.Set(ctx, key, val, r.ttl).Err(); err != nil {
		return ParseRedisError(err, origin)
	}

	return nil
}

func (r *RedisAsideAccess) getKey(roleID, resourceCode, action string) string {
	return fmt.Sprintf("%s:%s:%s:%s", "perms", roleID, resourceCode, action)
}
