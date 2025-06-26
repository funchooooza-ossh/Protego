package compositionAdapters

import (
	conns "github.com/funchooooza-ossh/protego/internal/composition/connections"
	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/funchooooza-ossh/protego/internal/contracts"
)

type Repositories struct {
	User        contracts.UserRepositoryInterface
	Session     contracts.CacheRepositoryInterface
	Counter     contracts.CounterRepositoryInterface
	Role        contracts.RoleRepositoryInterface
	LruAccess   contracts.LruCacheInterface
	RedisAccess contracts.CacheRepositoryInterface
	DbAccess    contracts.AccessRepositoryInterface
}

func NewRepositories(conns *conns.InfraConnections, cfg *config.Config) *Repositories {
	var db = conns.Queries
	var redis = conns.Redis

	user := newUserDbAdapter(db)
	role := newRoledDbAdapter(db)
	accessRepo := newAccessDbAdapter(db)

	lruAccessCache := newAccesLruAdapter(cfg)
	redisAccessCache := newAccessRedisAdapter(redis, cfg.AccessCacheTTL)

	session := newSessionRedisAdapter(redis, cfg.RefreshTtl)
	counter := newLoginCounterRedisAdapter(redis, cfg.LoginCounterTTL)

	return &Repositories{
		User:        user,
		Role:        role,
		Session:     session,
		LruAccess:   lruAccessCache,
		RedisAccess: redisAccessCache,
		DbAccess:    accessRepo,
		Counter:     counter,
	}
}
