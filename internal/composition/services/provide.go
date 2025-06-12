package compositiomServices

import (
	adapters "github.com/funchooooza-ossh/protego/internal/composition/adapters"
	compositionInfrastructure "github.com/funchooooza-ossh/protego/internal/composition/infrastructure"
	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/services"
	"github.com/funchooooza-ossh/protego/internal/tokens"
)

type Services struct {
	User  contracts.UserServiceInterface
	Token contracts.TokenServiceInterface
}

func NewServices(repos *adapters.Repositories, cfg *config.Config, infra *compositionInfrastructure.Infrastructure) *Services {
	userService := services.NewUserService(repos.User,
		repos.Role,
		repos.Counter,
		infra.Access,
		"user", //TODO env
		infra.Hasher,
	)
	jwt := tokens.NewJWTManager(cfg.JWTSecret)
	tokenService := services.NewTokenService(jwt, repos.Session, cfg.AccessTtl, cfg.RefreshTtl)

	return &Services{
		User:  userService,
		Token: tokenService,
	}
}
