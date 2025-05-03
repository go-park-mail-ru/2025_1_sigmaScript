package repository

import (
	"context"
	"database/sql"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/rs/zerolog/log"
)

const (
	getMovieInfoByIDQuery = `
SELECT
	m.id,
	m.name,
	m.original_name,
	m.about,
	m.poster,
	m.release_year,
	m.country,
	m.slogan,
	m.budget,
	m.box_office_us,
	m.box_office_global,
	m.box_office_russia,
	m.premier_russia,
	m.premier_global,
	m.rating,
	m.duration,
	string_agg(DISTINCT g.name, ', ' ORDER BY g.name) AS genres
FROM
	movie m
LEFT JOIN
	movie_genre mg ON m.id = mg.movie_id
LEFT JOIN
	genre g ON mg.genre_id = g.id
WHERE
	m.id = $1
GROUP BY
	m.id,
	m.name,
	m.original_name,
	m.about,
	m.poster,
	m.release_year,
	m.country,
	m.slogan,
	m.budget,
	m.box_office_us,
	m.box_office_global,
	m.box_office_russia,
	m.premier_russia,
	m.premier_global,
	m.rating,
	m.duration;
`

	getMovieStaffQuery = `
SELECT
	m.name as movie_title,
	p.id as id,
	p.full_name,
	p.photo
FROM movie m
LEFT JOIN movie_staff ms ON m.id = ms.movie_id
LEFT JOIN person p ON p.id = ms.staff_id
WHERE m.id = $1;
`

	getMovieReviewsQuery = `
SELECT
	m.name as movie_title,
	r.id,
	r.score,
	r.review_text,
	r.updated_at as created_at,
	u.login,
	u.avatar
FROM movie m
LEFT JOIN review r ON m.id = r.movie_id
LEFT JOIN "user" u ON r.user_id = u.id
WHERE m.id = $1;	
`
)

type MovieJSON struct {
	ID              int                `json:"id"`
	Name            string             `json:"name"`
	OriginalName    string             `json:"original_name,omitempty"`
	About           string             `json:"about,omitempty"`
	Poster          string             `json:"poster,omitempty"`
	ReleaseYear     int                `json:"release_year,omitempty"`
	Country         string             `json:"country,omitempty"`
	Slogan          string             `json:"slogan,omitempty"`
	Director        string             `json:"director,omitempty"`
	Budget          int64              `json:"budget,omitempty"`
	BoxOfficeUS     int64              `json:"box_office_us,omitempty"`
	BoxOfficeGlobal int64              `json:"box_office_global,omitempty"`
	BoxOfficeRussia int64              `json:"box_office_russia,omitempty"`
	PremierRussia   string             `json:"premier_russia,omitempty"`
	PremierGlobal   string             `json:"premier_global,omitempty"`
	Rating          float64            `json:"rating,omitempty"`
	Duration        string             `json:"duration,omitempty"`
	Genres          string             `json:"genres,omitempty"`
	Staff           []mocks.PersonJSON `json:"staff,omitempty"`
	Reviews         []mocks.ReviewJSON `json:"reviews,omitempty"`
}

const (
	DEFAULT_MOVIE_SCORE = 5
)

type MoviePostgresRepository struct {
	pgdb *sql.DB
}

func NewMoviePostgresRepository(movieDB *sql.DB) *MoviePostgresRepository {
	return &MoviePostgresRepository{pgdb: movieDB}
}

