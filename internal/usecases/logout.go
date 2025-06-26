package usecases

import (
	"context"
	"errors"

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

	err := u.tokenService.InvalidatePair(ctx, access)
	if err != nil {
		switch {
		case errors.Is(err, e.ErrInvalidInput),
			errors.Is(err, e.ErrUnauthorized):
			return e.ReturnErr(ctx, origin, err, e.Info)
		default:
			return e.ReturnErr(ctx, origin, err, e.Warn) //TODO clean error switching
		}
	}

	return nil
}
