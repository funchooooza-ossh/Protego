package httpHelpers

import (
	"strings"

	"github.com/funchooooza-ossh/protego/cmd/server/dto"
	"github.com/funchooooza-ossh/protego/internal/domain"
)

func ExtractParams(req dto.AuthRequest) string {
	action := domain.MapMethodToAction(strings.ToLower(req.Method))
	return string(action)
}
