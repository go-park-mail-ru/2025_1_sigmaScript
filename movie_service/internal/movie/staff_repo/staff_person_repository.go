package staff_repo

import (
	"context"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/rs/zerolog/log"
)

type StaffPersonRepository struct {
	db *mocks.Persons
}

func NewStaffPersonRepository(staffRepo *mocks.Persons) *StaffPersonRepository {
	return &StaffPersonRepository{db: staffRepo}
}

func (r *StaffPersonRepository) GetPersonFromRepoByID(ctx context.Context, personID int) (*mocks.PersonJSON, error) {
	logger := log.Ctx(ctx)

	for _, val := range *r.db {
		if val.ID == personID {
			return &val, nil
		}
	}

	logger.Err(errs.ErrPersonNotFound).Msg(errs.ErrPersonNotFound.Error())
	return nil, errs.ErrPersonNotFound
}
