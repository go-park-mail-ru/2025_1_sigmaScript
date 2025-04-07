package repository

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	synccredmap "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/sync_cred_map"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type SessionRepository struct {
	// sessionID --> username
	sessions synccredmap.SyncCredentialsMap
	cfg      *config.Cookie
}

func NewSessionRepository() *SessionRepository {
	res := &SessionRepository{
		sessions: *synccredmap.NewSyncCredentialsMap(),
	}

	return res
}

func (r *SessionRepository) GetSession(ctx context.Context, sessionID string) (string, error) {
	logger := log.Ctx(ctx)

	username, ok := r.sessions.Load(sessionID)
	if !ok {
		logger.Error().Err(errors.Wrap(errs.ErrSessionNotExists, errs.ErrMsgFailedToGetSession)).Msg(errs.ErrMsgSessionNotExists)
		return noData, errs.ErrSessionNotExists
	}
	return username, nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	logger := log.Ctx(ctx)

	_, ok := r.sessions.Load(sessionID)
	if !ok {
		logger.Error().Err(errors.Wrap(errs.ErrSessionNotExists, errs.ErrMsgFailedToGetSession)).Msg(errs.ErrMsgSessionNotExists)
		return errs.ErrSessionNotExists
	}

	r.sessions.Delete(sessionID)
	return nil
}

func (r *SessionRepository) StoreSession(ctx context.Context, newSessionID, login string) error {
	r.sessions.Store(newSessionID, login)
	return nil
}
