package adapters

import (
	"context"

	"github.com/funchooooza-ossh/protego/internal/db"
	"github.com/funchooooza-ossh/protego/internal/domain"
	"github.com/funchooooza-ossh/protego/internal/helpers"
	"github.com/funchooooza-ossh/protego/internal/mapper"
	"github.com/google/uuid"
)

type RoleRepository struct {
	q *db.UQueries
}

func NewRoleRepository(q *db.UQueries) *RoleRepository {
	return &RoleRepository{
		q: q,
	}
}

func (r *RoleRepository) Create(ctx context.Context, role *domain.Role) error {
	const origin = "role_repo.create"

	input := mapper.FromDomainRole(role)

	_, err := r.q.CreateRole(ctx, input)
	if err != nil {
		return ParseDBError(err, origin)
	}
	return nil
}

func (r *RoleRepository) GetByID(ctx context.Context, id string) (*domain.Role, error) {
	const origin = "role_repo.getById"

	uuidID := helpers.UUIDToPg(uuid.MustParse(id))
	roleRow, err := r.q.GetRoleByID(ctx, uuidID)
	if err != nil {
		return nil, ParseDBError(err, origin)
	}
	return mapper.ToDomainRole(roleRow), nil

}

func (r *RoleRepository) GetByCode(ctx context.Context, code string) (*domain.Role, error) {
	const origin = "role_repo.getByCode"

	roleRow, err := r.q.GetRoleByCode(ctx, code)
	if err != nil {
		return nil, ParseDBError(err, origin)
	}
	return mapper.ToDomainRole(roleRow), nil

}
