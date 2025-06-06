package usecases

import (
	"context"
	"log"

	"github.com/funchooooza-ossh/protego/internal/domain"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/funchooooza-ossh/protego/internal/services"
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

	log.Println("requested to create user")
	user, err := u.userService.CreateUser(ctx, email, password)
	if err != nil {
		return nil, e.ReturnErr(origin, err, e.Info)
	}

	log.Println("user created")
	return user, nil
}
