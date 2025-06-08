package compositiomServices

import (
	adapters "github.com/funchooooza-ossh/protego/internal/composition/adapters"
	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/services"
	"github.com/funchooooza-ossh/protego/internal/tokens"
)

type Services struct {
	User  contracts.UserServiceInterface
	Token contracts.TokenServiceInterface
}

func NewServices(repos *adapters.Repositories, cfg *config.Config) *Services {
	userService := services.NewUserService(repos.User,
		repos.Role,
		repos.Counter,
		repos.Access,
		"user", //TODO env
		cfg.PassCost,
	)
	jwt := tokens.NewJWTManager(cfg.JWTSecret)
	tokenService := services.NewTokenService(jwt, repos.Session, cfg.AccessTtl, cfg.RefreshTtl)

	return &Services{
		User:  userService,
		Token: tokenService,
	}
}
