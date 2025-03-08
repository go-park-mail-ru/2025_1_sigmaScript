package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/ds"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/messages"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var cfg, err = config.New()

const (
	RegisterUrl = "/auth/register/"
	LoginUrl    = "/auth/login/"
	LogoutUrl   = "/auth/logout/"
)

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
			expectedResponse: errs.ErrInvalidPasswordShort + ": " + errs.ErrPasswordTooShort,
			expectedStatus:   http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authHandler := NewAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie))
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
			expectedResponse: errs.ErrIncorrectLoginOrPasswordShort + ": " + bcrypt.ErrHashTooShort.Error(),
			expectedStatus:   http.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authHandler := NewAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie))
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
			expectedResponse:   errs.ErrUnauthorizedShort + ": " + http.ErrNoCookie.Error(),
			expectedStatus:     http.StatusUnauthorized,
		},
		{
			name:               "LogoutNoCookie",
			cookieCheckExpired: false,
			cookieChangeName:   false,
			deleteSessions:     true,
			checkMessage:       false,
			expectedResponse:   errs.ErrSessionNotExistsShort,
			expectedStatus:     http.StatusNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authHandler := NewAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie))
			registerUser(t, authHandler, registration)
			_, cookie := loginUser(t, authHandler, login)
			if test.cookieChangeName {
				cookie.Name = "something_else"
			}
			if test.deleteSessions {
				for k := range authHandler.sessions {
					delete(authHandler.sessions, k)
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
	auth.RegisterHandler(rr, req)
	return rr
}

func loginUser(t *testing.T, auth *AuthHandler, data any) (*httptest.ResponseRecorder, *http.Cookie) {
	rr, req := getResponseRequest(t, "POST", LoginUrl, data)
	auth.LoginHandler(rr, req)

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
	auth.LogoutHandler(rr, req)
	return rr
}

func checkResponseError(t *testing.T, rr *httptest.ResponseRecorder, expectedMessage string) {
	var got jsonutil.ErrorResponse
	expected := expectedMessage
	err = json.NewDecoder(rr.Body).Decode(&got)
	require.NoError(t, err, errs.ErrParseJSON)
	assert.True(t, reflect.DeepEqual(got.Error, expected))
}

func checkResponseMessage(t *testing.T, rr *httptest.ResponseRecorder, expectedMessage string) {
	var got ds.Response
	expected := expectedMessage
	err = json.NewDecoder(rr.Body).Decode(&got)
	require.NoError(t, err, errs.ErrParseJSON)
	assert.True(t, reflect.DeepEqual(got.Message, expected))
}
