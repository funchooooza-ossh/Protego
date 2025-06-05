package composition

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *config.Config, attempts int, delay time.Duration) (*redis.Client, error) {
	var rdb *redis.Client
	var err error

	for i := 0; i < attempts; i++ {
		rdb, err = connectRedisOnce(cfg)
		if err == nil {
			log.Printf("[infra] Redis connected successfully on attempt %d", i+1)
			return rdb, nil
		}

		log.Printf("[infra] Redis connection attempt %d failed: %v", i+1, err)
		time.Sleep(delay)
		delay *= 2 // можно зафиксировать, если не хочешь runaway delays
	}

	return nil, fmt.Errorf("failed to connect to Redis after %d attempts: %w", attempts, err)
}

func connectRedisOnce(cfg *config.Config) (*redis.Client, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[infra] panic during Redis connection: %v", r)
		}
	}()

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		DB:   0,

		// Timeouts
		DialTimeout:           time.Duration(cfg.RedisDialTimeout) * time.Millisecond,
		ReadTimeout:           time.Duration(cfg.RedisReadTimeout) * time.Millisecond,
		WriteTimeout:          time.Duration(cfg.RedisWriteTimeout) * time.Millisecond,
		ContextTimeoutEnabled: true,

		// Pool config
		PoolFIFO:       true,
		PoolSize:       cfg.RedisPool,
		MinIdleConns:   cfg.RedisIdleConns,
		MaxIdleConns:   400, //TODO env
		MaxActiveConns: 500,
		PoolTimeout:    100 * time.Millisecond,

		// Conn lifecycle
		ConnMaxIdleTime: 30 * time.Second,
		ConnMaxLifetime: 5 * time.Minute,

		DisableIdentity: true,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return rdb, nil
}
