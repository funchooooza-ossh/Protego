package contracts

import (
	"github.com/funchooooza-ossh/protego/internal/domain"
)

type JWTProvider interface {
	GenerateToken(*domain.TokenClaims) (string, error)
	VerifyToken(token string, strict bool) (*domain.TokenClaims, error)
}
