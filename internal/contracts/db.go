package contracts

import (
	"context"

	"github.com/funchooooza-ossh/protego/internal/domain"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

type RoleRepositoryInterface interface {
	Create(ctx context.Context, role *domain.Role) error
	GetByCode(ctx context.Context, code string) (*domain.Role, error)
}
