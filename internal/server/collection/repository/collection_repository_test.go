package repository

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCollectionRepository_GetMainPageCollectionsFromRepo(t *testing.T) {
	tests := []struct {
		name        string
		dbSetup     mocks.Collections
		expected    mocks.Collections
		expectedErr error
	}{
		{
			name:        "Success - Get collections",
			dbSetup:     mocks.MainPageCollections,
			expected:    mocks.MainPageCollections,
			expectedErr: nil,
		},
		{
			name:        "Success - Empty collections",
			dbSetup:     make(mocks.Collections),
			expected:    make(mocks.Collections),
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewCollectionRepository(&tt.dbSetup)
			collections, err := repo.GetMainPageCollectionsFromRepo(context.Background())

			assert.Equal(t, tt.expected, collections)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
