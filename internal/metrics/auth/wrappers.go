package authmetrics

import (
	"context"
	"time"

	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/helpers"
)

type AccessLruWithMetrics struct {
	cache contracts.LruCacheInterface
}

func NewAccessLruWithMetrics(cache contracts.LruCacheInterface) *AccessLruWithMetrics {
	return &AccessLruWithMetrics{cache: cache}
}

func (m *AccessLruWithMetrics) Get(key interface{}) (interface{}, bool) {
	start := time.Now()
	val, hit := m.cache.Get(key)
	helpers.ObserveDuration(LruLastDelay, LruDelaySummary, start)

	if hit {
		LruHits.Inc()
	} else {
		LruMisses.Inc()
	}
	return val, hit
}

func (m *AccessLruWithMetrics) Set(ctx context.Context, key interface{}, value any, cost int64) bool {
	start := time.Now()
	ok := m.cache.Set(ctx, key, value, cost)
	helpers.ObserveDuration(LruLastDelay, LruDelaySummary, start)

	return ok
}

type AccessRedisWithMetrics struct {
	cache contracts.CacheRepositoryInterface
}

func NewRedisCacheWithMetrics(impl contracts.CacheRepositoryInterface) *AccessRedisWithMetrics {
	return &AccessRedisWithMetrics{cache: impl}
}

func (m *AccessRedisWithMetrics) Get(ctx context.Context, key string) (string, error) {
	start := time.Now()
	val, err := m.cache.Get(ctx, key)
	helpers.ObserveDuration(RedisLastDelay, RedisDelaySummary, start)

	if err == nil && val != "" {
		RedisHits.Inc()
	} else {
		RedisMisses.Inc()
	}
	return val, err
}

func (m *AccessRedisWithMetrics) Set(ctx context.Context, key string, value string) error {
	start := time.Now()
	err := m.cache.Set(ctx, key, value)
	helpers.ObserveDuration(RedisLastDelay, RedisDelaySummary, start)

	return err
}

func (m *AccessRedisWithMetrics) Delete(ctx context.Context, key string) error {
	return m.cache.Delete(ctx, key)
}

type DBAccessWithMetrics struct {
	db contracts.AccessRepositoryInterface
}

func NewDelegateWithMetrics(impl contracts.AccessRepositoryInterface) *DBAccessWithMetrics {
	return &DBAccessWithMetrics{db: impl}
}

func (m *DBAccessWithMetrics) HasAccess(ctx context.Context, roleID, action, resourceCode string) (bool, error) {
	start := time.Now()
	result, err := m.db.HasAccess(ctx, roleID, action, resourceCode)
	DelegateRequests.Inc()
	helpers.ObserveDuration(DelegateLastDelay, DelegateDelaySummary, start)

	return result, err
}
