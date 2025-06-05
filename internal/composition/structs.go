package composition

import (
	"time"

	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/funchooooza-ossh/protego/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type InfraConnections struct {
	DBPool  *pgxpool.Pool
	Queries *db.Queries
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

	queries := db.New(pool)

	return &InfraConnections{
		DBPool:  pool,
		Queries: queries,
		Redis:   rdb,
	}, nil
}
