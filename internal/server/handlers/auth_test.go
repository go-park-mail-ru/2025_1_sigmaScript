package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/messages"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	synccredmap "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/sync_cred_map"
	"github.com/stretchr/testify/assert"
)

var cfg, err = config.New()

const (
	RegisterUrl = "/auth/register/"
	LoginUrl    = "/auth/login/"
	LogoutUrl   = "/auth/logout/"
)

func newAuthHandler(ctx context.Context) *AuthHandler {
	return &AuthHandler{
		users:    *synccredmap.NewSyncCredentialsMap(),
		sessions: *synccredmap.NewSyncCredentialsMap(),
		cfg:      config.FromCookieContext(ctx),
	}
}

func TestRegister(t *testing.T) {
	registration := models.RegisterData{
		Username:         "guestForTest",
		Password:         "guestPassword10",
		RepeatedPassword: "guestPassword10",
	}

	tests := []struct {
		name             string
		registerUser     bool
		data             models.RegisterData
		expectedStatus   int
		checkMessage     bool
		expectedResponse string
	}{
		{
			name:             "RegisterOK",
			registerUser:     false,
			data:             registration,
			checkMessage:     true,
			expectedResponse: messages.SuccessfulRegister,
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "RegisterAlreadyExists",
			registerUser:     true,
			data:             registration,
			checkMessage:     false,
			expectedResponse: errs.ErrAlreadyExistsShort,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:         "RegisterMismatch",
			registerUser: false,
			data: models.RegisterData{
				Username:         "guestForTest",
				Password:         "passwordFirST1",
				RepeatedPassword: "SeCond2password",
			},
			checkMessage:     false,
			expectedResponse: errs.ErrPasswordsMismatchShort,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:         "RegisterInvalidPassword",
			registerUser: false,
			checkMessage: false,
			data: models.RegisterData{
				Username:         "guestForTest",
				Password:         "123",
				RepeatedPassword: "123",
			},
			expectedResponse: fmt.Sprintf("%s: %s", errs.ErrInvalidPasswordShort, "Password too short"),
			expectedStatus:   http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authHandler := newAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie))
			if test.registerUser {
				registerUser(t, authHandler, registration)
			}
			rr := registerUser(t, authHandler, test.data)
			assertHeaders(t, test.expectedStatus, rr)
			if test.checkMessage {
				checkResponseMessage(t, rr, test.expectedResponse)
			} else {
				checkResponseError(t, rr, test.expectedResponse)
			}
		})
	}
}
func TestLogin(t *testing.T) {
	registration := models.RegisterData{
		Username:         "guestForTest",
		Password:         "guestPassword10",
		RepeatedPassword: "guestPassword10",
	}

	login := models.LoginData{
		Username: "guestForTest",
		Password: "guestPassword10",
	}

	tests := []struct {
		name             string
		registerUser     bool
		cookiesEnabled   bool
		checkMessage     bool
		expectedResponse string
		expectedStatus   int
	}{
		{
			name:             "LoginOK",
			registerUser:     true,
			cookiesEnabled:   true,
			checkMessage:     true,
			expectedResponse: messages.SuccessfulLogin,
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "LoginFail",
			registerUser:     false,
			cookiesEnabled:   false,
			checkMessage:     false,
			expectedResponse: fmt.Sprintf("%s: %s", errs.ErrIncorrectLoginOrPasswordShort, "user with this login does not exist"),
			expectedStatus:   http.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authHandler := newAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie))
			if test.registerUser {
				registerUser(t, authHandler, registration)
			}
			rr, cookie := loginUser(t, authHandler, login)
			assertHeaders(t, test.expectedStatus, rr)
			if test.checkMessage {
				checkResponseMessage(t, rr, test.expectedResponse)
			} else {
				checkResponseError(t, rr, test.expectedResponse)
			}
			if test.cookiesEnabled {
				assert.NotNil(t, cookie, errs.ErrSessionNotExists)
				assert.NotEmpty(t, cookie.Value, errs.ErrCookieEmpty)
				assert.True(t, cookie.HttpOnly, errs.ErrCookieHttpOnly)
			} else {
				assert.Empty(t, cookie, errs.ErrSessionCreated)
			}
		})
	}
}
func TestLogout(t *testing.T) {
	registration := models.RegisterData{
		Username:         "guestForTest",
		Password:         "guestPassword10",
		RepeatedPassword: "guestPassword10",
	}

	login := models.LoginData{
		Username: "guestForTest",
		Password: "guestPassword10",
	}

	tests := []struct {
		name               string
		cookieCheckExpired bool
		cookieChangeName   bool
		deleteSessions     bool
		checkMessage       bool
		expectedResponse   string
		expectedStatus     int
	}{
		{
			name:               "LogoutOK",
			cookieCheckExpired: true,
			cookieChangeName:   false,
			deleteSessions:     false,
			checkMessage:       true,
			expectedResponse:   messages.SuccessfulLogout,
			expectedStatus:     http.StatusOK,
		},
		{
			name:               "LogoutNoCookie",
			cookieCheckExpired: false,
			cookieChangeName:   true,
			deleteSessions:     false,
			checkMessage:       false,
			expectedResponse:   fmt.Sprintf("%s: %s", errs.ErrUnauthorizedShort, "http: named cookie not present"),
			expectedStatus:     http.StatusUnauthorized,
		},
		{
			name:               "LogoutNoCookie",
			cookieCheckExpired: false,
			cookieChangeName:   false,
			deleteSessions:     true,
			checkMessage:       false,
			expectedResponse:   fmt.Sprintf("%s: %s", errs.ErrSessionNotExistsShort, "session does not exist"),
			expectedStatus:     http.StatusNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authHandler := newAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie))
			registerUser(t, authHandler, registration)
			_, cookie := loginUser(t, authHandler, login)
			if test.cookieChangeName {
				cookie.Name = "something_else"
			}
			if test.deleteSessions {
				for k := range authHandler.sessions.Map() {
					delete(authHandler.sessions.Map(), k)
				}
			}
			rr := logoutUser(t, authHandler, login, cookie)
			assertHeaders(t, test.expectedStatus, rr)
			if test.checkMessage {
				checkResponseMessage(t, rr, test.expectedResponse)
			} else {
				checkResponseError(t, rr, test.expectedResponse)
			}
			if test.cookieCheckExpired {
				assert.True(t, cookie.MaxAge <= 0, errs.ErrCookieExpire)
			}
		})
	}
}

func registerUser(t *testing.T, auth *AuthHandler, data any) *httptest.ResponseRecorder {
	rr, req := getResponseRequest(t, "POST", RegisterUrl, data)
	auth.Register(rr, req)
	return rr
}

func loginUser(t *testing.T, auth *AuthHandler, data any) (*httptest.ResponseRecorder, *http.Cookie) {
	rr, req := getResponseRequest(t, "POST", LoginUrl, data)
	auth.Login(rr, req)

	var sessionCookie *http.Cookie
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == auth.cfg.SessionName {
			sessionCookie = cookie
			break
		}
	}
	return rr, sessionCookie
}

func logoutUser(t *testing.T, auth *AuthHandler, data any, cookie *http.Cookie) *httptest.ResponseRecorder {
	rr, req := getResponseRequest(t, "POST", LogoutUrl, data)
	req.AddCookie(cookie)
	auth.Logout(rr, req)
	return rr
}
