package auth

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	auth "github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/api/auth_api_v1/proto"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/config"
	delivery "github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/internal/auth/delivery"
	repoAuthSessions "github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/internal/auth/repository"
	service "github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/internal/auth/service"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/internal/interceptors"
)

// AppAuth is a root struct of auth_service
type AppAuth struct {
	logger *zerolog.Logger
	srv    *grpc.Server
	cfg    *config.Config
}

// New returns an instance of AppAuth
func New(isTest bool) (*AppAuth, error) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("error initialize app cfg: %w", err)
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.ChainUnaryInterceptors(
			interceptors.LoggerInterceptor,
			interceptors.AccessLogInterceptor,
		)),
	)

	reflection.Register(srv)

	sessRepo := repoAuthSessions.NewSessionRepository()

	sessServ := service.NewSessionService(sessRepo, config.SessionLength)
	auth.RegisterSessionRPCServer(srv, delivery.NewAuthServiceGRPCHandler(sessServ))

	return &AppAuth{
		srv:    srv,
		logger: &logger,
		cfg:    cfg,
	}, nil
}

// Run starts grpc server
func (a *AppAuth) Run() {
	lis, err := net.Listen("tcp", a.cfg.Listener.Port)
	if err != nil {
		a.logger.Fatal().Msgf("failed to setup listener: %v", err)
	}

	a.logger.Info().Msgf("starting server at %s", a.cfg.Listener.Port)

	defer func() {
		if err := a.GracefulShutdown(); err != nil {
			a.logger.Fatal().Msgf("failed to graceful shutdown: %v", err)
		}
	}()

	if err := a.srv.Serve(lis); err != nil {
		if errors.Is(err, grpc.ErrServerStopped) {
			a.logger.Info().Msg("server closed under request")
		} else {
			a.logger.Info().Msgf("server stopped: %v", err)
		}
	}
}

// GracefulShutdown gracefully shutdowns AppAuth
func (a *AppAuth) GracefulShutdown() error {
	a.logger.Info().Msg("Graceful shutdown")

	a.srv.GracefulStop()
	a.logger.Info().Msg("Auth service grpc shut down")

	return nil
}
