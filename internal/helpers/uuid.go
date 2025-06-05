package helpers

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDToPg(u uuid.UUID) pgtype.UUID {
	var pg pgtype.UUID
	_ = pg.Scan(u.String())
	return pg
}

func UUIDFromPg(p pgtype.UUID) uuid.UUID {
	u, err := uuid.FromBytes(p.Bytes[:])
	if err != nil {
		panic("invalid UUID in DB: " + err.Error())
	}
	return u
}
