package memory

import (
	"context"
	"errors"
	"sync"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
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

func (r *UserRepository) GetUser(ctx context.Context, login string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.rdb[login]
	if !ok {
		return nil, errors.New(errs.ErrIncorrectLogin)
	}

	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.rdb[user.Username]; ok {
		return errors.New(errs.ErrAlreadyExists)
	}

	r.rdb[user.Username] = user

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, login string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.rdb[login]; ok {
		delete(r.rdb, login)
		return nil
	}

	return errors.New(errs.ErrIncorrectLogin)
}
