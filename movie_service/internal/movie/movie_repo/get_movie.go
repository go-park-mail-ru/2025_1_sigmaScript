package movie_repo

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
	rating_kp,
	rating_imdb,
	m.duration,
	string_agg(DISTINCT g.name, ', ' ORDER BY g.name) AS genres,
	m.logo,
	m.backdrop
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
	rating_kp,
	rating_imdb,
	m.duration,
	m.logo,
	m.backdrop;
`

	getMovieStaffQuery = `
SELECT
	m.name as movie_title,
	p.id as id,
	p.full_name as full_name,
	p.photo as photo,
	ms.role as career
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
	resMovieRating := sql.NullFloat64{}
	resRatingKP := sql.NullFloat64{}
	resRatingIMDB := sql.NullFloat64{}
	resMovieDuration := sql.NullString{}
	resMovieGenres := sql.NullString{}
	resMovieLogo := sql.NullString{}
	resMovieBackdrop := sql.NullString{}

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
			&resRatingKP,
			&resRatingIMDB,
			&resMovieDuration,
			&resMovieGenres,
			&resMovieLogo,
			&resMovieBackdrop,
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
		var personCareer sql.NullString

		if err := execRowStaff.Scan(
			&movieName,
			&personID,
			&personFullName,
			&personPhoto,
			&personCareer,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in GetMovieFromRepoByID")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		// skip if not found
		if !personID.Valid {
			continue
		}

		person := mocks.PersonJSON{
			ID:       int(personID.Int64),
			FullName: personFullName.String,
			Photo:    personPhoto.String,
			Career:   personCareer.String,
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

		var reviewID sql.NullInt64
		var reviewScore sql.NullString
		var reviewText sql.NullString
		var reviewCreatedAt sql.NullString

		var userLogin sql.NullString
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

		// skip if bad user or review
		if !reviewID.Valid || !userLogin.Valid {
			continue
		}
		reviewUser := mocks.ReviewUserDataJSON{
			Login:  userLogin.String,
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
			ID:         int(reviewID.Int64),
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
		Rating:          resMovieRating.Float64,
		RatingKP:        resRatingKP.Float64,
		RatingIMDB:      resRatingIMDB.Float64,
		Duration:        resMovieDuration.String,
		Staff:           resStaff,
		Reviews:         resReviews,
		Genres:          resMovieGenres.String,
		Logo:            resMovieLogo.String,
		Backdrop:        resMovieBackdrop.String,
	}

	return &resMovie, nil
}
