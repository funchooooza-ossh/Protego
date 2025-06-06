package usecases

import (
	"context"
	"fmt"
	"log"

	"github.com/funchooooza-ossh/protego/internal/domain"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/funchooooza-ossh/protego/internal/services"
)

type LoginUsecase struct {
	userService  services.UserServiceInterface
	tokenService services.TokenServiceInterface
	maxAttempts  int
}

func NewLoginUsecase(
	userService services.UserServiceInterface,
	tokenService services.TokenServiceInterface,
	maxAttempts int,
) *LoginUsecase {
	return &LoginUsecase{
		userService:  userService,
		tokenService: tokenService,
		maxAttempts:  maxAttempts,
	}
}

func (u *LoginUsecase) Execute(ctx context.Context, email, password string) (string, string, error) {
	const origin = "login_usecase"

	user, err := u.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", e.ReturnErr(origin, err, e.Info)
	}

	userID := user.ID
	valid, err := u.userService.VerifyPassword(ctx, password, user.Password)
	if err != nil {
		return "", "", e.ReturnErr(origin, err, e.Warn)
	}
	if !valid {
		counter, err := u.userService.IncreaseCounter(ctx, userID)
		if err != nil {
			return "", "", e.ReturnErr(origin, err, e.Warn)
		}

		if counter >= u.maxAttempts {
			log.Printf("user %s has been blocked after %d failed attempts", userID, counter)

			// блокировку неважно логировать отдельно — best-effort
			e.BestEffort(origin, "BlockUser", u.userService.BlockUser(ctx, userID))

			return "", "", e.ReturnErr(
				origin,
				fmt.Errorf("%w:login", e.ErrTooManyRequests),
				e.Info,
			)
		}

		return "", "", e.ReturnErr(origin, fmt.Errorf("%w: credentials", e.ErrInvalidInput), e.Info)
	}

	e.BestEffort(origin, "DeleteCounter", u.userService.DeleteCounter(ctx, userID))

	claims := &domain.TokenClaims{
		UserID: userID,
		RoleID: user.Role.ID,
	}

	access, refresh, err := u.tokenService.CreatePair(ctx, claims)
	if err != nil {
		return "", "", e.ReturnErr(origin, err, e.Warn)
	}

	return access, refresh, nil
}
