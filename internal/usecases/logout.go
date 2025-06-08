package usecases

import (
	"context"

	"github.com/funchooooza-ossh/protego/internal/contracts"
	e "github.com/funchooooza-ossh/protego/internal/errors"
)

type LogoutUsecase struct {
	tokenService contracts.TokenServiceInterface
}

func NewLogoutUsecase(tokenService contracts.TokenServiceInterface) *LogoutUsecase {
	return &LogoutUsecase{
		tokenService: tokenService,
	}
}

func (u *LogoutUsecase) Execute(ctx context.Context, access string) error {
	const origin = "logout_usecase"

	if err := u.tokenService.InvalidatePair(ctx, access); err != nil {
		return e.ReturnErr(ctx, origin, err, e.Info)
	}

	return nil
}
