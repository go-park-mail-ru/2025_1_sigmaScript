package repository

import (
	"context"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/pkg/errors"
)

func (r *UserRepository) GetUser(ctx context.Context, login string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.rdb[login]
	if !ok {
		return nil, errors.New(errs.ErrIncorrectLogin)
	}

	return user, nil
}
