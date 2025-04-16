package delivery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	delivery_mocks "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/delivery/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMovieHandler_GetMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := delivery_mocks.NewMockMovieServiceInterface(ctrl)
	handler := NewMovieHandler(mockService)

	fightClub := mocks.ExistingMovies[1]
	matrix := mocks.ExistingMovies[2]

	tests := []struct {
		name         string
		movieID      string
		mockSetup    func()
		expectedCode int
		expectedBody string
	}{
		{
			name:    "OK. Get Fight Club",
			movieID: "1",
			mockSetup: func() {
				mockService.EXPECT().
					GetMovieByID(gomock.Any(), 1).
					Return(&fightClub, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"name":"Бойцовский клуб"`,
		},
		{
			name:    "OK. Get Matrix",
			movieID: "2",
			mockSetup: func() {
				mockService.EXPECT().
					GetMovieByID(gomock.Any(), 2).
					Return(&matrix, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"name":"Матрица"`,
		},
		{
			name:         "Fail. Invalid movie ID (string)",
			movieID:      "abc",
			mockSetup:    func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: `"error":"bad payload"`,
		},
		{
			name:    "Fail. Movie not found",
			movieID: "999",
			mockSetup: func() {
				mockService.EXPECT().
					GetMovieByID(gomock.Any(), 999).
					Return(nil, errs.ErrMovieNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: `"error":"not_found: movie by this id not found"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req, _ := http.NewRequest("GET", "/movie/"+tt.movieID, nil)
			req = mux.SetURLVars(req, map[string]string{"movie_id": tt.movieID})

			rr := httptest.NewRecorder()
			handler.GetMovie(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
