package compositionInfrastructure

import (
	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/infra"
	loginmetrics "github.com/funchooooza-ossh/protego/internal/metrics/login"
)

func newPasswordHasher(
	passwordCost int,
) contracts.PasswordHasherInterface {
	adapter := infra.NewBcryptHasher(passwordCost)
	return loginmetrics.NewPasswordHasherWithMetrics(adapter)
}
