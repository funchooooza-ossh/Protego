package contracts

import (
	"context"

	"github.com/funchooooza-ossh/protego/internal/domain"
)

type JWTProvider interface {
	GenerateToken(ctx context.Context, claims *domain.TokenClaims) (string, error)
	VerifyToken(ctx context.Context, token string) (*domain.TokenClaims, error)
}
