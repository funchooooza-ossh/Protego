package adapters

import (
	"context"

	"github.com/funchooooza-ossh/protego/internal/domain"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *domain.User) error
}
