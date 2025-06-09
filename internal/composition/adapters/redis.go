package compositionAdapters

import (
	"time"

	adapters "github.com/funchooooza-ossh/protego/internal/adapters/redis"
	"github.com/funchooooza-ossh/protego/internal/contracts"
	authmetrics "github.com/funchooooza-ossh/protego/internal/metrics/auth"
	loginmetrics "github.com/funchooooza-ossh/protego/internal/metrics/login"
	"github.com/redis/go-redis/v9"
)

func newAccessRedisAdapter(rdb *redis.Client, ttl time.Duration) contracts.CacheRepositoryInterface {
	adapter := adapters.NewCacheAccesRepository(rdb, ttl)

	return authmetrics.NewRedisCacheWithMetrics(adapter) // wrapped
}

func newSessionRedisAdapter(rdb *redis.Client, ttl time.Duration) contracts.CacheRepositoryInterface {
	adapter := adapters.NewSessionRepository(rdb, ttl)

	return adapter
}

func newLoginCounterRedisAdapter(rdb *redis.Client, ttl time.Duration) contracts.CounterRepositoryInterface {
	adapter := adapters.NewCounterRepository(rdb, ttl)

	return loginmetrics.NewCounterRepositoryWithMetrics(adapter)
}
