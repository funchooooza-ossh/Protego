package adapters

import (
	"context"

	"github.com/funchooooza-ossh/protego/internal/db"
	"github.com/funchooooza-ossh/protego/internal/domain"
	"github.com/funchooooza-ossh/protego/internal/helpers"
	"github.com/funchooooza-ossh/protego/internal/mapper"
	"github.com/google/uuid"
)

type UserRepository struct {
	q *db.UQueries
}

func NewUserRepository(q *db.UQueries) *UserRepository {
	return &UserRepository{
		q: q,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	const origin = "user_repo.create"

	input := mapper.FromDomainUser(user)

	_, err := r.q.CreateUser(ctx, input)
	if err != nil {
		return ParseDBError(ctx, err, origin)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	const origin = "user_repo.getById"

	uuidID := helpers.UUIDToPg(uuid.MustParse(id))
	userRow, err := r.q.GetUserByID(ctx, uuidID)
	if err != nil {
		return nil, ParseDBError(ctx, err, origin)
	}

	return mapper.ToDomainUser(userRow), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const origin = "user_repo.getByEmail"

	userRow, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ParseDBError(ctx, err, origin)
	}
	return mapper.ToDomainUser(userRow), nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	const origin = "user_repo.update"

	input := mapper.FromDomainUser(user)
	_, err := r.q.UpdateUser(ctx, input)
	if err != nil {
		return ParseDBError(ctx, err, origin)
	}

	return nil

}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	const origin = "user_repo.delete"

	uuidID := helpers.UUIDToPg(uuid.MustParse(id))

	if err := r.q.DeleteUser(ctx, uuidID); err != nil {
		return ParseDBError(ctx, err, origin)
	}
	return nil
}
