package compositionInfrastructure

import (
	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/infra"
)

func newJWTManager(secret string) contracts.JWTProvider {
	manager := infra.NewJWTManager(secret)
	return manager
}
