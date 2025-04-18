package delivery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	delivery_mocks "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/staff_person/delivery/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestStaffPersonHandler_GetPerson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := delivery_mocks.NewMockStaffPersonServiceInterface(ctrl)
	handler := NewStaffPersonHandler(mockService)

	keanuReeves := mocks.ExistingActors[11]

	tests := []struct {
		name         string
		personID     string
		mockSetup    func()
		expectedCode int
		expectedBody string
	}{
		{
			name:     "OK. Get Keanu Reeves",
			personID: "11",
			mockSetup: func() {
				mockService.EXPECT().
					GetPersonByID(gomock.Any(), 11).
					Return(&keanuReeves, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"full_name":"Киану Ривз"`,
		},
		{
			name:         "Fail. Invalid person ID (string)",
			personID:     "abc",
			mockSetup:    func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: `"error":"bad payload"`,
		},
		{
			name:     "Fail. Person not found",
			personID: "999",
			mockSetup: func() {
				mockService.EXPECT().
					GetPersonByID(gomock.Any(), 999).
					Return(nil, errs.ErrPersonNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: `"error":"not_found: person by this id not found"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req, _ := http.NewRequest("GET", "/person/"+tt.personID, nil)
			req = mux.SetURLVars(req, map[string]string{"person_id": tt.personID})

			rr := httptest.NewRecorder()
			handler.GetPerson(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
