package loginmetrics

import (
	"context"
	"errors"
	"time"

	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/domain"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/funchooooza-ossh/protego/internal/helpers"
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
	helpers.ObserveDuration(UserRepoCreateLastDelay, UserRepoCreateSummary, start)

	if err != nil && errors.Is(err, e.ErrAlreadyExists) {
		UserCreateAlreadyExists.Inc()
	}
	return err
}

func (m *UserRepositoryWithMetrics) GetByID(ctx context.Context, id string) (*domain.User, error) {
	UserRepoGetByIDCalls.Inc()
	start := time.Now()
	user, err := m.impl.GetByID(ctx, id)
	helpers.ObserveDuration(UserRepoGetByIDLastDelay, UserRepoGetByIDSummary, start)

	if err != nil && errors.Is(err, e.ErrNotFound) {
		UserGetByIDNotFound.Inc()
	}
	return user, err
}

func (m *UserRepositoryWithMetrics) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	UserRepoGetByEmailCalls.Inc()
	start := time.Now()
	user, err := m.impl.GetByEmail(ctx, email)
	helpers.ObserveDuration(UserRepoGetByEmailLastDelay, UserRepoGetByEmailSummary, start)

	if err != nil && errors.Is(err, e.ErrNotFound) {
		UserGetByEmailNotFound.Inc()
	}
	return user, err
}

func (m *UserRepositoryWithMetrics) Update(ctx context.Context, user *domain.User) error {
	UserRepoUpdateCalls.Inc()
	start := time.Now()
	err := m.impl.Update(ctx, user)
	helpers.ObserveDuration(UserRepoUpdateLastDelay, UserRepoUpdateSummary, start)

	if err != nil && errors.Is(err, e.ErrNotFound) {
		UserUpdateNotFound.Inc()
	}
	return err
}

func (m *UserRepositoryWithMetrics) Delete(ctx context.Context, id string) error {
	UserRepoDeleteCalls.Inc()
	start := time.Now()
	err := m.impl.Delete(ctx, id)
	helpers.ObserveDuration(UserRepoDeleteLastDelay, UserRepoDeleteSummary, start)

	if err != nil && errors.Is(err, e.ErrNotFound) {
		UserDeleteNotFound.Inc()
	}
	return err
}

type CounterRepositoryWithMetrics struct {
	impl contracts.CounterRepositoryInterface
}

func NewCounterRepositoryWithMetrics(impl contracts.CounterRepositoryInterface) *CounterRepositoryWithMetrics {
	return &CounterRepositoryWithMetrics{
		impl: impl,
	}
}

func (m *CounterRepositoryWithMetrics) Increment(ctx context.Context, key string) (int, error) {
	CounterIncrementCalls.Inc()
	start := time.Now()
	val, err := m.impl.Increment(ctx, key)
	helpers.ObserveDuration(CounterIncrementLastDelay, CounterIncrementSummary, start)

	return val, err
}

func (m *CounterRepositoryWithMetrics) Delete(ctx context.Context, key string) error {
	CounterDeleteCalls.Inc()
	start := time.Now()
	err := m.impl.Delete(ctx, key)
	helpers.ObserveDuration(CounterDeleteLastDelay, CounterDeleteSummary, start)

	return err
}
