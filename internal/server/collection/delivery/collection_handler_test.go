package delivery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/delivery/delivery_mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCollectionHandler_GetMainPageCollections(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := delivery_mocks.NewMockCollectionServiceInterface(ctrl)
	handler := NewCollectionHandler(mockService)

	tests := []struct {
		name         string
		mockSetup    func()
		expectedCode int
		expectedBody string
	}{
		{
			name: "OK. Get main page collections",
			mockSetup: func() {
				mockService.EXPECT().
					GetMainPageCollections(gomock.Any()).
					Return(mocks.MainPageCollections, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"Лучшие за всё время"`,
		},
		{
			name: "Fail. Collections not found",
			mockSetup: func() {
				mockService.EXPECT().
					GetMainPageCollections(gomock.Any()).
					Return(nil, errs.ErrCollectionNotExist)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: `collection does not exist`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req, _ := http.NewRequest("GET", "/collections/main", nil)
			rr := httptest.NewRecorder()

			handler.GetMainPageCollections(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
