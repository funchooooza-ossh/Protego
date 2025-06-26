package commonmetrics

import (
	"context"
	"time"

	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// HASHER
	VerifyPasswordDelay = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "verify_password_last_delay_seconds",
		Help:      "Last VerifyPassword delay in seconds",
		Namespace: "auth",
		Subsystem: "password",
	})
	HashPasswordDelay = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "has_password_last_delay_seconds",
		Help:      "Last hashPassword delay in seconds",
		Namespace: "auth",
		Subsystem: "password",
	})
	VerifyPasswordSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "verify_password_delay_seconds",
		Help:       "Summary of VerifyPassword delays",
		Namespace:  "auth",
		Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001},
		Subsystem:  "password",
	})
	HashPasswordSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "hash_password_delay_seconds",
		Help:       "Summary of hashPassword delays",
		Namespace:  "auth",
		Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001},
		Subsystem:  "password",
	})
)

func registerHash() {
	prometheus.MustRegister(
		VerifyPasswordDelay, HashPasswordDelay,
		VerifyPasswordSummary, HashPasswordSummary,
	)
}

type PasswordHasherWithMetrics struct {
	impl contracts.PasswordHasherInterface
}

func NewPasswordHasherWithMetrics(impl contracts.PasswordHasherInterface) *PasswordHasherWithMetrics {
	return &PasswordHasherWithMetrics{
		impl: impl,
	}
}

func (m *PasswordHasherWithMetrics) Hash(ctx context.Context, unhashed string) ([]byte, error) {
	start := time.Now()
	val, err := m.impl.Hash(ctx, unhashed)

	helpers.ObserveDuration(HashPasswordDelay, HashPasswordSummary, start)
	return val, err
}

func (m *PasswordHasherWithMetrics) Verify(ctx context.Context, unhashed, hashed string) (bool, error) {
	start := time.Now()
	val, err := m.impl.Verify(ctx, unhashed, hashed)

	helpers.ObserveDuration(VerifyPasswordDelay, VerifyPasswordSummary, start)
	return val, err
}
