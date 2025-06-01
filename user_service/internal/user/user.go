package user

import (
	"context"
	"net"
	"os"
	"time"

	user "github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/api/user_api_v1/proto"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/config"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/db"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/internal/interceptors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/internal/user/delivery"
	repoUser "github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/internal/user/repository"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/internal/user/service"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	logger *zerolog.Logger
	srv    *grpc.Server
	cfg    *config.Config
}

func New(isTest bool) (*App, error) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.ChainUnaryInterceptors(
			interceptors.LoggerInterceptor,
			interceptors.AccessLogInterceptor,
		)),
	)

	reflection.Register(srv)

	cfg, err := config.New()
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create config")
		return nil, err
	}

	ctx := config.WrapDatabaseContext(context.Background(), &cfg.Database)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	database, err := db.SetupDatabase(ctxWithTimeout, cancel)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot setup database")
		return nil, err
	}

	userRepo := repoUser.NewUserRepository(database)

	userServ := service.NewUserService(userRepo)
	user.RegisterUserServiceServer(srv, delivery.NewUserServiceGRPCHandler(userServ))

	return &App{
		srv:    srv,
		logger: &logger,
		cfg:    cfg,
	}, nil
}

func (a *App) Run() {
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

func (a *App) GracefulShutdown() error {
	a.logger.Info().Msg("Graceful shutdown")

	a.srv.GracefulStop()
	a.logger.Info().Msg("Auth service grpc shut down")

	return nil
}
