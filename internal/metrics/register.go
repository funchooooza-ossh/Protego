package appmetrics

import (
	authmetrics "github.com/funchooooza-ossh/protego/internal/metrics/auth"
	lifespanmetrics "github.com/funchooooza-ossh/protego/internal/metrics/lifespan"
)

func Register() {
	authmetrics.Register()
	lifespanmetrics.Register()
}
