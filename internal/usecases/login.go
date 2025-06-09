package usecases

import (
	"context"
	"fmt"

	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/domain"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/funchooooza-ossh/protego/internal/logger"
	"go.uber.org/zap/zapcore"
)

type LoginUsecase struct {
	userService  contracts.UserServiceInterface
	tokenService contracts.TokenServiceInterface
	maxAttempts  int
}

func NewLoginUsecase(
	userService contracts.UserServiceInterface,
	tokenService contracts.TokenServiceInterface,
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

	var loginCounter = 0

	user, err := u.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", e.ReturnErr(ctx, origin, err, e.Info)
	}

	userID := user.ID
	valid, err := u.userService.VerifyPassword(ctx, password, user.Password)
	if err != nil {
		e.BestEffort(ctx, origin, "verify password", err)
	}
	if !valid {
		loginCounter, err := u.userService.IncreaseCounter(ctx, userID)
		if err != nil {
			return "", "", e.ReturnErr(ctx, origin, err, e.Warn)
		}

		if loginCounter >= u.maxAttempts {
			logger.Log(ctx, zapcore.InfoLevel, fmt.Sprintf("user %s has been blocked after %d failed attempts", userID, loginCounter))

			// блокировку неважно логировать отдельно — best-effort
			e.BestEffort(ctx, origin, "BlockUser", u.userService.BlockUser(ctx, userID))

			return "", "", e.ReturnErr(
				ctx,
				origin,
				fmt.Errorf("%w:login", e.ErrTooManyRequests),
				e.Info,
			)
		}

		return "", "", e.ReturnErr(ctx, origin, fmt.Errorf("%w: credentials", e.ErrInvalidInput), e.Info)
	}
	if loginCounter > 0 {
		e.BestEffort(ctx, origin, "DeleteCounter", u.userService.DeleteCounter(ctx, userID))

	}

	claims := &domain.TokenClaims{
		UserID: userID,
		RoleID: user.Role.ID,
	}

	access, refresh, err := u.tokenService.CreatePair(ctx, claims)
	if err != nil {
		return "", "", e.ReturnErr(ctx, origin, err, e.Warn)
	}

	return access, refresh, nil
}
