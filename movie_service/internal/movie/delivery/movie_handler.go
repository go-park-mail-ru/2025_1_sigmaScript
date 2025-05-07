// movie_service/internal/movie/delivery/movie_service_handler.go
package delivery

import (
	"context"

	pb "github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/api/movie_api_v1/proto"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/internal/movie/delivery/adapter"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/internal/movie/delivery/interfaces"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type MovieServiceGRPCHandler struct {
	pb.UnimplementedMovieRPCServer
	movieService      interfaces.MovieServiceInterface
	genreService      interfaces.GenreServiceInterface
	searchService     interfaces.SearchServiceInterface
	staffService      interfaces.StaffPersonServiceInterface
	collectionService interfaces.CollectionServiceInterface
}

func NewMovieServiceGRPCHandler(movieService interfaces.MovieServiceInterface,
	genreService interfaces.GenreServiceInterface,
	searchService interfaces.SearchServiceInterface,
	staffService interfaces.StaffPersonServiceInterface,
	collectionService interfaces.CollectionServiceInterface,
) *MovieServiceGRPCHandler {
	return &MovieServiceGRPCHandler{
		movieService:      movieService,
		genreService:      genreService,
		searchService:     searchService,
		staffService:      staffService,
		collectionService: collectionService,
	}
}

func (h *MovieServiceGRPCHandler) GetMovieByID(ctx context.Context, req *pb.GetMovieByIDRequest) (*pb.GetMovieByIDResponse, error) {
	logger := log.Ctx(ctx)

	if req == nil {
		logger.Error().Msg("GetMovieByID: request is nil")
		return nil, errors.New("bad request: request is nil")
	}
	logger.Info().Int32("movie_id", req.MovieId).Msg("GetMovieByID request received")

	movieData, err := h.movieService.GetMovieByID(ctx, int(req.MovieId))
	if err != nil {
		logger.Error().Err(err).Msg("GetMovieByID: movie_service error")
		return nil, err
	}

	logger.Info().Interface("movie_data_name", movieData.Name).Msg("GetMovieByID: movie found")
	return adapter.ToDescGetMovieByIDResponse(movieData), nil
}

func (h *MovieServiceGRPCHandler) GetAllReviewsOfMovieByID(ctx context.Context, req *pb.GetAllReviewsOfMovieByIDRequest) (*pb.GetAllReviewsOfMovieByIDResponse, error) {
	logger := log.Ctx(ctx)

	if req == nil {
		logger.Error().Msg("GetAllReviewsOfMovieByID: request is nil")
		return nil, errors.New("bad request: request is nil")
	}
	logger.Info().Int32("movie_id", req.MovieId).Msg("GetAllReviewsOfMovieByID request received")

	reviewsData, err := h.movieService.GetAllReviewsOfMovieByID(ctx, int(req.MovieId))
	if err != nil {
		logger.Error().Err(err).Msg("GetAllReviewsOfMovieByID: movie_service error")
		return nil, err
	}

	if reviewsData != nil {
		logger.Info().Int("reviews_count", len(*reviewsData)).Msg("GetAllReviewsOfMovieByID: reviews fetched")
	} else {
		logger.Info().Msg("GetAllReviewsOfMovieByID: no reviews found or movie_service returned nil slice pointer")
	}
	return adapter.ToDescGetAllReviewsOfMovieByIDResponse(reviewsData), nil
}

func (h *MovieServiceGRPCHandler) CreateNewMovieReview(ctx context.Context, req *pb.CreateNewMovieReviewRequest) (*pb.CreateNewMovieReviewResponse, error) {
	logger := log.Ctx(ctx)

	if req == nil || req.NewReview == nil {
		logger.Error().Msg("CreateNewMovieReview: request or new_review field is nil")
		return nil, errors.New("bad request: request or new_review is nil")
	}
	logger.Info().Str("user_id", req.UserId).Str("movie_id", req.MovieId).Msg("CreateNewMovieReview request received")

	srvNewReviewData := adapter.ToSrvNewReviewDataFromDesc(req)

	createdReview, err := h.movieService.CreateNewMovieReview(ctx, req.UserId, req.MovieId, srvNewReviewData)
	if err != nil {
		logger.Error().Err(err).Msg("CreateNewMovieReview: movie_service error")
		return nil, err
	}

	logger.Info().Msg("CreateNewMovieReview: review created successfully")
	return adapter.ToDescCreateNewMovieReviewResponse(createdReview), nil
}

func (h *MovieServiceGRPCHandler) UpdateMovieReview(ctx context.Context, req *pb.UpdateMovieReviewRequest) (*pb.UpdateMovieReviewResponse, error) {
	logger := log.Ctx(ctx)

	if req == nil || req.NewReview == nil {
		logger.Error().Msg("UpdateMovieReview: request or new_review field is nil")
		return nil, errors.New("bad request: request or new_review is nil")
	}
	logger.Info().Str("user_id", req.UserId).Str("movie_id", req.MovieId).Msg("UpdateMovieReview request received")

	srvNewReviewData := adapter.ToSrvUpdateReviewDataFromDesc(req)

	updatedReview, err := h.movieService.UpdateMovieReview(ctx, req.UserId, req.MovieId, srvNewReviewData)
	if err != nil {
		logger.Error().Err(err).Msg("UpdateMovieReview: movie_service error")
		return nil, err
	}

	logger.Info().Msg("UpdateMovieReview: review updated successfully")
	return adapter.ToDescUpdateMovieReviewResponse(updatedReview), nil
}

func (h *MovieServiceGRPCHandler) GetPersonByID(ctx context.Context, req *pb.GetPersonByIDRequest) (*pb.GetPersonByIDResponse, error) {
	logger := log.Ctx(ctx)
	if req == nil {
		return nil, errors.New("GetPersonByID: empty request")
	}
	pid := adapter.ToSrvPersonID(req)
	person, err := h.staffService.GetPersonByID(ctx, pid)
	if err != nil {
		logger.Error().Err(err).Msg("GetPersonByID failed")
		return nil, err
	}
	return adapter.ToDescGetPersonByIDResponse(person), nil
}

func (h *MovieServiceGRPCHandler) GetGenreByID(ctx context.Context, req *pb.GetGenreByIDRequest) (*pb.GetGenreByIDResponse, error) {
	logger := log.Ctx(ctx)
	if req == nil {
		return nil, errors.New("GetGenreByID: empty request")
	}
	gid := adapter.ToSrvGenreID(req)
	genre, err := h.genreService.GetGenreByID(ctx, gid)
	if err != nil {
		logger.Error().Err(err).Msg("GetGenreByID failed")
		return nil, err
	}
	return adapter.ToDescGetGenreByIDResponse(genre), nil
}

func (h *MovieServiceGRPCHandler) GetAllGenres(ctx context.Context, _ *pb.GetAllGenresRequest) (*pb.GetAllGenresResponse, error) {
	logger := log.Ctx(ctx)
	genres, err := h.genreService.GetAllGenres(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("GetAllGenres failed")
		return nil, err
	}
	return adapter.ToDescGetAllGenresResponse(genres), nil
}

func (h *MovieServiceGRPCHandler) SearchActorsAndMovies(ctx context.Context, req *pb.SearchActorsAndMoviesRequest) (*pb.SearchActorsAndMoviesResponse, error) {
	logger := log.Ctx(ctx)
	if req == nil {
		return nil, errors.New("SearchActorsAndMovies: empty request")
	}
	query := adapter.ToSrvSearchQuery(req)
	sr, err := h.searchService.SearchActorsAndMovies(ctx, query)
	if err != nil {
		logger.Error().Err(err).Msg("SearchActorsAndMovies failed")
		return nil, err
	}
	return adapter.ToDescSearchActorsAndMoviesResponse(sr), nil
}

func (h *MovieServiceGRPCHandler) GetMainPageCollections(ctx context.Context, _ *pb.GetMainPageCollectionsRequest) (*pb.GetMainPageCollectionsResponse, error) {
	logger := log.Ctx(ctx)
	logger.Info().Msg("GetMainPageCollections request received")

	cols, err := h.collectionService.GetMainPageCollections(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("service.GetMainPageCollections failed")
		return nil, errors.Wrap(err, "could not fetch main page collections")
	}

	return adapter.ToDescGetMainPageCollectionsResponse(cols), nil
}
