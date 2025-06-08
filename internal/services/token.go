package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	adapters "github.com/funchooooza-ossh/protego/internal/adapters"
	"github.com/funchooooza-ossh/protego/internal/domain"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/funchooooza-ossh/protego/internal/tokens"
	"github.com/google/uuid"
)

type TokenService struct {
	jwtManager *tokens.JWTManager
	tokenRepo  adapters.CacheRepositoryInterface
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewTokenService(
	jwt *tokens.JWTManager,
	repo adapters.CacheRepositoryInterface,
	accessTTL, refreshTTL time.Duration,
) *TokenService {
	return &TokenService{
		jwtManager: jwt,
		tokenRepo:  repo,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *TokenService) CreatePair(ctx context.Context, baseClaims *domain.TokenClaims) (string, string, error) {
	const origin = "token_service.CreatePair"

	now := time.Now().Unix()
	jti := uuid.NewString() // generate unique string for session

	accessClaims := &domain.TokenClaims{
		UserID:    baseClaims.UserID,
		RoleID:    baseClaims.RoleID,
		Type:      "access",
		JTI:       jti,
		IssuedAt:  now,
		ExpiresAt: now + int64(s.accessTTL.Seconds()),
	}

	refreshClaims := &domain.TokenClaims{
		UserID:    baseClaims.UserID,
		RoleID:    baseClaims.RoleID,
		Type:      "refresh",
		JTI:       jti,
		IssuedAt:  now,
		ExpiresAt: now + int64(s.refreshTTL.Seconds()),
	}
	// generate both tokens
	accessToken, err := s.jwtManager.GenerateToken(accessClaims)
	if err != nil {
		err = fmt.Errorf("%w creating access token", e.ErrInternal)
		return "", "", e.ReturnErr(ctx, origin, err, e.Warn)
	}
	refreshToken, err := s.jwtManager.GenerateToken(refreshClaims)
	if err != nil {
		err = fmt.Errorf("%w creating refresh token", e.ErrInternal)
		return "", "", e.ReturnErr(ctx, origin, err, e.Warn)
	}
	//store session into redis
	sessionKey := s.sessionKey(baseClaims.UserID)
	if err := s.tokenRepo.Set(ctx, sessionKey, jti); err != nil {
		return "", "", e.ReturnErr(ctx, origin, err, e.Warn)
	}

	return accessToken, refreshToken, nil

}

func (s *TokenService) RefreshAccess(ctx context.Context, refresh string) (*domain.TokenClaims, string, error) {
	const origin = "token_service.RefreshAccess"

	//1. Парсим refresh
	refreshClaims, err := s.jwtManager.VerifyToken(refresh, true)
	if err != nil {
		err = fmt.Errorf("%w %s", e.ErrInvalidInput, err.Error()) // сессией управляет хранимый id сессии
		return nil, "", e.ReturnErr(ctx, origin, err, e.Info)

	}

	// 2. Проверка jti-сессии
	sessionKey := s.sessionKey(refreshClaims.UserID)
	storedJTI, err := s.tokenRepo.Get(ctx, sessionKey)
	switch {
	case errors.Is(err, e.ErrNotFound), storedJTI == "":
		return nil, "", e.ReturnErr(ctx, origin, e.ErrUnauthorized, e.Info)

	case err != nil:
		return nil, "", e.ReturnErr(ctx, origin, err, e.Warn)

	case storedJTI != refreshClaims.JTI:
		return nil, "", e.ReturnErr(ctx, origin, e.ErrUnauthorized, e.Info)
	}

	//3.Создаем новый токен
	now := time.Now().Unix()
	newAccessClaims := &domain.TokenClaims{
		UserID:    refreshClaims.UserID,
		RoleID:    refreshClaims.RoleID,
		Type:      "access",
		JTI:       refreshClaims.JTI,
		IssuedAt:  now,
		ExpiresAt: now + int64(s.accessTTL.Seconds()),
	}

	newAccess, err := s.jwtManager.GenerateToken(newAccessClaims)
	if err != nil {
		err = fmt.Errorf("%w creating access token", e.ErrInternal)
		return nil, "", e.ReturnErr(ctx, origin, err, e.Warn)

	}

	return newAccessClaims, newAccess, nil
}

func (s *TokenService) InvalidatePair(ctx context.Context, token string) error {
	const origin = "token_service.InvalidatePair"

	claims, err := s.jwtManager.VerifyToken(token, true)
	if err != nil {
		err = fmt.Errorf("%w: token", e.ErrInvalidInput)
		return e.ReturnErr(ctx, origin, err, e.Info)
	}
	sessionKey := s.sessionKey(claims.UserID)

	if err = s.tokenRepo.Delete(ctx, sessionKey); err != nil {
		return e.ReturnErr(ctx, origin, err, e.Warn)
	}
	return nil

}

func (s *TokenService) sessionKey(id string) string {
	return "session:" + id
}

func (s *TokenService) GetClaimsFromToken(ctx context.Context, token string) (*domain.TokenClaims, error) {
	const origin = "token_service.GetClaimsFromToken"
	// Парсим JWT
	claims, err := s.jwtManager.VerifyToken(token, true)
	if err != nil {
		err = fmt.Errorf("%w: token", e.ErrInvalidInput)
		return nil, e.ReturnErr(ctx, origin, err, e.Info) // token протух или не наш
	}

	// Возвращаем данные о пользователе
	return claims, nil
}
