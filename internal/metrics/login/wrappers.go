package loginmetrics

import (
	"context"
	"errors"
	"time"

	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/domain"
	e "github.com/funchooooza-ossh/protego/internal/errors"
)

type UserRepositoryWithMetrics struct {
	impl contracts.UserRepositoryInterface
}

func NewUserRepositoryWithMetrics(impl contracts.UserRepositoryInterface) *UserRepositoryWithMetrics {
	return &UserRepositoryWithMetrics{
		impl: impl,
	}
}

func (m *UserRepositoryWithMetrics) Create(ctx context.Context, user *domain.User) error {
	UserRepoCreateCalls.Inc()

	start := time.Now()
	err := m.impl.Create(ctx, user)
	duration := time.Since(start).Seconds()

	UserRepoCreateLastDelay.Set(duration)
	UserRepoCreateSummary.Observe(duration)

	if err != nil && errors.Is(err, e.ErrAlreadyExists) {
		UserCreateAlreadyExists.Inc()
	}
	return err
}

func (m *UserRepositoryWithMetrics) GetByID(ctx context.Context, id string) (*domain.User, error) {
	UserRepoGetByIDCalls.Inc()

	start := time.Now()
	user, err := m.impl.GetByID(ctx, id)
	duration := time.Since(start).Seconds()

	UserRepoGetByIDLastDelay.Set(duration)
	UserRepoGetByIDSummary.Observe(duration)

	if err != nil && errors.Is(err, e.ErrNotFound) {
		UserGetByIDNotFound.Inc()
	}

	return user, err

}

func (m *UserRepositoryWithMetrics) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	UserRepoGetByEmailCalls.Inc()

	start := time.Now()
	user, err := m.impl.GetByEmail(ctx, email)
	duration := time.Since(start).Seconds()

	UserRepoGetByEmailLastDelay.Set(duration)
	UserRepoGetByEmailSummary.Observe(duration)

	if err != nil && errors.Is(err, e.ErrNotFound) {
		UserGetByEmailNotFound.Inc()
	}

	return user, err

}

func (m *UserRepositoryWithMetrics) Update(ctx context.Context, user *domain.User) error {
	UserRepoUpdateCalls.Inc()

	start := time.Now()
	err := m.impl.Update(ctx, user)
	duration := time.Since(start).Seconds()

	UserRepoUpdateLastDelay.Set(duration)
	UserRepoUpdateSummary.Observe(duration)

	if err != nil && errors.Is(err, e.ErrNotFound) {
		UserUpdateNotFound.Inc()
	}
	return err
}

func (m *UserRepositoryWithMetrics) Delete(ctx context.Context, id string) error {
	UserRepoDeleteCalls.Inc()

	start := time.Now()
	err := m.impl.Delete(ctx, id)
	duration := time.Since(start).Seconds()

	UserRepoDeleteLastDelay.Set(duration)
	UserRepoDeleteSummary.Observe(duration)

	if err != nil && errors.Is(err, e.ErrNotFound) {
		UserDeleteNotFound.Inc()
	}
	return err

}
