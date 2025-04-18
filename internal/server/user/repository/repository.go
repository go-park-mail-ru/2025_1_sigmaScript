package repository

import (
	"database/sql"
	"sync"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
)

type UserRepository struct {
	mu   sync.RWMutex
	rdb  map[string]*models.User
	pgdb *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		rdb:  make(map[string]*models.User),
		pgdb: db,
	}
}
