package repository

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/rs/zerolog/log"
)

const (
	getActorInfoByIDQuery = `
select
	p.full_name,
	p.en_full_name,
	p.photo,
	p.about,
	p.sex,
	p.growth,
	p.birthday,
	p.death,
	string_agg(DISTINCT CAST(c.career AS TEXT), ', ' ORDER BY CAST(c.career AS TEXT)) AS career,
	string_agg(DISTINCT g.name, ', ' ORDER BY g.name) AS genres
FROM career_person cp
JOIN person p ON cp.person_id = p.id
JOIN career c ON cp.career_id = c.id
JOIN person_genre pg ON pg.person_id = p.id
JOIN genre g ON pg.genre_id = g.id
where p.id = $1
GROUP BY p.id, p.full_name, p.en_full_name, p.photo, p.about, p.sex, p.growth, p.birthday, p.death;
`

	getTotalActorMoviesQuery = `
SELECT
	m.id as id,
	m.name as title,
	m.poster as preview_url,
	m.duration as duration,
	COUNT(m.id) OVER () AS total_actor_movies
FROM movie_staff ms
JOIN person p ON p.id = ms.staff_id
JOIN movie m ON m.id = ms.movie_id
WHERE p.id = $1;
`
)

// StaffPersonRepository collect and process data of staff person
type StaffPersonPostgresRepository struct {
	pgdb *sql.DB
}

// NewStaffPersonRepository returns new instance of StaffPersonRepository
func NewStaffPersonPostgresRepository(staffRepo *sql.DB) *StaffPersonPostgresRepository {
	return &StaffPersonPostgresRepository{pgdb: staffRepo}
}

// GetPersonFromRepoByID obtains person info from repo by id
func (r *StaffPersonPostgresRepository) GetPersonFromRepoByID(ctx context.Context, personID int) (*mocks.PersonJSON, error) {
	logger := log.Ctx(ctx)
	logger.Info().Msgf("Get Person by id %d from postgres repo", personID)

	var resPerson mocks.PersonJSON
	var personFullName sql.NullString
	var personEnFullName sql.NullString
	var personPhoto sql.NullString
	var personAbout sql.NullString
	var personSex sql.NullString
	var personGrowth sql.NullString
	var personBirthday sql.NullString
	var personDeath sql.NullString
	var personCareer sql.NullString
	var personGenres sql.NullString
	var personTotalFilms sql.NullString

	personMovieCollection := mocks.Collection{}

	execRowPerson, err := r.pgdb.Query(getActorInfoByIDQuery, personID)
	if err != nil {
		errMsg := errors.Wrap(err, "error in query statement in GetPersonFromRepoByID")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}
	defer func() {
		if closeErr := execRowPerson.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed_to_close_statement")
			return
		}
	}()

	for execRowPerson.Next() {
		// null the repeating values
		personCareer = sql.NullString{}
		personGenres = sql.NullString{}

		if err := execRowPerson.Scan(
			&personFullName,
			&personEnFullName,
			&personPhoto,
			&personAbout,
			&personSex,
			&personGrowth,
			&personBirthday,
			&personDeath,
			&personCareer,
			&personGenres,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in get person info in GetPersonFromRepoByID")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}
	}
	if execErr := execRowPerson.Err(); execErr != nil {
		errMsg := errors.Wrap(execErr, "error in query next in get movies in GetPersonFromRepoByID")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	// getting all actor movies
	execRow, err := r.pgdb.Query(getTotalActorMoviesQuery, personID)
	if err != nil {
		errMsg := errors.Wrap(err, "error in query statement in GetPersonFromRepoByID")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}
	defer func() {
		if closeErr := execRow.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed_to_close_statement")
			return
		}
	}()

	for execRow.Next() {
		var movieID int
		var movieTitle sql.NullString
		var moviePreviewURL sql.NullString
		var movieDuration sql.NullString

		if err := execRow.Scan(
			&movieID,
			&movieTitle,
			&moviePreviewURL,
			&movieDuration,
			&personTotalFilms,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in get movies in GetPersonFromRepoByID")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		movie := mocks.Movie{
			ID:         movieID,
			Title:      movieTitle.String,
			PreviewURL: moviePreviewURL.String,
			Duration:   movieDuration.String,
		}
		personMovieCollection[movie.ID] = movie
	}
	if execErr := execRow.Err(); execErr != nil {
		errMsg := errors.Wrap(execErr, "error in query next in get movies in GetPersonFromRepoByID")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	resPerson = mocks.PersonJSON{
		ID:              personID,
		FullName:        personFullName.String,
		EnFullName:      personEnFullName.String,
		Photo:           personPhoto.String,
		About:           personAbout.String,
		Sex:             personSex.String,
		Growth:          personGrowth.String,
		Birthday:        personBirthday.String,
		Death:           personDeath.String,
		Career:          personCareer.String,
		Genres:          personGenres.String,
		TotalFilms:      personTotalFilms.String,
		MovieCollection: personMovieCollection,
	}

	return &resPerson, nil
}
