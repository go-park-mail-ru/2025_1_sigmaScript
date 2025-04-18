package service

import (
	"context"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	service_mocks "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/service/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCollectionService_GetMainPageCollections(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := service_mocks.NewMockCollectionRepositoryInterface(ctrl)
	service := NewCollectionService(mockRepo)

	tests := []struct {
		name        string
		mockSetup   func()
		expected    mocks.Collections
		expectedErr error
	}{
		{
			name: "OK. Get collections",
			mockSetup: func() {
				mockRepo.EXPECT().
					GetMainPageCollectionsFromRepo(gomock.Any()).
					Return(mocks.MainPageCollections, nil)
			},
			expected:    mocks.MainPageCollections,
			expectedErr: nil,
		},
		{
			name: "Fail. Empty collections",
			mockSetup: func() {
				mockRepo.EXPECT().
					GetMainPageCollectionsFromRepo(gomock.Any()).
					Return(nil, errs.ErrCollectionNotExist)
			},
			expected:    nil,
			expectedErr: errs.ErrCollectionNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			collections, err := service.GetMainPageCollections(context.Background())

			assert.Equal(t, tt.expected, collections)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
