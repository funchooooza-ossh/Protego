package usecases

import (
	"context"
	"fmt"

	"github.com/funchooooza-ossh/protego/internal/domain"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/funchooooza-ossh/protego/internal/services"
)

type AuthUsecase struct {
	userService  services.UserServiceInterface
	tokenService services.TokenServiceInterface
}

func NewAuthUsecase(userService services.UserServiceInterface, tokenService services.TokenServiceInterface) *AuthUsecase {
	return &AuthUsecase{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (u *AuthUsecase) Execute(ctx context.Context, resource, action, access, refresh string) (bool, string, error) {
	const origin = "auth_usecase"

	claims, newAccess, err := u.tryGetClaims(ctx, access, refresh)
	if err != nil {
		return false, "", e.ReturnErr(origin, err, e.Info)
	}

	allowed, err := u.userService.HasPermission(ctx, claims.RoleID, action, resource)
	if err != nil {
		return false, newAccess, e.ReturnErr(origin, err, e.Info)
	}

	return allowed, newAccess, nil
}

func (u *AuthUsecase) tryGetClaims(ctx context.Context, access, refresh string) (*domain.TokenClaims, string, error) {
	const origin = "auth_usecase"
	// пробуем получить claims из access token
	claims, err := u.tokenService.GetClaimsFromToken(ctx, access)
	if err == nil {
		return claims, access, nil
	}

	// fallback на refresh
	newClaims, newAccess, err := u.tokenService.RefreshAccess(ctx, refresh)
	if err != nil {
		// если refresh выкинет ошибку => мы не сможем узнать от кого запрос = 401.
		err = fmt.Errorf("%w:%s", e.ErrUnauthorized, err.Error())
		return nil, "", e.ReturnErr(origin, err, e.Info)
	}
	// возвращаем обновленный access и claims
	return newClaims, newAccess, nil
}
