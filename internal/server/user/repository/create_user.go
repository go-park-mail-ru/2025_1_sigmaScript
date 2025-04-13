package repository

import (
	"context"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/pkg/errors"
)

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.rdb[user.Username]; ok {
		return errors.New(errs.ErrAlreadyExists)
	}

	r.rdb[user.Username] = user

	return nil
}
