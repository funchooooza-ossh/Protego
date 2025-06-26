package compositionAdapters

import (
	adapters "github.com/funchooooza-ossh/protego/internal/adapters/lru"
	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/funchooooza-ossh/protego/internal/contracts"
	authmetrics "github.com/funchooooza-ossh/protego/internal/metrics/auth"
)

func newAccesLruAdapter(cfg *config.Config) contracts.LruCacheInterface {
	lru := adapters.NewLruAccessCacheRepository(
		1e7, //TODO env
		1e6,
		64,
		cfg.AccessCacheTTL,
	)

	return authmetrics.NewAccessLruWithMetrics(lru) //wrapped

}
