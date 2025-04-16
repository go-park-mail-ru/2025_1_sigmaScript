package delivery

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	authMocks "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/delivery/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	movieMocks "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/delivery/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestReviewHandler_GetAllReviewsOfMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := authMocks.NewMockUserServiceInterface(ctrl)
	mockSessionService := authMocks.NewMockSessionServiceInterface(ctrl)
	mockMovieService := movieMocks.NewMockMovieServiceInterface(ctrl)

	handler := NewReviewHandler(mockUserService, mockSessionService, mockMovieService)

	fightClubReviews := mocks.ExistingMovies[1].Reviews
	emptyReviews := mocks.ExistingMovies[2].Reviews

	tests := []struct {
		name         string
		movieID      string
		mockSetup    func()
		expectedCode int
		expectedBody string
	}{
		{
			name:    "Success. Get reviews for Fight Club",
			movieID: "1",
			mockSetup: func() {
				mockMovieService.EXPECT().
					GetAllReviewsOfMovieByID(gomock.Any(), 1).
					Return(&fightClubReviews, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"KinoKritik77"`,
		},
		{
			name:    "Success. Get empty reviews for Matrix",
			movieID: "2",
			mockSetup: func() {
				mockMovieService.EXPECT().
					GetAllReviewsOfMovieByID(gomock.Any(), 2).
					Return(&emptyReviews, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `Потрясающая игра актеров и неожиданный финал.`,
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
				mockMovieService.EXPECT().
					GetAllReviewsOfMovieByID(gomock.Any(), 999).
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

			req, _ := http.NewRequest("GET", "/movies/"+tt.movieID+"/reviews", nil)
			req = mux.SetURLVars(req, map[string]string{"movie_id": tt.movieID})

			rr := httptest.NewRecorder()
			handler.GetAllReviewsOfMovie(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestReviewHandler_CreateReview(t *testing.T) {
	// TODO
	// Само создание отзыва

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := authMocks.NewMockUserServiceInterface(ctrl)
	mockSessionService := authMocks.NewMockSessionServiceInterface(ctrl)
	mockMovieService := movieMocks.NewMockMovieServiceInterface(ctrl)

	handler := NewReviewHandler(mockUserService, mockSessionService, mockMovieService)

	tests := []struct {
		name           string
		movieID        string
		sessionID      string
		requestBody    string
		mockSetup      func()
		expectedCode   int
		expectedBody   string
		expectedHeader string
	}{
		{
			name:         "Fail. No session cookie",
			movieID:      "1",
			sessionID:    "",
			requestBody:  `{}`,
			mockSetup:    func() {},
			expectedCode: http.StatusUnauthorized,
			expectedBody: `"error":"unauthorized"`,
		},
		{
			name:      "Fail. Invalid session",
			movieID:   "1",
			sessionID: "invalid_session",
			requestBody: `{
				"score": 9,
				"review_text": "Great movie!"
			}`,
			mockSetup: func() {
				mockSessionService.EXPECT().
					GetSession(gomock.Any(), "invalid_session").
					Return("", errs.ErrSessionNotExists)
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: `session does not exist`,
		},
		{
			name:      "Fail. Invalid score (too low)",
			movieID:   "1",
			sessionID: "valid_session",
			requestBody: `{
				"score": 0,
				"review_text": "Bad movie!"
			}`,
			mockSetup: func() {
				mockSessionService.EXPECT().
					GetSession(gomock.Any(), "valid_session").
					Return("test_user", nil)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `"error":"bad payload"`,
		},
		{
			name:      "Fail. Invalid score (too high)",
			movieID:   "1",
			sessionID: "valid_session",
			requestBody: `{
				"score": 11,
				"review_text": "Awesome movie!"
			}`,
			mockSetup: func() {
				mockSessionService.EXPECT().
					GetSession(gomock.Any(), "valid_session").
					Return("test_user", nil)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `"error":"bad payload"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req, _ := http.NewRequest("POST", "/movies/"+tt.movieID+"/reviews", bytes.NewBufferString(tt.requestBody))
			req = mux.SetURLVars(req, map[string]string{"movie_id": tt.movieID})

			if tt.sessionID != "" {
				req.AddCookie(&http.Cookie{Name: "session_id", Value: tt.sessionID})
			}

			rr := httptest.NewRecorder()
			handler.CreateReview(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedBody)
			}
			if tt.expectedHeader != "" {
				assert.Equal(t, tt.expectedHeader, rr.Header().Get("Content-Type"))
			}
		})
	}
}
