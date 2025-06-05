package mapper

import (
	"github.com/funchooooza-ossh/protego/internal/db"
	"github.com/funchooooza-ossh/protego/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ToDomainUser(row db.User, roleCode string) *domain.User {
	uid := uuidFromPg(row.ID)
	roleID := uuidFromPg(row.RoleID)

	return &domain.User{
		ID:       uid.String(),
		Email:    row.Email,
		Password: row.Password,
		Blocked:  row.Blocked,
		Role: domain.Role{
			ID:   roleID.String(),
			Code: roleCode,
		},
	}
}

func FromDomainUser(u *domain.User) db.CreateUserParams {
	return db.CreateUserParams{
		ID:       uuidToPg(uuid.MustParse(u.ID)),
		Email:    u.Email,
		Password: u.Password,
		RoleID:   uuidToPg(uuid.MustParse(u.Role.ID)),
		Blocked:  u.Blocked,
	}
}

func uuidToPg(u uuid.UUID) pgtype.UUID {
	var pg pgtype.UUID
	_ = pg.Scan(u.String())
	return pg
}

func uuidFromPg(p pgtype.UUID) uuid.UUID {
	u, err := uuid.FromBytes(p.Bytes[:])
	if err != nil {
		panic("invalid UUID in DB: " + err.Error())
	}
	return u
}
