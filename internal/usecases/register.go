package usecases

import (
	"context"
	"errors"

	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/domain"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/funchooooza-ossh/protego/internal/logger"
	"go.uber.org/zap/zapcore"
)

type RegisterUsecase struct {
	userService contracts.UserServiceInterface
}

func NewRegisterUsecase(service contracts.UserServiceInterface) *RegisterUsecase {
	return &RegisterUsecase{
		userService: service,
	}
}

func (u *RegisterUsecase) Execute(ctx context.Context, email, password string) (*domain.User, error) {
	const origin = "register_usecase"

	logger.Log(ctx, zapcore.InfoLevel, "register request")

	user, err := u.userService.CreateUser(ctx, email, password)
	if err != nil {
		switch {
		case errors.Is(err, e.ErrAlreadyExists):
			return nil, e.ReturnErr(ctx, origin, err, e.Info)
		default:
			return nil, e.ReturnErr(ctx, origin, err, e.Warn) //TODO clean error switching
		}
	}

	logger.Log(ctx, zapcore.InfoLevel, "register success")
	return user, nil
}
