package repository

import (
	"database/sql"
)

type UserRepository struct {
	pgdb *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		pgdb: db,
	}
}
