package compositionInfrastructure

import (
	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/infra"
)

func newCacheAsideAccess(
	redisAcces contracts.CacheRepositoryInterface,
	dbAccess contracts.AccessRepositoryInterface,
	lruAccess contracts.LruCacheInterface,
) contracts.AccessRepositoryInterface {
	return infra.NewAccessCacheAside(
		redisAcces,
		dbAccess,
		lruAccess,
	)
}
