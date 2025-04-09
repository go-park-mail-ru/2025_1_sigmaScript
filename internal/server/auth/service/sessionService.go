package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/session"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	noData = ""
)

type SessionRepositoryInterface interface {
	StoreSession(ctx context.Context, newSessionID, login string) error
	DeleteSession(ctx context.Context, sessionID string) error
	GetSession(ctx context.Context, sessionID string) (string, error)
}

type SessionService struct {
	sessionRepo   SessionRepositoryInterface
	sessionLength int
}

func NewSessionService(ctx context.Context, sessionRepo SessionRepositoryInterface) *SessionService {
	return &SessionService{
		sessionRepo:   sessionRepo,
		sessionLength: config.FromCookieContext(ctx).SessionLength,
	}
}

// CreateSession method creates new sessionID and stores username by sessionID
func (s *SessionService) CreateSession(ctx context.Context, username string) (string, error) {
	logger := log.Ctx(ctx)

	newSessionID, err := session.GenerateSessionID(s.sessionLength)
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrMsgGenerateSession)).Msg(errors.Wrap(err, errs.ErrMsgGenerateSession).Error())
		return noData, errs.ErrGenerateSession
	}

	errRepo := s.sessionRepo.StoreSession(ctx, newSessionID, username)
	if errRepo != nil {
		logger.Error().Err(errRepo).Msg(errRepo.Error())
		return noData, errRepo
	}
	logger.Info().Msg("Session created")
	return newSessionID, nil
}

// DeleteSession method deletes session by sessionID
func (s *SessionService) DeleteSession(ctx context.Context, sessionID string) error {
	logger := log.Ctx(ctx)

	errRepo := s.sessionRepo.DeleteSession(ctx, sessionID)
	if errRepo != nil {
		logger.Error().Err(errRepo).Msg(errRepo.Error())
		return errRepo
	}

	logger.Info().Msg("session successfully deleted")

	return nil
}

// GetSession method gets session by sessionID
func (s *SessionService) GetSession(ctx context.Context, sessionID string) (string, error) {
	logger := log.Ctx(ctx)

	username, errRepo := s.sessionRepo.GetSession(ctx, sessionID)
	if errRepo != nil {
		logger.Error().Err(errors.Wrap(errRepo, errs.ErrMsgSessionNotExists)).Msg(errRepo.Error())
		return noData, errRepo
	}

	return username, nil
}
