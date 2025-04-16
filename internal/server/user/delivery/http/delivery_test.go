package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/ds"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/delivery/http/dto"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	mocks "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/delivery/mocks"
)

func TestUserHandler_UpdateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dummyCookie := &config.Cookie{
		SessionName: "session_id",
		HTTPOnly:    true,
		Secure:      false,
		SameSite:    http.SameSiteLaxMode,
		Path:        "/",
	}
	ctx := config.WrapCookieContext(context.Background(), dummyCookie)

	mockUserSvc := mocks.NewMockUserServiceInterface(ctrl)
	mockSessionSvc := mocks.NewMockSessionServiceInterface(ctrl)

	updateReq := dto.UpdateUserRequest{
		Username:            "newusername",
		OldPassword:         "oldpassword",
		NewPassword:         "newpassword",
		RepeatedNewPassword: "newpassword",
		Avatar:              "newavatar.png",
	}
	jsonReq, err := json.Marshal(updateReq)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonReq))
	req = req.WithContext(ctx)

	req.AddCookie(&http.Cookie{Name: "session_id", Value: "oldsession"})
	rec := httptest.NewRecorder()

	mockSessionSvc.
		EXPECT().
		GetSession(gomock.Any(), "oldsession").
		Return("oldusername", nil).
		Times(1)

	hashedOld, err := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.DefaultCost)
	assert.NoError(t, err)
	existingUser := &models.User{
		Username:       "oldusername",
		HashedPassword: string(hashedOld),
		Avatar:         "oldavatar.png",
		CreatedAt:      time.Now().Add(-time.Hour),
	}
	mockUserSvc.
		EXPECT().
		GetUser(gomock.Any(), "oldusername").
		Return(existingUser, nil).
		Times(1)

	mockUserSvc.
		EXPECT().
		UpdateUser(gomock.Any(), "oldusername", gomock.Any()).
		Return(nil).
		Times(1)

	mockSessionSvc.
		EXPECT().
		DeleteSession(gomock.Any(), "oldsession").
		Return(nil).
		Times(1)

	mockSessionSvc.
		EXPECT().
		CreateSession(gomock.Any(), updateReq.Username).
		Return("newsession", nil).
		Times(1)

	handler := NewUserHandler(ctx, mockUserSvc, mockSessionSvc)
	handler.UpdateUser(rec, req)

	res := rec.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var resp ds.Response
	err = json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "successfully updated user", resp.Message)

	found := false
	for _, c := range res.Cookies() {
		if c.Name == dummyCookie.SessionName && c.Expires.After(time.Now()) {
			found = true
			assert.Equal(t, "newsession", c.Value)
		}
	}
	assert.True(t, found, "expected new session cookie to be set")
}

