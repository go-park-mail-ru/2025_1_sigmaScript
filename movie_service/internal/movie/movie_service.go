package movie

import (
	"context"
	"errors"
	"net"
	"os"
	"time"

	repoCollection "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/repository"
	serviceCollection "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/service"
	movie_pb "github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/api/movie_api_v1/proto"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/config"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/db"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/interceptors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/internal/movie/delivery"
	repoGenre "github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/internal/movie/genres_repo"
	serviceGenre "github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/internal/movie/genres_service"
	repoMovie "github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/internal/movie/movie_repo"
	service "github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/internal/movie/movie_service"
	repoSearch "github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/internal/movie/search_repo"
	serviceSearch "github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/internal/movie/search_service"
	repoStaff "github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/internal/movie/staff_repo"
	serviceStaff "github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/internal/movie/staff_service"
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

	// cfg := config.ConfigPgDB{
	// 	Listener: config.Listener{
	// 		Port: ":8083",
	// 	},
	// 	Databases: config.Databases{
	// 		Postgres: config.Postgres{
	// 			Host:            "127.0.0.1",
	// 			Port:            5433,
	// 			User:            "filmlk_user",
	// 			Password:        "filmlk_password",
	// 			Name:            "filmlk",
	// 			MaxOpenConns:    100,
	// 			MaxIdleConns:    30,
	// 			ConnMaxLifetime: 300,
	// 			ConnMaxIdleTime: 60,
	// 		},
	// 		LocalStorage: config.LocalAvatarsStorage{
	// 			UserAvatarsFullPath:   "/Users/propolisss/frontend/2025_1_sigmaScript/public/static/avatars/",
	// 			UserAvatarsStaticPath: "/static/avatars/",
	// 		},
	// 	},
	// }

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

	staffPersonRepo := repoStaff.NewStaffPersonPostgresRepository(database)
	staffPersonService := serviceStaff.NewStaffPersonService(staffPersonRepo)

	collectionRepo := repoCollection.NewCollectionPostgresRepository(database)
	collectionService := serviceCollection.NewCollectionService(collectionRepo)

	genreRepo := repoGenre.NewGenreRepository(database)
	genreService := serviceGenre.NewGenreService(genreRepo)

	searchRepo := repoSearch.NewSearchRepository(database)
	searchService := serviceSearch.NewSearchService(searchRepo)

	movieRepo := repoMovie.NewMoviePostgresRepository(database)
	movieSvc := service.NewMovieService(movieRepo)
	movie_pb.RegisterMovieRPCServer(srv, delivery.NewMovieServiceGRPCHandler(
		movieSvc,
		genreService,
		searchService,
		staffPersonService,
		collectionService,
	))

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
	a.logger.Info().Msg("Auth search_service grpc shut down")

	return nil
}
