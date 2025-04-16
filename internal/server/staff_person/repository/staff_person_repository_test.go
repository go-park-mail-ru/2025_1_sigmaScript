package repository

import (
	"context"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestStaffPersonRepository_GetPersonFromRepoByID(t *testing.T) {
	db := mocks.ExistingActors
	repo := NewStaffPersonRepository(&db)

	tests := []struct {
		name        string
		personID    int
		expected    *mocks.PersonJSON
		expectedErr error
	}{
		{
			name:     "OK. Get Keanu Reeves",
			personID: 11,
			expected: func() *mocks.PersonJSON {
				person := mocks.ExistingActors[11]
				return &person
			}(),
			expectedErr: nil,
		},
		{
			name:        "Fail. Person not found",
			personID:    999,
			expected:    nil,
			expectedErr: errs.ErrPersonNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			person, err := repo.GetPersonFromRepoByID(context.Background(), tt.personID)

			assert.Equal(t, tt.expected, person)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
