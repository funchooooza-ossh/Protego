package composition

import (
	"context"
	"fmt"
	"time"

	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseDsn())
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	poolCfg.MaxConns = 150 //TODO env
	poolCfg.MinConns = 30
	poolCfg.MaxConnLifetime = 5 * time.Minute
	poolCfg.MaxConnIdleTime = 30 * time.Second
	poolCfg.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("pgx pool connect: %w", err)
	}

	// Проверка соединения
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pg ping: %w", err)
	}

	return pool, nil
}
