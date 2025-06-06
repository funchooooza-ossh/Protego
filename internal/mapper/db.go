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

func ToDomainRole(row db.Role) *domain.Role {
	uid := helpers.UUIDFromPg(row.ID)

	return &domain.Role{
		ID:   uid.String(),
		Code: row.Code,
	}
}

func FromDomainRole(r *domain.Role) db.RoleParams {
	return db.RoleParams{
		ID:   helpers.UUIDToPg(uuid.MustParse(r.ID)),
		Code: r.Code,
	}
}
