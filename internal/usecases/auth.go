package usecases

import (
	"context"

	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/domain"
	e "github.com/funchooooza-ossh/protego/internal/errors"
)

type AuthUsecase struct {
	userService  contracts.UserServiceInterface
	tokenService contracts.TokenServiceInterface
}

func NewAuthUsecase(userService contracts.UserServiceInterface, tokenService contracts.TokenServiceInterface) *AuthUsecase {
	return &AuthUsecase{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (u *AuthUsecase) Execute(ctx context.Context, resource, action, access, refresh string) (bool, string, error) {
	const origin = "auth_usecase"

	claims, newAccess, err := u.tryGetClaims(ctx, access, refresh)
	if err != nil {
		return false, "", e.ReturnErr(ctx, origin, err, e.Info)
	}

	allowed, err := u.userService.HasPermission(ctx, claims.RoleID, action, resource)
	if err != nil {
		return false, newAccess, e.ReturnErr(ctx, origin, err, e.Info)
	}

	return allowed, newAccess, nil
}

func (u *AuthUsecase) tryGetClaims(ctx context.Context, access, refresh string) (*domain.TokenClaims, string, error) {
	const origin = "auth_usecase"

	if access != "" {
		claims, err := u.tokenService.GetClaimsFromToken(ctx, access)
		if err == nil {
			return claims, "", nil
		}
	}

	if refresh != "" {
		claims, newAccess, err := u.tokenService.RefreshAccess(ctx, refresh)
		if err == nil {
			return claims, newAccess, nil
		}
		// если refresh есть, но он невалиден → 401
		return nil, "", e.ReturnErr(ctx, origin, e.ErrUnauthorized, e.Info)
	}

	// оба токена отсутствуют или невалидны
	return nil, "", e.ReturnErr(ctx, origin, e.ErrUnauthorized, e.Info)
}
