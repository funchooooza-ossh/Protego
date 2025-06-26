package compositionInfrastructure

import (
	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/infra"
	m "github.com/funchooooza-ossh/protego/internal/metrics/common"
)

func newPasswordHasher(
	passwordCost int,
) contracts.PasswordHasherInterface {
	hasher := infra.NewBcryptHasher(passwordCost)
	return m.NewPasswordHasherWithMetrics(hasher)
}
