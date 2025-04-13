package repository

import (
	"context"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/rs/zerolog/log"
)

type MovieRepository struct {
	db *mocks.Movies
}

func NewMovieRepository(movieDB *mocks.Movies) *MovieRepository {
	return &MovieRepository{db: movieDB}
}

func (r *MovieRepository) GetMovieFromRepoByID(ctx context.Context, movieID int) (*mocks.MovieJSON, error) {
	logger := log.Ctx(ctx)

	movie, exists := (*r.db)[movieID]
	if !exists {
		logger.Err(errs.ErrMovieNotFound).Msg(errs.ErrMovieNotFound.Error())
		return nil, errs.ErrMovieNotFound
	}

	return &movie, nil
}
