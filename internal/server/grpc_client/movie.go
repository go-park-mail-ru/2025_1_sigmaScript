// movie_service/internal/movie/client/movie_client.go
package client

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	pb "github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/api/movie_api_v1/proto"
)

//go:generate mockgen -source=movie_client.go -destination=../movie/service/mocks/mock.go

// MovieClientInterface defines client methods for MovieRPC service.
type MovieClientInterface interface {
	GetMovieByID(ctx context.Context, movieID int) (*mocks.MovieJSON, error)
	GetAllReviewsOfMovieByID(ctx context.Context, movieID int) ([]mocks.ReviewJSON, error)
	CreateNewMovieReview(ctx context.Context, userID, movieID string, newReview mocks.NewReviewDataJSON) (mocks.NewReviewDataJSON, error)
	UpdateMovieReview(ctx context.Context, userID, movieID string, newReview mocks.NewReviewDataJSON) (mocks.NewReviewDataJSON, error)
	GetPersonByID(ctx context.Context, personID int) (*mocks.PersonJSON, error)
	GetGenreByID(ctx context.Context, genreID string) (*mocks.Genre, error)
	GetAllGenres(ctx context.Context) ([]mocks.Genre, error)
	SearchActorsAndMovies(ctx context.Context, query string) (*models.SearchResponseJSON, error)
	GetMainPageCollections(ctx context.Context) (mocks.Collections, error)
}

// MovieClient is a gRPC client for MovieRPC.
type MovieClient struct {
	svc pb.MovieRPCClient
}

// NewMovieClient returns a new MovieClient.
func NewMovieClient(svc pb.MovieRPCClient) *MovieClient {
	return &MovieClient{svc: svc}
}

func (cl *MovieClient) GetMovieByID(ctx context.Context, movieID int) (*mocks.MovieJSON, error) {
	resp, err := cl.svc.GetMovieByID(ctx, &pb.GetMovieByIDRequest{MovieId: int32(movieID)})
	if err != nil {
		return nil, err
	}
	pbM := resp.GetMovie()
	if pbM == nil {
		return nil, nil
	}
	return &mocks.MovieJSON{
		ID:              int(pbM.GetId()),
		Name:            pbM.GetName(),
		OriginalName:    pbM.GetOriginalName(),
		About:           pbM.GetAbout(),
		Poster:          pbM.GetPoster(),
		PromoURL:        pbM.GetPromoUrl(),
		ReleaseYear:     pbM.GetReleaseYear(),
		Country:         pbM.GetCountry(),
		Slogan:          pbM.GetSlogan(),
		Director:        pbM.GetDirector(),
		Budget:          pbM.GetBudget(),
		BoxOfficeUS:     pbM.GetBoxOfficeUs(),
		BoxOfficeGlobal: pbM.GetBoxOfficeGlobal(),
		BoxOfficeRussia: pbM.GetBoxOfficeRussia(),
		PremierRussia:   pbM.GetPremierRussia(),
		PremierGlobal:   pbM.GetPremierGlobal(),
		Rating:          pbM.GetRating(),
		Duration:        pbM.GetDuration(),
		Genres:          pbM.GetGenres(),
		Staff: func() []mocks.PersonJSON {
			out := make([]mocks.PersonJSON, 0, len(pbM.GetStaff()))
			for _, p := range pbM.GetStaff() {
				out = append(out, mocks.PersonJSON{
					ID:         int(p.GetId()),
					FullName:   p.GetFullName(),
					EnFullName: p.GetEnFullName(),
					Photo:      p.GetPhoto(),
					About:      p.GetAbout(),
					Sex:        p.GetSex(),
					Growth:     p.GetGrowth(),
					Birthday:   p.GetBirthday(),
					Death:      p.GetDeath(),
					Career:     p.GetCareer(),
					Genres:     p.GetGenres(),
					TotalFilms: p.GetTotalFilms(),
					MovieCollection: func() mocks.Collection {
						coll := make(mocks.Collection)
						for k, mv := range p.GetMovieCollection() {
							coll[int(k)] = mocks.Movie{
								ID:          int(mv.GetId()),
								Title:       mv.GetName(),
								PreviewURL:  mv.GetPoster(),
								Duration:    mv.GetDuration(),
								ReleaseDate: mv.GetPremierRussia(),
								Rating:      mv.GetRating(),
							}
						}
						return coll
					}(),
				})
			}
			return out
		}(),
		Reviews: func() []mocks.ReviewJSON {
			out := make([]mocks.ReviewJSON, 0, len(pbM.GetReviews()))
			for _, r := range pbM.GetReviews() {
				out = append(out, mocks.ReviewJSON{
					ID:         int(r.GetId()),
					ReviewText: r.GetReviewText(),
					Score:      r.GetScore(),
					CreatedAt:  r.GetCreatedAt(),
					User: mocks.ReviewUserDataJSON{
						Login:  r.GetUser().GetLogin(),
						Avatar: r.GetUser().GetAvatar(),
					},
				})
			}
			return out
		}(),
		RatingKP:   pbM.GetRatingKp(),
		RatingIMDB: pbM.GetRatingImdb(),
		Watchability: func() []mocks.WatchProviderJSON {
			out := make([]mocks.WatchProviderJSON, 0, len(pbM.GetWatchability()))
			for _, w := range pbM.GetWatchability() {
				out = append(out, mocks.WatchProviderJSON{
					Name: w.GetName(),
					Url:  w.GetUrl(),
				})
			}
			return out
		}(),
		Logo:     pbM.GetLogo(),
		Backdrop: pbM.GetBackdrop(),
	}, nil
}

