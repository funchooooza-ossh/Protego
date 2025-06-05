package adapters

import (
	"context"

	"github.com/funchooooza-ossh/protego/internal/domain"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
