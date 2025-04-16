package service

import (
	"context"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/staff_person/service/service_mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestStaffPersonService_GetPersonByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := service_mocks.NewMockStaffPersonRepositoryInterface(ctrl)
	service := NewStaffPersonService(mockRepo)

	keanuReeves := mocks.ExistingActors[11]

	tests := []struct {
		name        string
		personID    int
		mockSetup   func()
		expected    *mocks.PersonJSON
		expectedErr error
	}{
		{
			name:     "Success - Get Keanu Reeves",
			personID: 11,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetPersonFromRepoByID(gomock.Any(), 11).
					Return(&keanuReeves, nil)
			},
			expected:    &keanuReeves,
			expectedErr: nil,
		},
		{
			name:     "Fail - Person not found",
			personID: 999,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetPersonFromRepoByID(gomock.Any(), 999).
					Return(nil, errs.ErrPersonNotFound)
			},
			expected:    nil,
			expectedErr: errs.ErrPersonNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			person, err := service.GetPersonByID(context.Background(), tt.personID)

			assert.Equal(t, tt.expected, person)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
