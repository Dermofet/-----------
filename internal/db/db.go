package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	QueryTimeout = 10 * time.Second
)

type source struct {
	db *sqlx.DB
}

func NewSource(db *sqlx.DB) *source {
	return &source{
		db: db,
	}
}
