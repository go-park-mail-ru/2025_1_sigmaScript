package handlers

import (
	"bytes"
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
)

var cfg, err = config.New()

var registration = models.RegisterData{
	Username:         "george",
	Password:         "georgesPassword10",
	RepeatedPassword: "georgesPassword10",
}

var login = models.LoginData{
	Username: "george",
	Password: "georgesPassword10",
}

func getResponseRequest(t *testing.T, method, target string, data any) (*httptest.ResponseRecorder, *http.Request) {
	var req *http.Request
	jsonData, err := json.Marshal(data)
	require.NoError(t, err, errs.ErrParseJSON, errs.ErrEncodeJSON)
	jsonReader := bytes.NewReader(jsonData)
	if method == "GET" {
		req = httptest.NewRequest(method, target, nil)
	} else {
		req = httptest.NewRequest(method, target, jsonReader)
	}
	rr := httptest.NewRecorder()
	return rr, req
}

func registerUser(t *testing.T, auth *AuthHandler, data any) *httptest.ResponseRecorder {
	rr, req := getResponseRequest(t, "POST", "/auth/register/", data)
	auth.RegisterHandler(rr, req)
	return rr
}

func loginUser(t *testing.T, auth *AuthHandler, data any) (*httptest.ResponseRecorder, *http.Cookie) {
	rr, req := getResponseRequest(t, "POST", "/auth/login/", data)
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
	rr, req := getResponseRequest(t, "POST", "/auth/logout/", data)
	req.AddCookie(cookie)
	auth.LogoutHandler(rr, req)
	return rr
}

func assertHeaders(t *testing.T, code int, rr *httptest.ResponseRecorder) {
	assert.Equal(t, code, rr.Code, errs.ErrWrongResponseCode)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), errs.ErrWrongHeaders)
}

func checkResponceError(t *testing.T, rr *httptest.ResponseRecorder, expectedMessage string) {
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
func TestRegisterOK(t *testing.T) {
	authHandler := NewAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie))
	rr := registerUser(t, authHandler, registration)
	assertHeaders(t, http.StatusOK, rr)
	checkResponseMessage(t, rr, messages.SuccessfulRegister)
}

func TestRegisterAlreadyExists(t *testing.T) {
	authHandler := NewAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie))
	registerUser(t, authHandler, registration)
	rr := registerUser(t, authHandler, registration)
	assertHeaders(t, http.StatusBadRequest, rr)
	checkResponceError(t, rr, errs.ErrAlreadyExistsShort)
}

func TestRegisterMissmatch(t *testing.T) {
	authHandler := NewAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie))
	data := models.RegisterData{
		Username:         "Mismatch",
		Password:         "passwordFiRst1",
		RepeatedPassword: "SeCond2password",
	}
	rr := registerUser(t, authHandler, data)
	assertHeaders(t, http.StatusBadRequest, rr)
	checkResponceError(t, rr, errs.ErrPasswordsMismatchShort)
}

func TestRegisterInvalidPassword(t *testing.T) {
	authHandler := NewAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie))
	data := models.RegisterData{
		Username:         "invalidPassword",
		Password:         "123",
		RepeatedPassword: "123",
	}
	rr := registerUser(t, authHandler, data)
	assertHeaders(t, http.StatusBadRequest, rr)
	checkResponceError(t, rr, errs.ErrInvalidPasswordShort+": "+errs.ErrPasswordTooShort)
}

func TestLoginOK(t *testing.T) {
	authHandler := NewAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie))
	registerUser(t, authHandler, registration)
	rr, cookie := loginUser(t, authHandler, login)
	assertHeaders(t, http.StatusOK, rr)
	checkResponseMessage(t, rr, messages.SuccessfulLogin)
	assert.NotNil(t, cookie, errs.ErrSessionNotExists)
	assert.NotEmpty(t, cookie.Value, errs.ErrCookieEmpty)
	assert.True(t, cookie.HttpOnly, errs.ErrCookieHttpOnly)
}

func TestLoginFail(t *testing.T) {
	authHandler := NewAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie))
	rr, cookie := loginUser(t, authHandler, login)
	assertHeaders(t, http.StatusUnauthorized, rr)
	expected := "not_found: crypto/bcrypt: hashedSecret too short to be a bcrypted password"
	checkResponceError(t, rr, expected)
	assert.Empty(t, cookie, errs.ErrSessionCreated)
}

func TestLogout(t *testing.T) {
	authHandler := NewAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie))
	registerUser(t, authHandler, registration)
	_, cookie := loginUser(t, authHandler, login)
	rr := logoutUser(t, authHandler, login, cookie)
	assertHeaders(t, http.StatusOK, rr)
	checkResponseMessage(t, rr, messages.SuccessfulLogout)
	assert.True(t, cookie.MaxAge <= 0, errs.ErrCookieExpire)
}

func TestLogoutNoCookie(t *testing.T) {
	authHandler := NewAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie))
	registerUser(t, authHandler, registration)
	_, cookie := loginUser(t, authHandler, login)
	cookie.Name = "something_else"
	rr := logoutUser(t, authHandler, login, cookie)
	assertHeaders(t, http.StatusUnauthorized, rr)
	checkResponceError(t, rr, errs.ErrUnauthorizedShort+": "+http.ErrNoCookie.Error())
}

func TestLogoutNoSession(t *testing.T) {
	var authHandler = NewAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie))
	registerUser(t, authHandler, registration)
	_, cookie := loginUser(t, authHandler, login)

	for k := range authHandler.sessions {
		delete(authHandler.sessions, k)
	}

	rr := logoutUser(t, authHandler, login, cookie)
	assertHeaders(t, http.StatusNotFound, rr)
	checkResponceError(t, rr, errs.ErrSessionNotExistsShort)
}
