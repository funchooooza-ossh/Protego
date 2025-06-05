package adapters

import (
	"context"

	"github.com/funchooooza-ossh/protego/internal/db"
	"github.com/funchooooza-ossh/protego/internal/domain"
	"github.com/funchooooza-ossh/protego/internal/mapper"
)

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) *UserRepository {
	return &UserRepository{
		q: q,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	const origin = "user_repo.create"

	input := mapper.FromDomainUser(user)

	_, err := r.q.CreateUser(ctx, input) // created user not needed
	if err != nil {
		return ParseDBError(err, origin)
	}

	return nil
}
