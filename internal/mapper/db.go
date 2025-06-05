package mapper

import (
	"github.com/funchooooza-ossh/protego/internal/db"
	"github.com/funchooooza-ossh/protego/internal/domain"
	"github.com/funchooooza-ossh/protego/internal/helpers"
	"github.com/google/uuid"
)

func ToDomainUser(row db.User) *domain.User {
	uid := helpers.UUIDFromPg(row.ID)
	roleID := helpers.UUIDFromPg(row.RoleID)

	return &domain.User{
		ID:       uid.String(),
		Email:    row.Email,
		Password: row.Password,
		Blocked:  row.Blocked,
		Role: domain.Role{
			ID: roleID.String(),
		},
	}
}

func FromDomainUser(u *domain.User) db.UserParams {
	return db.UserParams{
		ID:       helpers.UUIDToPg(uuid.MustParse(u.ID)),
		Email:    u.Email,
		Password: u.Password,
		RoleID:   helpers.UUIDToPg(uuid.MustParse(u.Role.ID)),
		Blocked:  u.Blocked,
	}
}
