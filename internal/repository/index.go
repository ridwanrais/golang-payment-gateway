package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

type Repository interface {
	// account
}

func NewRepository(d *pgxpool.Pool) Repository {
	return &repository{db: d}
}
