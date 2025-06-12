package infra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/funchooooza-ossh/protego/internal/domain"
	e "github.com/funchooooza-ossh/protego/internal/errors"
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

func (m *JWTManager) GenerateToken(ctx context.Context, claims *domain.TokenClaims) (string, error) {
	const origin = "jwt.generate_token"

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

	signed, err := token.SignedString([]byte(m.secret))
	if err != nil {
		err = fmt.Errorf("%w: %w", e.ErrInternal, err)
		return "", e.ReturnErr(ctx, origin, err, e.Error)
	}

	return signed, nil
}
func (m *JWTManager) VerifyToken(ctx context.Context, tokenStr string) (*domain.TokenClaims, error) {
	const origin = "jwt.verify_token"

	parser := jwt.NewParser()

	token, err := parser.ParseWithClaims(tokenStr, &jwtClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(m.secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, e.ReturnErr(ctx, origin, e.ErrTokenExpired, e.Info)
		}
		if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, e.ReturnErr(ctx, origin, e.ErrTokenInvalid, e.Info)
		}
		err = fmt.Errorf("%w: %w", e.ErrInternal, err)
		return nil, e.ReturnErr(ctx, origin, err, e.Error)
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, e.ReturnErr(ctx, origin, e.ErrTokenInvalid, e.Info)
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
