package repository

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/common"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	synccredmap "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/sync_cred_map"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type UserRepository struct {
	// username --> hashedPass
	users synccredmap.SyncCredentialsMap
	cfg   *config.Cookie
}

func NewUserRepository() *UserRepository {
	res := &UserRepository{
		users: *synccredmap.NewSyncCredentialsMap(),
	}

	return res
}

func (r *UserRepository) CreateUser(ctx context.Context, login, hashedPass string) error {
	logger := log.Ctx(ctx)

	if _, exists := r.users.Load(login); exists {
		logger.Error().Err(errors.New(errs.ErrAlreadyExists)).Msg(common.MsgUserWithNameAlreadyExists)
		return errors.New(errs.ErrAlreadyExists)
	}
	r.users.Store(login, string(hashedPass))
	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, login string) (hashedPass string, errRepo error) {
	logger := log.Ctx(ctx)

	hashedPass, exists := r.users.Load(login)
	if exists {
		return hashedPass, nil
	}
	err := errors.New(errs.ErrIncorrectLogin)
	logger.Error().Err(errors.Wrap(err, errs.ErrIncorrectLoginOrPassword)).Msg(err.Error())
	return noData, err
}
