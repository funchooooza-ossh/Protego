package usecases

import (
	"context"

	"github.com/funchooooza-ossh/protego/internal/domain"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/funchooooza-ossh/protego/internal/logger"
	"github.com/funchooooza-ossh/protego/internal/services"
	"go.uber.org/zap/zapcore"
)

type RegisterUsecase struct {
	userService services.UserServiceInterface
}

func NewRegisterUsecase(service services.UserServiceInterface) *RegisterUsecase {
	return &RegisterUsecase{
		userService: service,
	}
}

func (u *RegisterUsecase) Execute(ctx context.Context, email, password string) (*domain.User, error) {
	const origin = "register_usecase"

	logger.Log(ctx, zapcore.InfoLevel, "register request")
	user, err := u.userService.CreateUser(ctx, email, password)
	if err != nil {
		return nil, e.ReturnErr(ctx, origin, err, e.Info)
	}

	logger.Log(ctx, zapcore.InfoLevel, "register success")
	return user, nil
}
