package infra

import (
	"errors"
	"time"

	"github.com/funchooooza-ossh/protego/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret string
}

type jwtClaims struct {
	UserID string `json:"sub"`
	Role   string `json:"role"`
	Type   string `json:"type"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{secret: secret}
}
func (m *JWTManager) GenerateToken(claims *domain.TokenClaims) (string, error) {
	jwtC := jwtClaims{
		UserID: claims.UserID,
		Role:   claims.RoleID,
		Type:   string(claims.Type),
		JTI:    claims.JTI,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Unix(claims.IssuedAt, 0)),
			ExpiresAt: jwt.NewNumericDate(time.Unix(claims.ExpiresAt, 0)),
			Subject:   claims.UserID,
			ID:        claims.JTI,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtC)
	return token.SignedString([]byte(m.secret))
}

func (m *JWTManager) VerifyToken(tokenStr string, strict bool) (*domain.TokenClaims, error) {
	var parser *jwt.Parser
	if strict {
		parser = jwt.NewParser()
	} else {
		parser = jwt.NewParser(jwt.WithoutClaimsValidation())
	}

	token, err := parser.ParseWithClaims(tokenStr, &jwtClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(m.secret), nil
	})
	if err != nil {
		return nil, err // включает подпись, формат, exp (если strict)
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid && strict {
		return nil, errors.New("invalid claims or token")
	}

	return &domain.TokenClaims{
		UserID:    claims.UserID,
		RoleID:    claims.Role,
		Type:      domain.TokenType(claims.Type),
		JTI:       claims.JTI,
		IssuedAt:  claims.IssuedAt.Unix(),
		ExpiresAt: claims.ExpiresAt.Unix(),
	}, nil
}
