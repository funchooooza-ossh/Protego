package compositionAdapters

import (
	dbadapters "github.com/funchooooza-ossh/protego/internal/adapters/db"
	lruadapters "github.com/funchooooza-ossh/protego/internal/adapters/lru"
	redisadapters "github.com/funchooooza-ossh/protego/internal/adapters/redis"
	conns "github.com/funchooooza-ossh/protego/internal/composition/connections"
	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/infra"
	authmetrics "github.com/funchooooza-ossh/protego/internal/metrics/auth"
)

type Repositories struct {
	User    contracts.UserRepositoryInterface
	Session contracts.CacheRepositoryInterface
	Counter contracts.CounterRepositoryInterface
	Role    contracts.RoleRepositoryInterface
	Access  contracts.AccessRepositoryInterface
}

func NewRepositories(conns *conns.InfraConnections, cfg *config.Config) *Repositories {
	user := dbadapters.NewUserRepository(conns.Queries)
	role := dbadapters.NewRoleRepository(conns.Queries)

	dbAccess := dbadapters.NewAccessRepository(conns.Queries)
	dbAccessWithMetrics := authmetrics.NewDelegateWithMetrics(dbAccess)

	lruAccessCache := lruadapters.NewLruAccessCacheRepository(
		1e7, //TODO env
		1e6,
		64,
		cfg.AccessCacheTTL,
	)
	lruAccessCacheWithMetrics := authmetrics.NewAccessLruWithMetrics(lruAccessCache)

	redisAccessCache := redisadapters.NewCacheAccesRepository(conns.Redis, cfg.AccessCacheTTL)
	redisAccessCacheWithMetrics := authmetrics.NewRedisCacheWithMetrics(redisAccessCache)

	cacheAccess := infra.NewAccessCacheAside(redisAccessCacheWithMetrics, dbAccessWithMetrics, lruAccessCacheWithMetrics) //TODO infra layer struct
	session := redisadapters.NewSessionRepository(conns.Redis, cfg.RefreshTtl)
	counter := redisadapters.NewCounterRepository(conns.Redis, cfg.LoginCounterTTL)

	return &Repositories{
		User:    user,
		Role:    role,
		Session: session,
		Access:  cacheAccess,
		Counter: counter,
	}
}
