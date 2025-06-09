package appmetrics

import (
	authmetrics "github.com/funchooooza-ossh/protego/internal/metrics/auth"
	lifespanmetrics "github.com/funchooooza-ossh/protego/internal/metrics/lifespan"
	loginmetrics "github.com/funchooooza-ossh/protego/internal/metrics/login"
)

func Register() {
	authmetrics.Register()
	lifespanmetrics.Register()
	loginmetrics.Register()
}
