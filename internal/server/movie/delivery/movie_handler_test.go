package delivery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/delivery/delivery_mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMovieHandler_GetMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 1. Подготавливаем моки
	mockService := delivery_mocks.NewMockMovieServiceInterface(ctrl)
	handler := NewMovieHandler(mockService)

	// 2. Тестовые данные
	fightClub := mocks.ExistingMovies[0] // "Бойцовский клуб"

	tests := []struct {
		name         string
		movieID      string // ID из URL (может быть некорректным)
		mockSetup    func() // Настройка моков сервиса
		expectedCode int    // Ожидаемый HTTP-статус
		expectedBody string // Ожидаемое тело ответа (частичное совпадение)
	}{
		// Успешные сценарии
		{
			name:    "Success - Get existing movie",
			movieID: "0",
			mockSetup: func() {
				mockService.EXPECT().
					GetMovieByID(gomock.Any(), 0).
					Return(&fightClub, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"name":"Бойцовский клуб"`,
		},

		// Ошибки клиента
		{
			name:         "Fail - Invalid movie ID (string)",
			movieID:      "abc",
			mockSetup:    func() {}, // Сервис не должен вызываться
			expectedCode: http.StatusBadRequest,
			expectedBody: `"error":"bad payload"`,
		},
		{
			name:    "Fail - Movie not found",
			movieID: "999",
			mockSetup: func() {
				mockService.EXPECT().
					GetMovieByID(gomock.Any(), 999).
					Return(nil, errs.ErrMovieNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: `not_found`,
		},
		{
			name:    "Fail - Internal server error",
			movieID: "0",
			mockSetup: func() {
				mockService.EXPECT().
					GetMovieByID(gomock.Any(), 0).
					Return(nil, assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `"error":"something went wrong"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Настраиваем моки
			tt.mockSetup()

			// 2. Создаем HTTP-запрос
			req, _ := http.NewRequest("GET", "/movie/"+tt.movieID, nil)
			req = mux.SetURLVars(req, map[string]string{"movie_id": tt.movieID})

			// 3. Записываем ответ
			rr := httptest.NewRecorder()

			// 4. Вызываем обработчик
			handler.GetMovie(rr, req)

			// 5. Проверяем статус код
			assert.Equal(t, tt.expectedCode, rr.Code)

			// 6. Проверяем тело ответа (для JSON)
			if tt.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
