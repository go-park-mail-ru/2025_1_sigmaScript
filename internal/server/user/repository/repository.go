package repository

import (
	"sync"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
)

type UserRepository struct {
	mu  sync.RWMutex
	rdb map[string]*models.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		rdb: make(map[string]*models.User),
	}
}
