package repository

import (
	"context"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/internal/errors"
	synccredmap "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/sync_cred_map"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	noData = ""
)

type SessionRepository struct {
	// sessionID --> username
	rdb synccredmap.SyncCredentialsMap
}

func NewSessionRepository() *SessionRepository {
	res := &SessionRepository{
		rdb: *synccredmap.NewSyncCredentialsMap(),
	}

	return res
}

func (r *SessionRepository) GetSession(ctx context.Context, sessionID string) (string, error) {
	logger := log.Ctx(ctx)

	username, ok := r.rdb.Load(sessionID)
	if !ok {
		logger.Error().Err(errors.Wrap(errs.ErrSessionNotExists, errs.ErrMsgFailedToGetSession)).Msg(errs.ErrMsgSessionNotExists)
		return noData, errs.ErrSessionNotExists
	}
	return username, nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	logger := log.Ctx(ctx)

	_, ok := r.rdb.Load(sessionID)
	if !ok {
		logger.Error().Err(errors.Wrap(errs.ErrSessionNotExists, errs.ErrMsgFailedToGetSession)).Msg(errs.ErrMsgSessionNotExists)
		return errs.ErrSessionNotExists
	}

	r.rdb.Delete(sessionID)
	return nil
}

func (r *SessionRepository) StoreSession(ctx context.Context, newSessionID, login string) error {
	r.rdb.Store(newSessionID, login)
	return nil
}
