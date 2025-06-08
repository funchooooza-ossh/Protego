package compositionInfrastructure

import (
	compositionAdapters "github.com/funchooooza-ossh/protego/internal/composition/adapters"
	"github.com/funchooooza-ossh/protego/internal/contracts"
)

type Infrastructure struct {
	Access contracts.AccessRepositoryInterface
}

func NewInfrastructure(adapters *compositionAdapters.Repositories) *Infrastructure {

	cacheAsideAccess := newCacheAsideAccess(adapters.RedisAccess, adapters.DbAccess, adapters.LruAccess)

	return &Infrastructure{
		Access: cacheAsideAccess,
	}
}
