package appmetrics

import (
	authmetrics "github.com/funchooooza-ossh/protego/internal/metrics/auth"
	commonmetrics "github.com/funchooooza-ossh/protego/internal/metrics/common"
	lifespanmetrics "github.com/funchooooza-ossh/protego/internal/metrics/lifespan"
	loginmetrics "github.com/funchooooza-ossh/protego/internal/metrics/login"
)

func Register() {
	authmetrics.Register()
	lifespanmetrics.Register()
	loginmetrics.Register()
	commonmetrics.Register()
}
