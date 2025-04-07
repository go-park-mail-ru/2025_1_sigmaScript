package repository

import (
	"context"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/rs/zerolog/log"
)

// StaffPersonRepository collect and process data of staff person
type StaffPersonRepository struct {
	staffRepo *mocks.Persons
}

// NewStaffPersonRepository returns new instance of StaffPersonRepository
func NewStaffPersonRepository(staffRepo *mocks.Persons) *StaffPersonRepository {
	return &StaffPersonRepository{staffRepo: staffRepo}
}

// GetPersonFromRepoByID obtains person info from repo by id
func (r *StaffPersonRepository) GetPersonFromRepoByID(ctx context.Context, personID int) (*mocks.PersonJSON, error) {
	logger := log.Ctx(ctx)

	for _, val := range *r.staffRepo {
		if val.ID == personID {
			return &val, nil
		}
	}

	logger.Err(errs.ErrPersonNotFound).Msg(errs.ErrPersonNotFound.Error())
	return nil, errs.ErrPersonNotFound
}
