package repository

import (
	"context"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/pkg/errors"
)

func (r *UserRepository) DeleteUser(ctx context.Context, login string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.rdb[login]; ok {
		delete(r.rdb, login)
		return nil
	}

	return errors.New(errs.ErrIncorrectLogin)
}
