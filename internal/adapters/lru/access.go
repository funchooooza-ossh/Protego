package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/funchooooza-ossh/protego/internal/logger"
	"go.uber.org/zap/zapcore"
)

type LruAccessCacheRepository struct {
	lru *ristretto.Cache
	ttl time.Duration //seconds
}

func NewLruAccessCacheRepository(numCounters, maxCost int64, bufferItems int64, ttl time.Duration) *LruAccessCacheRepository {
	lru, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: numCounters,
		MaxCost:     maxCost,
		BufferItems: bufferItems,
	})
	return &LruAccessCacheRepository{
		lru: lru,
		ttl: ttl,
	}
}

func (r *LruAccessCacheRepository) Set(ctx context.Context, key interface{}, value any, cost int64) bool {

	val := r.lru.SetWithTTL(key, value, cost, r.ttl)
	if !val {
		logger.Log(ctx, zapcore.DebugLevel, fmt.Sprintf("lru not added for key:%s", key))
	}

	return val

}

func (r *LruAccessCacheRepository) Get(key interface{}) (interface{}, bool) {
	val, hit := r.lru.Get(key)
	return val, hit
}
