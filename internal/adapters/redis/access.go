package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/funchooooza-ossh/protego/internal/adapters"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type RedisAsideAccess struct {
	rdb      *redis.Client
	delegate adapters.AccessRepositoryInterface

	ttl        time.Duration
	group      singleflight.Group
	localCache *ristretto.Cache
}

func NewRedisAside(rdb *redis.Client, delegate adapters.AccessRepositoryInterface, ttl time.Duration) *RedisAsideAccess {
	lru, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7, // количество ключей для подсчёта частоты (примерно 10x от MaxCost)
		MaxCost:     1e6, // ограничение по "стоимости" (можно считать как количество записей)
		BufferItems: 64,  //TODO env размер буфера (оптимально 64–128)
	})
	return &RedisAsideAccess{
		rdb:        rdb,
		delegate:   delegate,
		ttl:        ttl,
		localCache: lru,
	}
}

func (r *RedisAsideAccess) HasAccess(ctx context.Context, roleID, action, resourceCode string) (bool, error) {
	const origin = "cache_aside.has_access"

	key := r.getKey(roleID, resourceCode, action)

	cached, err := r.get(ctx, key)
	if err != nil {
		e.LogErr(ctx, origin, err, e.Info)
	}
	if cached != nil {
		return *cached, nil
	}

	val, err, _ := r.group.Do(key, func() (any, error) {
		allowed, err := r.delegate.HasAccess(ctx, roleID, action, resourceCode)
		if err != nil {
			return nil, e.ReturnErr(ctx, origin, err, e.Warn)
		}
		r.cacheAccess(ctx, key, origin, allowed)
		return allowed, nil

	})

	if err != nil {
		return false, err
	}

	return val.(bool), nil

}

func (r *RedisAsideAccess) get(ctx context.Context, key string) (*bool, error) {
	const origin = "cache_aside.get"

	//1. Чтение из LRU кэша
	if val, ok := r.localCache.Get(key); ok {
		result := val.(bool)
		return &result, nil
	}

	//2. Если нет в LRU — идём в Redis
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return nil, nil
		default:
			return nil, ParseRedisError(ctx, err, origin)
		}
	}

	result := (val == "1")

	//3. Кэшируем обратно в LRU
	r.localCache.SetWithTTL(key, result, 1, r.ttl)

	return &result, nil
}

func (r *RedisAsideAccess) cacheAccess(ctx context.Context, key, origin string, allowed bool) {
	// LRU cache — sync, контекст не нужен
	r.localCache.SetWithTTL(key, allowed, 1, r.ttl)

	// Асинхронный пуш в Redis
	go func(parentCtx context.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				e.LogErr(parentCtx, origin, fmt.Errorf("panic in cacheAccess: %v", rec), e.Info)
			}
		}()

		// Создаем scoped timeout контекст на базе родительского
		ctx, cancel := context.WithTimeout(parentCtx, 300*time.Millisecond)
		defer cancel()

		if err := r.set(ctx, key, allowed); err != nil {
			e.LogErr(ctx, origin, err, e.Info)
		}
	}(ctx)
}

func (r *RedisAsideAccess) set(ctx context.Context, key string, value bool) error {
	const origin = "cache_aside.set"

	val := "0"
	if value {
		val = "1"
	}

	if err := r.rdb.Set(ctx, key, val, r.ttl).Err(); err != nil {
		return ParseRedisError(ctx, err, origin)
	}

	return nil
}

func (r *RedisAsideAccess) getKey(roleID, resourceCode, action string) string {
	return fmt.Sprintf("%s:%s:%s:%s", "perms", roleID, resourceCode, action)
}
