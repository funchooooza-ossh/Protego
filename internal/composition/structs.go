package composition

import (
	"time"

	"github.com/funchooooza-ossh/protego/internal/adapters"
	dbadapters "github.com/funchooooza-ossh/protego/internal/adapters/db"
	redisadapters "github.com/funchooooza-ossh/protego/internal/adapters/redis"
	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/funchooooza-ossh/protego/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type InfraConnections struct {
	DBPool  *pgxpool.Pool
	Queries *db.UQueries
	Redis   *redis.Client
}

func NewInfraConnections(cfg *config.Config) (*InfraConnections, error) {
	pool, err := ConnectDB(cfg)
	if err != nil {
		return nil, err
	}

	rdb, err := ConnectRedis(cfg, 3, time.Second)
	if err != nil {
		return nil, err
	}

	queries := db.NewQueries(pool)

	return &InfraConnections{
		DBPool:  pool,
		Queries: queries,
		Redis:   rdb,
	}, nil
}

type Repositories struct {
	User    adapters.UserRepositoryInterface
	Session adapters.CacheRepositoryInterface
	Counter adapters.CounterRepositoryInterface
	Role    adapters.RoleRepositoryInterface
	Access  adapters.AccessRepositoryInterface
}

func NewRepositories(conns *InfraConnections, cfg *config.Config) *Repositories {
	user := dbadapters.NewUserRepository(conns.Queries)
	role := dbadapters.NewRoleRepository(conns.Queries)
	access := dbadapters.NewAccessRepository(conns.Queries)

	cacheAccess := redisadapters.NewRedisAside(conns.Redis, access, cfg.AccessCacheTTL)
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
