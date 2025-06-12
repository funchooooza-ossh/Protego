package compositionInfrastructure

import (
	compositionAdapters "github.com/funchooooza-ossh/protego/internal/composition/adapters"
	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/funchooooza-ossh/protego/internal/contracts"
)

type Infrastructure struct {
	Access     contracts.AccessRepositoryInterface
	Hasher     contracts.PasswordHasherInterface
	JWTManager contracts.JWTProvider
}

func NewInfrastructure(adapters *compositionAdapters.Repositories, cfg *config.Config) *Infrastructure {

	cacheAsideAccess := newCacheAsideAccess(adapters.RedisAccess, adapters.DbAccess, adapters.LruAccess)
	hasher := newPasswordHasher(cfg.PassCost)
	jwt := newJWTManager(cfg.JWTSecret)
	return &Infrastructure{
		Access:     cacheAsideAccess,
		Hasher:     hasher,
		JWTManager: jwt,
	}
}
