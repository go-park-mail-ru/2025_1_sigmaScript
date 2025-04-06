package repository

import (
	"context"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	synccredmap "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/sync_cred_map"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	noData = ""
)

type AuthRepositoryInterface interface {
	CreateUser(ctx context.Context, login, hashedPass string) *errs.RepoError
	GetUser(ctx context.Context, login string) (hashedPass string, errRepo *errs.RepoError)
	StoreSession(ctx context.Context, newSessionID, login string) *errs.RepoError
	DeleteSession(ctx context.Context, sessionID string) *errs.RepoError
	GetSession(ctx context.Context, sessionID string) (string, *errs.RepoError)
}

type AuthRepository struct {
	// username --> hashedPass
	users synccredmap.SyncCredentialsMap
	// sessionID --> username
	sessions synccredmap.SyncCredentialsMap
}

func NewAuthRepository() AuthRepositoryInterface {
	res := &AuthRepository{
		users:    *synccredmap.NewSyncCredentialsMap(),
		sessions: *synccredmap.NewSyncCredentialsMap(),
	}

	return res
}

func (ar *AuthRepository) GetSession(ctx context.Context, sessionID string) (string, *errs.RepoError) {
	logger := log.Ctx(ctx)
	// repo
	username, ok := ar.sessions.Load(sessionID)
	if !ok {
		err := errors.New("failed to get session")
		logger.Error().Err(errors.Wrap(err, errs.ErrSessionNotExists)).Msg(errors.Wrap(err, errs.ErrSessionNotExists).Error())
		return noData, &errs.RepoError{
			Msg:   err.Error(),
			Error: errors.Wrap(err, errs.ErrSessionNotExists),
		}
	}

	return username, nil
}

func (ar *AuthRepository) DeleteSession(ctx context.Context, sessionID string) *errs.RepoError {
	ar.sessions.Delete(sessionID)

	return nil
}

func (ar *AuthRepository) StoreSession(ctx context.Context, newSessionID, login string) *errs.RepoError {
	ar.sessions.Store(newSessionID, login)

	return nil
}

func (ar *AuthRepository) CreateUser(ctx context.Context, login, hashedPass string) *errs.RepoError {
	logger := log.Ctx(ctx)

	if _, exists := ar.users.Load(login); exists {
		msg := "user with that name already exists"
		logger.Error().Err(errors.New(errs.ErrAlreadyExists)).Msg(msg)
		return &errs.RepoError{
			Msg:   msg,
			Error: errors.New(errs.ErrAlreadyExists),
		}
	}

	ar.users.Store(login, string(hashedPass))

	return nil
}

func (ar *AuthRepository) GetUser(ctx context.Context, login string) (hashedPass string, errRepo *errs.RepoError) {
	logger := log.Ctx(ctx)

	// repo login
	hashedPass, exists := ar.users.Load(login)
	if !exists {
		errMsg := errors.New(errs.ErrIncorrectLogin)
		logger.Error().Err(errors.Wrap(errMsg, errs.ErrIncorrectLoginOrPassword)).Msg(errMsg.Error())
		return noData, &errs.RepoError{
			Msg:   errMsg.Error(),
			Error: errors.Wrap(errMsg, errs.ErrIncorrectLoginOrPasswordShort),
		}
	}
	return hashedPass, nil
}
