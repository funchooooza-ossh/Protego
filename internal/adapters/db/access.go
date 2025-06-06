package adapters

import (
	"context"

	"github.com/funchooooza-ossh/protego/internal/db"
	"github.com/funchooooza-ossh/protego/internal/helpers"
	"github.com/google/uuid"
)

type AccessRepository struct {
	q *db.UQueries
}

func NewAccessRepository(q *db.UQueries) *AccessRepository {
	return &AccessRepository{
		q: q,
	}
}

func (r *AccessRepository) HasAccess(ctx context.Context, RoleID, action, resourceCode string) (bool, error) {
	const origin = "access_repo"

	uuidID := helpers.UUIDToPg(uuid.MustParse(RoleID))

	allowed, err := r.q.HasAccess(ctx, uuidID, action, resourceCode)

	if err != nil {
		return false, ParseDBError(err, origin)
	}

	return allowed, nil

}