func (r *MoviePostgresRepository) GetMovieFromRepoByID(ctx context.Context, movieID int) (*mocks.MovieJSON, error) {
	logger := log.Ctx(ctx)
	logger.Info().Msgf("Get Movie data from postgres repo by id %d ", movieID)

	resMovie := mocks.MovieJSON{}
	resMovieID := sql.NullString{}
	resMovieName := sql.NullString{}
	resMovieOriginalName := sql.NullString{}
	resMovieAbout := sql.NullString{}
	resMoviePoster := sql.NullString{}
	resMovieReleaseYear := sql.NullString{}
	resMovieCountry := sql.NullString{}
	resMovieSlogan := sql.NullString{}
	resMovieBudget := sql.NullInt64{}
	resMovieBoxOfficeUS := sql.NullInt64{}
	resMovieBoxOfficeGlobal := sql.NullInt64{}
	resMovieBoxOfficeRussia := sql.NullInt64{}
	resMoviePremierRussia := sql.NullString{}
	resMoviePremierGlobal := sql.NullString{}
	resMovieRating := sql.NullString{}
	resMovieDuration := sql.NullString{}
	resMovieGenres := sql.NullString{}

	// get movie info
	execRow, err := r.pgdb.Query(getMovieInfoByIDQuery, movieID)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "error in query statement in GetMovieFromRepoByID").Error())
		return nil, errors.Wrap(err, "error in prepare query statement in GetMovieFromRepoByID")
	}
	defer func() {
		if closeErr := execRow.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed_to_close_statement")
			return
		}
	}()

	for execRow.Next() {
		if err := execRow.Scan(
			&resMovieID,
			&resMovieName,
			&resMovieOriginalName,
			&resMovieAbout,
			&resMoviePoster,
			&resMovieReleaseYear,
			&resMovieCountry,
			&resMovieSlogan,
			&resMovieBudget,
			&resMovieBoxOfficeUS,
			&resMovieBoxOfficeGlobal,
			&resMovieBoxOfficeRussia,
			&resMoviePremierRussia,
			&resMoviePremierGlobal,
			&resMovieRating,
			&resMovieDuration,
			&resMovieGenres,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in GetMovieFromRepoByID")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}
	}
	if execErr := execRow.Err(); execErr != nil {
		errMsg := errors.Wrap(execErr, "error in query next in GetMovieFromRepoByID")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	resMovieRatingFloat, err := strconv.ParseFloat(resMovieRating.String, 64)
	if err != nil {
		errMsg := errors.Wrap(err, "error in query statement in GetMovieFromRepoByID")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	// get staff
	resStaff := []mocks.PersonJSON{}
	execRowStaff, err := r.pgdb.Query(getMovieStaffQuery, movieID)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "error in query statement in GetMovieFromRepoByID").Error())
		return nil, errors.Wrap(err, "error in prepare query statement in GetMovieFromRepoByID")
	}
	defer func() {
		if closeErr := execRowStaff.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed_to_close_statement")
			return
		}
	}()

	for execRowStaff.Next() {
		var movieName string

		var personID sql.NullInt64
		var personFullName sql.NullString
		var personPhoto sql.NullString

		if err := execRowStaff.Scan(
			&movieName,
			&personID,
			&personFullName,
			&personPhoto,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in GetMovieFromRepoByID")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		person := mocks.PersonJSON{
			ID:       int(personID.Int64),
			FullName: personFullName.String,
			Photo:    personPhoto.String,
		}

		resStaff = append(resStaff, person)
	}
	if execErr := execRowStaff.Err(); execErr != nil {
		errMsg := errors.Wrap(execErr, "error in query next in GetMovieFromRepoByID")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	// get reviews
	resReviews := []mocks.ReviewJSON{}
	execRowReviews, err := r.pgdb.Query(getMovieReviewsQuery, movieID)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "error in query statement in GetMovieFromRepoByID").Error())
		return nil, errors.Wrap(err, "error in prepare query statement in GetMovieFromRepoByID")
	}
	defer func() {
		if closeErr := execRowReviews.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed_to_close_statement")
			return
		}
	}()

	for execRowReviews.Next() {
		var movieName string

		var reviewID int
		var reviewScore sql.NullString
		var reviewText sql.NullString
		var reviewCreatedAt sql.NullString

		var userLogin string
		var userAvatar sql.NullString

		if err := execRowReviews.Scan(
			&movieName,
			&reviewID,
			&reviewScore,
			&reviewText,
			&reviewCreatedAt,
			&userLogin,
			&userAvatar,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in GetMovieFromRepoByID")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		reviewUser := mocks.ReviewUserDataJSON{
			Login:  userLogin,
			Avatar: userAvatar.String,
		}

		// convert score from string to float
		reviewScoreFloat, err := strconv.ParseFloat(reviewScore.String, 64)
		if err != nil {
			errMsg := errors.Wrap(err, "error in query statement in GetMovieFromRepoByID")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		review := mocks.ReviewJSON{
			ID:         reviewID,
			Score:      reviewScoreFloat,
			ReviewText: reviewText.String,
			CreatedAt:  reviewCreatedAt.String,
			User:       reviewUser,
		}

		resReviews = append(resReviews, review)
	}
	if execErr := execRowReviews.Err(); execErr != nil {
		errMsg := errors.Wrap(execErr, "error in query next in GetMovieFromRepoByID")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	// result
	resMovie = mocks.MovieJSON{
		ID:              movieID,
		Name:            resMovieName.String,
		OriginalName:    resMovieOriginalName.String,
		About:           resMovieAbout.String,
		Poster:          resMoviePoster.String,
		ReleaseYear:     resMovieReleaseYear.String,
		Country:         resMovieCountry.String,
		Slogan:          resMovieSlogan.String,
		Budget:          resMovieBudget.Int64,
		BoxOfficeUS:     resMovieBoxOfficeUS.Int64,
		BoxOfficeGlobal: resMovieBoxOfficeGlobal.Int64,
		BoxOfficeRussia: resMovieBoxOfficeRussia.Int64,
		PremierRussia:   resMoviePremierRussia.String,
		PremierGlobal:   resMoviePremierGlobal.String,
		Rating:          resMovieRatingFloat,
		Duration:        resMovieDuration.String,
		Staff:           resStaff,
		Reviews:         resReviews,
	}

	return &resMovie, nil
}
