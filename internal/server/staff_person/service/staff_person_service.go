package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mocks.go -package=service_mocks StaffPersonRepositoryInterface
type StaffPersonRepositoryInterface interface {
	GetPersonFromRepoByID(ctx context.Context, personID int) (*mocks.PersonJSON, error)
}

// StaffPersonService collect and process data of staff person
type StaffPersonService struct {
	staffPersonRepo StaffPersonRepositoryInterface
}

// NewStaffPersonService returns new instance of StaffPersonService
func NewStaffPersonService(staffPersonRepo StaffPersonRepositoryInterface) *StaffPersonService {
	return &StaffPersonService{
		staffPersonRepo: staffPersonRepo,
	}
}

// GetPersonByID obtains and formats person info gotten by id
func (s *StaffPersonService) GetPersonByID(ctx context.Context, personID int) (*mocks.PersonJSON, error) {
	logger := log.Ctx(ctx)

	personJSON, err := s.staffPersonRepo.GetPersonFromRepoByID(ctx, personID)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return nil, err
	}

	return personJSON, nil
}
