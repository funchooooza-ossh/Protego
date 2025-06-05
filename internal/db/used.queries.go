package db

import (
	"github.com/jackc/pgx/v5"
)

type UQueries struct {
	db DBTX
}

func (q *UQueries) WithTx(tx pgx.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}

func NewQueries(db DBTX) *UQueries {
	return &UQueries{db: db}
}
