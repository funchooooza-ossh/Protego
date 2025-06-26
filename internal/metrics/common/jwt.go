package commonmetrics

import (
	"context"
	"errors"
	"time"

	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/domain"
	apperrors "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/funchooooza-ossh/protego/internal/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	JWTGenerateTokenDelay = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "generate_token_last_delay_seconds",
		Help:      "Last delay of generate jwt token in seconds",
		Namespace: "auth",
		Subsystem: "jwt",
	})
	JWTVerifyTokenDelay = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "verify_token_last_delay_seconds",
		Help:      "Last delay of veryfing jwt token in seconds",
		Namespace: "auth",
		Subsystem: "jwt",
	})
	JWTGenerateTokenSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "generate_token_summary",
		Help:       "Summary of generate jwt token delays",
		Namespace:  "auth",
		Subsystem:  "jwt",
		Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001},
	})
	JWTVerifyTokenSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "verify_token_summary",
		Help:       "Summary of veryfing jwt token delays",
		Namespace:  "auth",
		Subsystem:  "jwt",
		Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001},
	})

	JWTGenerateTokenErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "generate_token_errors_count",
		Help:      "Total number of generate jwt errors",
		Namespace: "auth",
		Subsystem: "jwt",
	})
	JWTInvalidToken = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "verify_token_failed_count",
		Help:      "Total number of verify token fails cause of invalid token",
		Namespace: "auth",
		Subsystem: "jwt",
	})
	JWTVerifyTokenErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "verify_token_errors_count",
		Help:      "Total number of verify token errors",
		Namespace: "auth",
		Subsystem: "jwt",
	})
)

func registerJWT() {
	prometheus.MustRegister(
		JWTGenerateTokenErrors, JWTInvalidToken, JWTVerifyTokenErrors,
		JWTGenerateTokenDelay, JWTVerifyTokenDelay,
		JWTGenerateTokenSummary, JWTVerifyTokenSummary,
	)
}

type JWTManagerWithMetrics struct {
	impl contracts.JWTProvider
}

func NewJWTManagerWithMetrics(impl contracts.JWTProvider) *JWTManagerWithMetrics {
	return &JWTManagerWithMetrics{
		impl: impl,
	}
}

func (m *JWTManagerWithMetrics) GenerateToken(ctx context.Context, claims *domain.TokenClaims) (string, error) {
	start := time.Now()
	val, err := m.impl.GenerateToken(ctx, claims)
	helpers.ObserveDuration(JWTGenerateTokenDelay, JWTGenerateTokenSummary, start)
	if err != nil {
		JWTGenerateTokenErrors.Inc()
	}
	return val, err
}

func (m *JWTManagerWithMetrics) VerifyToken(ctx context.Context, tokenStr string) (*domain.TokenClaims, error) {
	start := time.Now()

	val, err := m.impl.VerifyToken(ctx, tokenStr)
	helpers.ObserveDuration(JWTVerifyTokenDelay, JWTVerifyTokenSummary, start)
	if err != nil {
		if errors.Is(err, apperrors.ErrTokenInvalid) {
			JWTInvalidToken.Inc()
			return val, err
		}
		JWTVerifyTokenErrors.Inc()

	}
	return val, err
}
