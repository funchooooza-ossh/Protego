package usecases

import (
	"context"

	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/funchooooza-ossh/protego/internal/services"
)

type LogoutUsecase struct {
	tokenService services.TokenServiceInterface
}

func NewLogoutUsecase(tokenService services.TokenServiceInterface) *LogoutUsecase {
	return &LogoutUsecase{
		tokenService: tokenService,
	}
}

func (u *LogoutUsecase) Execute(ctx context.Context, access string) error {
	const origin = "logout_usecase"

	if err := u.tokenService.InvalidatePair(ctx, access); err != nil {
		return e.ReturnErr(origin, err, e.Info)
	}

	return nil
}
