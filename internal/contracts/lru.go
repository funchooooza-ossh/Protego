package contracts

import "context"

type LruCacheInterface interface {
	Set(ctx context.Context, key interface{}, value any, cost int64) bool
	Get(key interface{}) (interface{}, bool)
}
