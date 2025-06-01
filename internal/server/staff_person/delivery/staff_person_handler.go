package delivery

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mocks.go -package=delivery_mocks StaffPersonServiceInterface
type StaffPersonServiceInterface interface {
	GetPersonByID(ctx context.Context, personID int) (*mocks.PersonJSON, error)
}

// StaffPersonHandler handles requests to staff person: actor, director, etc
type StaffPersonHandler struct {
	staffPersonService StaffPersonServiceInterface
}

// NewStaffPersonHandler returns new instance of StaffPersonHandler
func NewStaffPersonHandler(staffPersonService StaffPersonServiceInterface) *StaffPersonHandler {
	return &StaffPersonHandler{
		staffPersonService: staffPersonService,
	}
}

// GetPerson handles GET request to obtain person info by id
func (h *StaffPersonHandler) GetPerson(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	personID, err := strconv.Atoi(mux.Vars(r)["person_id"])
	if err != nil {
		errMsg := errors.Wrap(err, "getPersonByID action: bad request")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errs.ErrBadPayload, errs.ErrBadPayload)
		return
	}

	logger.Info().Msgf("getting person by id: %d", personID)
	personJSON, err := h.staffPersonService.GetPersonByID(r.Context(), personID)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		if errors.Is(err, errs.ErrPersonNotFound) {
			jsonutil.SendError(r.Context(), w, http.StatusNotFound, errors.Wrap(err, errs.ErrNotFoundShort).Error(), err.Error())
			return
		}
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong, errs.ErrSomethingWentWrong)
		return
	}
	logger.Info().Msgf("successfully got person data by id: %d", personID)

	if err := jsonutil.SendJSON(r.Context(), w, personJSON); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}
}
