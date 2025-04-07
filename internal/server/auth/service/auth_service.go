package service

// import (
// 	"context"

// 	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
// 	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
// 	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
// 	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/session"
// 	"github.com/pkg/errors"
// 	"github.com/rs/zerolog/log"
// 	"golang.org/x/crypto/bcrypt"
// )

// // type UserRepositoryInterface interface {
// // 	CreateUser(ctx context.Context, login, hashedPass string) error
// // 	GetUser(ctx context.Context, login string) (hashedPass string, errRepo error)
// // }

// // type SessionRepositoryInterface interface {
// // 	StoreSession(ctx context.Context, newSessionID, login string) error
// // 	DeleteSession(ctx context.Context, sessionID string) error
// // 	GetSession(ctx context.Context, sessionID string) (string, error)
// // }

// type AuthService struct {
// 	userRepo      UserRepositoryInterface
// 	sessionRepo   SessionRepositoryInterface
// 	sessionLength int
// }

// func NewAuthService(ctx context.Context, userRepo UserRepositoryInterface,
// 	sessionRepo SessionRepositoryInterface) *AuthService {
// 	return &AuthService{
// 		userRepo:      userRepo,
// 		sessionRepo:   sessionRepo,
// 		sessionLength: config.FromCookieContext(ctx).SessionLength,
// 	}
// }

// // Register method registers user with given parameters
// func (s *AuthService) Register(ctx context.Context, regUser models.RegisterData) error {
// 	logger := log.Ctx(ctx)

// 	hashedPass, err := bcrypt.GenerateFromPassword([]byte(regUser.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		logger.Error().Err(errors.Wrap(err, errs.ErrBcrypt)).Msg(errs.ErrInvalidPassword)
// 		return errors.New(errs.ErrInvalidPassword)
// 	}

// 	errRepo := s.userRepo.CreateUser(ctx, regUser.Username, string(hashedPass))
// 	if errRepo != nil {
// 		logger.Error().Err(errRepo).Msg(errRepo.Error())
// 		return errRepo
// 	}
// 	return nil
// }

// // CreateSession method creates new sessionID and stores username by sessionID
// func (s *AuthService) CreateSession(ctx context.Context, username string) (string, error) {
// 	logger := log.Ctx(ctx)

// 	newSessionID, err := session.GenerateSessionID(s.sessionLength)
// 	if err != nil {
// 		logger.Error().Err(errors.Wrap(err, errs.ErrMsgGenerateSession)).Msg(errors.Wrap(err, errs.ErrMsgGenerateSession).Error())
// 		return noData, errs.ErrGenerateSession
// 	}

// 	errRepo := s.sessionRepo.StoreSession(ctx, newSessionID, username)
// 	if errRepo != nil {
// 		logger.Error().Err(errRepo).Msg(errRepo.Error())
// 		return noData, errRepo
// 	}
// 	logger.Info().Msg("Session created")
// 	return newSessionID, nil
// }

// // DeleteSession method deletes session by sessionID
// func (s *AuthService) DeleteSession(ctx context.Context, sessionID string) error {
// 	logger := log.Ctx(ctx)

// 	errRepo := s.sessionRepo.DeleteSession(ctx, sessionID)
// 	if errRepo != nil {
// 		logger.Error().Err(errRepo).Msg(errRepo.Error())
// 		return errRepo
// 	}

// 	logger.Info().Msg("session successfully deleted")

// 	return nil
// }

// // GetSession method gets session by sessionID
// func (s *AuthService) GetSession(ctx context.Context, sessionID string) (string, error) {
// 	logger := log.Ctx(ctx)

// 	username, errRepo := s.sessionRepo.GetSession(ctx, sessionID)
// 	if errRepo != nil {
// 		logger.Error().Err(errors.Wrap(errRepo, errs.ErrMsgSessionNotExists)).Msg(errRepo.Error())
// 		return noData, errRepo
// 	}

// 	return username, nil
// }

// // Login method checks if user with given credentials exists
// func (s *AuthService) Login(ctx context.Context, loginData models.LoginData) error {
// 	logger := log.Ctx(ctx)

// 	hashedPass, errRepo := s.userRepo.GetUser(ctx, loginData.Username)
// 	if errRepo != nil {
// 		logger.Error().Err(errors.Wrap(errRepo, errs.ErrIncorrectLoginOrPassword)).Msg(errRepo.Error())
// 		return errRepo
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(loginData.Password)); err != nil {
// 		logger.Error().Err(errors.Wrap(err, errs.ErrIncorrectLoginOrPassword)).Msg(errs.ErrIncorrectPassword)
// 		return errors.New(errs.ErrIncorrectPassword)
// 	}

// 	return nil
// }
