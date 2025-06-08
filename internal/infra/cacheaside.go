package infra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/funchooooza-ossh/protego/internal/adapters"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/funchooooza-ossh/protego/internal/logger"
	"golang.org/x/sync/singleflight"
)

type AccessCacheAside struct {
	redisCache adapters.CacheRepositoryInterface
	delegate   adapters.AccessRepositoryInterface

	group    singleflight.Group
	lruCache adapters.LruCacheInterface
}

func NewAccessCacheAside(
	redisCache adapters.CacheRepositoryInterface,
	delegate adapters.AccessRepositoryInterface,
	lru adapters.LruCacheInterface,
) *AccessCacheAside {
	return &AccessCacheAside{
		redisCache: redisCache,
		delegate:   delegate,
		lruCache:   lru,
	}
}

func (r *AccessCacheAside) HasAccess(ctx context.Context, roleID, action, resourceCode string) (bool, error) {
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

func (r *AccessCacheAside) get(ctx context.Context, key string) (*bool, error) {
	const origin = "cache_aside.get"

	//1. Чтение из LRU кэша
	if val, ok := r.lruCache.Get(key); ok {
		result := val.(bool)
		return &result, nil
	}

	//2. Если нет в LRU — идём в Redis
	val, err := r.redisCache.Get(ctx, key)
	if err != nil {
		switch {
		case errors.Is(err, e.ErrNotFound):
			return nil, nil
		default:
			return nil, e.ReturnErr(ctx, origin, err, e.Info)
		}
	}

	result := (val == "1")
	//3.Ставим в LRU
	r.lruCache.Set(ctx, key, result, 1)

	return &result, nil
}

func (r *AccessCacheAside) cacheAccess(ctx context.Context, key, origin string, allowed bool) {
	// LRU cache — sync, контекст не нужен
	r.lruCache.Set(ctx, key, allowed, 1)

	// Асинхронный пуш в Redis
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				e.LogErr(ctx, origin, fmt.Errorf("panic in cacheAccess: %v", rec), e.Info)
			}
		}()

		// Detached ctx, но переносим request_id
		bg := context.Background()
		if rid := logger.GetRequestID(ctx); rid != "" {
			bg = context.WithValue(bg, logger.CtxKeyRequestID{}, rid)
		}

		// Scoped timeout
		ctxWithTimeout, cancel := context.WithTimeout(bg, 300*time.Millisecond)
		defer cancel()

		if err := r.set(ctxWithTimeout, key, allowed); err != nil {
			e.BestEffort(ctxWithTimeout, origin, "set redis", err)
		}
	}()
}

func (r *AccessCacheAside) set(ctx context.Context, key string, value bool) error {
	const origin = "cache_aside.set"

	val := "0"
	if value {
		val = "1"
	}

	if err := r.redisCache.Set(ctx, key, val); err != nil {
		return e.ReturnErr(ctx, origin, err, e.Info)
	}

	return nil
}

func (r *AccessCacheAside) getKey(roleID, resourceCode, action string) string {
	return fmt.Sprintf("%s:%s:%s:%s", "perms", roleID, resourceCode, action)
}