func (cl *MovieClient) GetAllReviewsOfMovieByID(ctx context.Context, movieID int) (*[]mocks.ReviewJSON, error) {
	resp, err := cl.svc.GetAllReviewsOfMovieByID(ctx, &pb.GetAllReviewsOfMovieByIDRequest{MovieId: int32(movieID)})
	if err != nil {
		return nil, err
	}
	var out []mocks.ReviewJSON
	for _, r := range resp.GetReviews() {
		out = append(out, mocks.ReviewJSON{
			ID:         int(r.GetId()),
			ReviewText: r.GetReviewText(),
			Score:      r.GetScore(),
			CreatedAt:  r.GetCreatedAt(),
			User: mocks.ReviewUserDataJSON{
				Login: r.GetUser().GetLogin(),
			},
		})
	}
	return &out, nil
}

func (cl *MovieClient) CreateNewMovieReview(ctx context.Context, userID, movieID string, newReview mocks.NewReviewDataJSON) (*mocks.NewReviewDataJSON, error) {
	resp, err := cl.svc.CreateNewMovieReview(ctx, &pb.CreateNewMovieReviewRequest{
		UserId:    userID,
		MovieId:   movieID,
		NewReview: &pb.NewReviewDataJSON{ReviewText: newReview.ReviewText, Score: newReview.Score},
	})
	if err != nil {
		return &mocks.NewReviewDataJSON{}, err
	}
	nb := resp.GetReview()
	return &mocks.NewReviewDataJSON{ReviewText: nb.GetReviewText(), Score: nb.GetScore()}, nil
}

func (cl *MovieClient) UpdateMovieReview(ctx context.Context, userID, movieID string, newReview mocks.NewReviewDataJSON) (*mocks.NewReviewDataJSON, error) {
	resp, err := cl.svc.UpdateMovieReview(ctx, &pb.UpdateMovieReviewRequest{
		UserId:    userID,
		MovieId:   movieID,
		NewReview: &pb.NewReviewDataJSON{ReviewText: newReview.ReviewText, Score: newReview.Score},
	})
	if err != nil {
		return &mocks.NewReviewDataJSON{}, err
	}
	nb := resp.GetReview()
	return &mocks.NewReviewDataJSON{ReviewText: nb.GetReviewText(), Score: nb.GetScore()}, nil
}

func (cl *MovieClient) GetPersonByID(ctx context.Context, personID int) (*mocks.PersonJSON, error) {
	resp, err := cl.svc.GetPersonByID(ctx, &pb.GetPersonByIDRequest{PersonId: int32(personID)})
	if err != nil {
		return nil, err
	}
	pbP := resp.GetPerson()
	if pbP == nil {
		return nil, nil
	}
	coll := make(mocks.Collection)
	for k, mv := range pbP.GetMovieCollection() {
		coll[int(k)] = mocks.Movie{
			ID:          int(mv.GetId()),
			Title:       mv.GetName(),
			PreviewURL:  mv.GetPoster(),
			Duration:    mv.GetDuration(),
			ReleaseDate: mv.GetReleaseYear(),
			Rating:      mv.GetRating(),
		}
	}
	return &mocks.PersonJSON{
		ID:              int(pbP.GetId()),
		FullName:        pbP.GetFullName(),
		EnFullName:      pbP.GetEnFullName(),
		Photo:           pbP.GetPhoto(),
		About:           pbP.GetAbout(),
		Sex:             pbP.GetSex(),
		Growth:          pbP.GetGrowth(),
		Birthday:        pbP.GetBirthday(),
		Death:           pbP.GetDeath(),
		Career:          pbP.GetCareer(),
		Genres:          pbP.GetGenres(),
		TotalFilms:      pbP.GetTotalFilms(),
		MovieCollection: coll,
	}, nil
}

