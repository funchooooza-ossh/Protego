package services

import (
	"context"

	"github.com/funchooooza-ossh/protego/internal/domain"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, email, password string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	HasPermission(ctx context.Context, roleID string, resourceCode string, action string) (bool, error)
	VerifyPassword(ctx context.Context, password string, hashed string) (bool, error)
	IncreaseCounter(ctx context.Context, id string) (int, error)
	DeleteCounter(ctx context.Context, id string) error
	BlockUser(ctx context.Context, id string) error
}

type TokenServiceInterface interface {
	CreatePair(ctx context.Context, claims *domain.TokenClaims) (access string, refresh string, err error)
	GetClaimsFromToken(ctx context.Context, token string) (*domain.TokenClaims, error)
	RefreshAccess(ctx context.Context, refreshToken string) (*domain.TokenClaims, string, error)
	InvalidatePair(ctx context.Context, token string) error
}