func TestUserHandler_UpdateUser_MissingSessionCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dummyCookie := &config.Cookie{
		SessionName: "session_id",
		HTTPOnly:    true,
		Secure:      false,
		SameSite:    http.SameSiteLaxMode,
		Path:        "/",
	}
	ctx := config.WrapCookieContext(context.Background(), dummyCookie)

	mockUserSvc := mocks.NewMockUserServiceInterface(ctrl)
	mockSessionSvc := mocks.NewMockSessionServiceInterface(ctrl)

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{}`))
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler := NewUserHandler(ctx, mockUserSvc, mockSessionSvc)
	handler.UpdateUser(rec, req)

	res := rec.Result()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestUserHandler_UpdateUser_ErrorPaths(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dummyCookie := &config.Cookie{
		SessionName: "session_id",
		HTTPOnly:    true,
		Secure:      false,
		SameSite:    http.SameSiteLaxMode,
		Path:        "/",
	}
	ctx := config.WrapCookieContext(context.Background(), dummyCookie)

	type testCase struct {
		name                    string
		sessionSvcSetup         func(mockSvc *mocks.MockSessionServiceInterface)
		userSvcSetup            func(mockSvc *mocks.MockUserServiceInterface)
		requestBody             string
		preExistingUserPassword string
		expectedStatus          int
		expectedError           string
	}

	validRequest := dto.UpdateUserRequest{
		Username:            "newusername",
		OldPassword:         "oldpassword",
		NewPassword:         "newpassword",
		RepeatedNewPassword: "newpassword",
		Avatar:              "newavatar.png",
	}
	validJSON, err := json.Marshal(validRequest)
	assert.NoError(t, err)

	tests := []testCase{
		{
			name: "GetSession error",
			sessionSvcSetup: func(m *mocks.MockSessionServiceInterface) {
				m.EXPECT().GetSession(gomock.Any(), "oldsession").Return("", errors.New(errs.ErrMsgSessionNotExists)).Times(1)
			},
			requestBody:    string(validJSON),
			expectedStatus: http.StatusUnauthorized,
			expectedError:  errs.ErrMsgSessionNotExists,
		},
		{
			name:           "JSON parsing error",
			requestBody:    "not a json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  errs.ErrParseJSONShort,
		},
		{
			name: "passwords mismatch",
			requestBody: func() string {
				req := validRequest
				req.RepeatedNewPassword = "different"
				b, err := json.Marshal(req)
				assert.NoError(t, err)
				return string(b)
			}(),
			expectedStatus: http.StatusBadRequest,
			expectedError:  errs.ErrPasswordsMismatchShort,
		},
		{
			name: "short password",
			requestBody: func() string {
				req := validRequest
				req.NewPassword = "short"
				req.RepeatedNewPassword = "short"
				b, _ := json.Marshal(req)
				return string(b)
			}(),
			expectedStatus: http.StatusBadRequest,
			expectedError:  errs.ErrInvalidPasswordShort,
		},
		{
			name: "long login",
			requestBody: func() string {
				req := validRequest
				req.Username = "longlonglonglonglonglonglonglonglonglonglonglonglogin"
				b, _ := json.Marshal(req)
				return string(b)
			}(),
			expectedStatus: http.StatusBadRequest,
			expectedError:  errs.ErrInvalidLoginShort,
		},
		{
			name: "GetUser error",
			sessionSvcSetup: func(m *mocks.MockSessionServiceInterface) {
				m.EXPECT().GetSession(gomock.Any(), "oldsession").Return("oldusername", nil).Times(1)
			},
			userSvcSetup: func(m *mocks.MockUserServiceInterface) {
				m.EXPECT().GetUser(gomock.Any(), "oldusername").Return(nil, errors.New(errs.ErrIncorrectLogin)).Times(1)
			},
			requestBody:    string(validJSON),
			expectedStatus: http.StatusBadRequest,
			expectedError:  "error getting user",
		},
		{
			name: "old password mismatch",
			sessionSvcSetup: func(m *mocks.MockSessionServiceInterface) {
				m.EXPECT().GetSession(gomock.Any(), "oldsession").Return("oldusername", nil).Times(1)
			},
			userSvcSetup: func(m *mocks.MockUserServiceInterface) {
				hashed, err := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.DefaultCost)
				assert.NoError(t, err)
				m.EXPECT().GetUser(gomock.Any(), "oldusername").Return(&models.User{
					Username:       "oldusername",
					HashedPassword: string(hashed),
					Avatar:         "oldavatar.png",
				}, nil).Times(1)
			},
			requestBody: func() string {
				req := validRequest
				req.OldPassword = "wrong old"
				b, _ := json.Marshal(req)
				return string(b)
			}(),
			expectedStatus: http.StatusBadRequest,
			expectedError:  errs.ErrInvalidPasswordShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserSvc := mocks.NewMockUserServiceInterface(ctrl)
			mockSessionSvc := mocks.NewMockSessionServiceInterface(ctrl)

			if tt.sessionSvcSetup != nil {
				tt.sessionSvcSetup(mockSessionSvc)
			} else {
				mockSessionSvc.EXPECT().GetSession(gomock.Any(), "oldsession").Return("oldusername", nil).AnyTimes()
			}

			if tt.userSvcSetup != nil {
				tt.userSvcSetup(mockUserSvc)
			}

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(tt.requestBody)))
			req = req.WithContext(ctx)

			req.AddCookie(&http.Cookie{Name: "session_id", Value: "oldsession"})
			rec := httptest.NewRecorder()

			handler := NewUserHandler(ctx, mockUserSvc, mockSessionSvc)
			handler.UpdateUser(rec, req)

			res := rec.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}