func (cl *MovieClient) GetGenreByID(ctx context.Context, genreID string) (*mocks.Genre, error) {
	resp, err := cl.svc.GetGenreByID(ctx, &pb.GetGenreByIDRequest{GenreId: genreID})
	if err != nil {
		return nil, err
	}
	pbG := resp.GetGenre()
	if pbG == nil {
		return nil, nil
	}
	movies := make([]mocks.Movie, 0, len(pbG.GetMovies()))
	for _, mv := range pbG.GetMovies() {
		movies = append(movies, mocks.Movie{
			ID:          int(mv.GetId()),
			Title:       mv.GetName(),
			PreviewURL:  mv.GetPoster(),
			Duration:    mv.GetDuration(),
			ReleaseDate: mv.GetReleaseYear(),
			Rating:      mv.GetRating(),
		})
	}
	return &mocks.Genre{ID: pbG.GetId(), Name: pbG.GetName(), Movies: movies}, nil
}

func (cl *MovieClient) GetAllGenres(ctx context.Context) (*[]mocks.Genre, error) {
	resp, err := cl.svc.GetAllGenres(ctx, &pb.GetAllGenresRequest{})
	if err != nil {
		return nil, err
	}
	var out []mocks.Genre
	for _, g := range resp.GetGenres() {
		movies := make([]mocks.Movie, 0, len(g.GetMovies()))
		for _, mv := range g.GetMovies() {
			movies = append(movies, mocks.Movie{
				ID:          int(mv.GetId()),
				Title:       mv.GetName(),
				PreviewURL:  mv.GetPoster(),
				Duration:    mv.GetDuration(),
				ReleaseDate: mv.GetReleaseYear(),
				Rating:      mv.GetRating(),
			})
		}
		out = append(out, mocks.Genre{ID: g.GetId(), Name: g.GetName(), Movies: movies})
	}
	return &out, nil
}

func (cl *MovieClient) SearchActorsAndMovies(ctx context.Context, query string) (*models.SearchResponseJSON, error) {
	resp, err := cl.svc.SearchActorsAndMovies(ctx, &pb.SearchActorsAndMoviesRequest{Query: query})
	if err != nil {
		return nil, err
	}
	result := &models.SearchResponseJSON{}
	for _, m := range resp.GetMovies() {
		result.MovieCollection = append(result.MovieCollection, mocks.Movie{
			ID:          int(m.GetId()),
			Title:       m.GetName(),
			PreviewURL:  m.GetPoster(),
			Duration:    m.GetDuration(),
			ReleaseDate: m.GetReleaseYear(),
			Rating:      m.GetRating(),
		})
	}
	for _, a := range resp.GetActors() {
		result.Actors = append(result.Actors, mocks.PersonJSON{
			ID:         int(a.GetId()),
			FullName:   a.GetFullName(),
			EnFullName: a.GetEnFullName(),
			Photo:      a.GetPhoto(),
			About:      a.GetAbout(),
		})
	}
	return result, nil
}

func (cl *MovieClient) GetMainPageCollections(ctx context.Context) (mocks.Collections, error) {
	resp, err := cl.svc.GetMainPageCollections(ctx, &pb.GetMainPageCollectionsRequest{})
	if err != nil {
		return nil, err
	}
	cols := make(mocks.Collections)
	for _, c := range resp.GetCollections() {
		coll := make(mocks.Collection)
		for _, m := range c.GetMovies() {
			coll[int(m.GetId())] = mocks.Movie{
				ID:          int(m.GetId()),
				Title:       m.GetName(),
				PreviewURL:  m.GetPoster(),
				Duration:    m.GetDuration(),
				ReleaseDate: m.GetReleaseYear(),
				Rating:      m.GetRating(),
			}
		}
		cols[c.GetName()] = coll
	}
	return cols, nil
}
