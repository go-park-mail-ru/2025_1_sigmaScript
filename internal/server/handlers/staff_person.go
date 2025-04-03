package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
)

type StaffPersonInterface interface {
	GetPersonByID(w http.ResponseWriter, r *http.Request)
}

// StaffPersonHandler handles requests to staff person: actor, director, etc
type StaffPersonHandler struct {
	staffRepo *mocks.Persons
}

// NewStaffPersonHandler returns new instance of StaffPersonHandler
func NewStaffPersonHandler(staffRepo *mocks.Persons) *StaffPersonHandler {
	return &StaffPersonHandler{staffRepo: staffRepo}
}

// GetPersonByID handles GET request to obtain person info by id: actor, director, etc
func (sph *StaffPersonHandler) GetPersonByID(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	personID, err := strconv.Atoi(mux.Vars(r)["person_id"])
	if err != nil {
		errMsg := errors.Wrapf(err, "getPersonByID action: bad request: %w", err)
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		jsonutil.SendError(w, http.StatusBadRequest, errs.ErrBadPayload, errs.ErrBadPayload)
		return
	}
	logger.Info().Msgf("getting person by id: %d", personID)

	personJSON, err := sph.getPersonFromRepoByID(personID)
	if err != nil {
		if errors.Is(err, errs.ErrPersonNotFound) {
			logger.Error().Err(err).Msg(err.Error())
			jsonutil.SendError(w, http.StatusNotFound, err.Error(), err.Error())
			return
		}
		logger.Error().Err(err).Msg(err.Error())
		jsonutil.SendError(w, http.StatusInternalServerError, errs.ErrSomethingWentWrong, errs.ErrSomethingWentWrong)
		return
	}

	logger.Info().Msgf("successfully got person data by id: %d", personID)
	if err := jsonutil.SendJSON(w, personJSON); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}

}

func (sph *StaffPersonHandler) getPersonFromRepoByID(personID int) (*mocks.PersonJSON, error) {
	for _, val := range *sph.staffRepo {
		if val.ID == personID {
			return &val, nil
		}
	}
	return nil, errs.ErrPersonNotFound
}
