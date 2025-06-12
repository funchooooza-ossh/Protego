package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/domain"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/google/uuid"
)

type TokenService struct {
	jwtManager contracts.JWTProvider
	tokenRepo  contracts.CacheRepositoryInterface
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewTokenService(
	jwt contracts.JWTProvider,
	repo contracts.CacheRepositoryInterface,
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
	accessToken, err := s.jwtManager.GenerateToken(ctx, accessClaims)
	if err != nil {
		return "", "", e.ReturnErr(ctx, origin, fmt.Errorf("creating access token: %w", err), e.Error)
	}

	refreshToken, err := s.jwtManager.GenerateToken(ctx, refreshClaims)
	if err != nil {
		return "", "", e.ReturnErr(ctx, origin, fmt.Errorf("creating refresh token: %w", err), e.Error)
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

	// 1. Парсим refresh
	refreshClaims, err := s.jwtManager.VerifyToken(ctx, refresh)
	if err != nil {
		return nil, "", classifyJWTError(ctx, origin, err)
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

	// 3. Генерация access-токена
	now := time.Now().Unix()
	newAccessClaims := &domain.TokenClaims{
		UserID:    refreshClaims.UserID,
		RoleID:    refreshClaims.RoleID,
		Type:      "access",
		JTI:       refreshClaims.JTI,
		IssuedAt:  now,
		ExpiresAt: now + int64(s.accessTTL.Seconds()),
	}

	newAccess, err := s.jwtManager.GenerateToken(ctx, newAccessClaims)
	if err != nil {
		return nil, "", e.ReturnErr(ctx, origin, fmt.Errorf("creating access token: %w", err), e.Error)
	}

	return newAccessClaims, newAccess, nil
}

func (s *TokenService) InvalidatePair(ctx context.Context, token string) error {
	const origin = "token_service.InvalidatePair"

	claims, err := s.jwtManager.VerifyToken(ctx, token)
	if err != nil {
		return classifyJWTError(ctx, origin, err)
	}

	sessionKey := s.sessionKey(claims.UserID)

	if err = s.tokenRepo.Delete(ctx, sessionKey); err != nil {
		return e.ReturnErr(ctx, origin, err, e.Warn)
	}

	return nil
}

func (s *TokenService) GetClaimsFromToken(ctx context.Context, token string) (*domain.TokenClaims, error) {
	const origin = "token_service.GetClaimsFromToken"
	// Парсим JWT
	claims, err := s.jwtManager.VerifyToken(ctx, token)
	if err != nil {
		return nil, classifyJWTError(ctx, origin, err)
	}

	// Возвращаем данные о пользователе
	return claims, nil
}

func (s *TokenService) sessionKey(id string) string {
	return "session:" + id
}

// classifyJWTError маппит ошибки jwt-парсинга в семантически корректные user-facing ошибки
func classifyJWTError(ctx context.Context, origin string, err error) error {
	switch {
	case errors.Is(err, e.ErrTokenExpired),
		errors.Is(err, e.ErrTokenInvalid):
		return e.ReturnErr(ctx, origin, e.ErrInvalidInput, e.Info)
	default:
		return e.ReturnErr(ctx, origin, err, e.Warn)
	}
}
