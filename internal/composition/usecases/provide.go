package compositionUsecases

import (
	services "github.com/funchooooza-ossh/protego/internal/composition/services"
	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/usecases"
)

type Usecaess struct {
	Register contracts.RegisterUsecaseInterface
	Login    contracts.LoginUsecaseInterface
	Auth     contracts.AuthUsecaseInterface
	Logout   contracts.LogoutUsecaseInterface
}

func NewUsecases(servs *services.Services, cfg *config.Config) *Usecaess {
	Register := usecases.NewRegisterUsecase(servs.User)
	Login := usecases.NewLoginUsecase(servs.User, servs.Token, cfg.LoginAttempts)
	Auth := usecases.NewAuthUsecase(servs.User, servs.Token)
	Logout := usecases.NewLogoutUsecase(servs.Token)

	return &Usecaess{
		Register: Register,
		Login:    Login,
		Auth:     Auth,
		Logout:   Logout,
	}

}
