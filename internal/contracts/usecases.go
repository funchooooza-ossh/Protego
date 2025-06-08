package contracts

import (
	"context"

	"github.com/funchooooza-ossh/protego/internal/domain"
)

type LoginUsecaseInterface interface {
	Execute(ctx context.Context, email, password string) (access_token, refresh_token string, err error)
}

type RegisterUsecaseInterface interface {
	Execute(ctx context.Context, email, password string) (*domain.User, error)
}

type AuthUsecaseInterface interface {
	Execute(ctx context.Context, resource, action, access, refresh string) (bool, string, error)
}

type LogoutUsecaseInterface interface {
	Execute(ctx context.Context, access string) error
}
