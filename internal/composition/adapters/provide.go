package compositionAdapters

import (
	dbadapters "github.com/funchooooza-ossh/protego/internal/adapters/db"
	lruadapters "github.com/funchooooza-ossh/protego/internal/adapters/lru"
	redisadapters "github.com/funchooooza-ossh/protego/internal/adapters/redis"
	conns "github.com/funchooooza-ossh/protego/internal/composition/connections"
	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/infra"
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
	access := dbadapters.NewAccessRepository(conns.Queries)

	lruAccessCache := lruadapters.NewLruAccessCacheRepository(
		1e7, //TODO env
		1e6,
		64,
		cfg.AccessCacheTTL,
	)

	redisAccessCache := redisadapters.NewCacheAccesRepository(conns.Redis, cfg.AccessCacheTTL)
	cacheAccess := infra.NewAccessCacheAside(redisAccessCache, access, lruAccessCache) //TODO infra layer struct
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
