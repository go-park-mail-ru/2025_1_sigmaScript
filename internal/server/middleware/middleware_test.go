package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/common"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	csrftoken "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/csrf_token"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRequestWithLoggerMiddleware_AssignsRequestID(t *testing.T) {
	handlerCalled := false

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		requestID := r.Context().Value(requestIDKey)
		assert.NotNil(t, requestID)
		w.WriteHeader(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{"msg":"hello"}`))
	rec := httptest.NewRecorder()

	middleware := RequestWithLoggerMiddleware(handler)
	middleware.ServeHTTP(rec, req)

	assert.True(t, handlerCalled)
	assert.Equal(t, http.StatusTeapot, rec.Code)
	assert.NotEmpty(t, rec.Header().Get("X-Request-ID"))
}

func TestRequestWithLoggerMiddleware_UsesRequestIDHeader(t *testing.T) {
	customID := "custom-req-id-42"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(requestIDKey)
		assert.Equal(t, customID, requestID)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", customID)
	rec := httptest.NewRecorder()

	RequestWithLoggerMiddleware(handler).ServeHTTP(rec, req)

	assert.Equal(t, customID, rec.Header().Get("X-Request-ID"))
}

func TestRequestWithLoggerMiddleware_ReadsBody(t *testing.T) {
	const bodyStr = `{"test":"value"}`

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, bodyStr, string(data))
	})

	req := httptest.NewRequest(http.MethodPost, "/body", strings.NewReader(bodyStr))
	rec := httptest.NewRecorder()

	RequestWithLoggerMiddleware(handler).ServeHTTP(rec, req)
}

func TestMiddlewareCors_SetsCORSHeaders(t *testing.T) {
	viper.Set(kinolkHostEnv, "http://localhost:3000")
	viper.Set(kinolkAllowedMethodsEnv, "GET, POST, OPTIONS")
	viper.Set(kinolkAllowCredentialsEnv, "true")
	viper.Set(kinolkAllowedHeadersEnv, "Content-Type, Authorization")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	MiddlewareCors(handler).ServeHTTP(rec, req)

	assert.Equal(t, "http://localhost:3000", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, OPTIONS", rec.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "Content-Type, Authorization", rec.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestMiddlewareCors_OPTIONSMethod(t *testing.T) {
	viper.Set(kinolkHostEnv, "http://localhost:3000")
	viper.Set(kinolkAllowedMethodsEnv, "GET, POST, OPTIONS")
	viper.Set(kinolkAllowCredentialsEnv, "true")
	viper.Set(kinolkAllowedHeadersEnv, "Content-Type")

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	rec := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	MiddlewareCors(handler).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.False(t, handlerCalled)
}

func TestPreventPanicMiddleware_HandlesPanic(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	PreventPanicMiddleware(handler).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Internal server error")
}

func TestPreventPanicMiddleware_NoPanic(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	PreventPanicMiddleware(handler).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTeapot, rec.Code)
}

func TestMiddlewareCSRF_PostMethod(t *testing.T) {
	viper.Set(kinolkHostEnv, "http://localhost:3000")
	viper.Set(kinolkAllowedMethodsEnv, "GET, POST, OPTIONS")
	viper.Set(kinolkAllowCredentialsEnv, "true")
	viper.Set(kinolkAllowedHeadersEnv, "Content-Type")

	req := httptest.NewRequest(http.MethodPost, "/movie/1/reviews", nil)
	rec := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	CsrfTokenMiddleware(handler).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.False(t, handlerCalled)
}

func TestMiddlewareCSRF_GetMethod(t *testing.T) {
	viper.Set(kinolkHostEnv, "http://localhost:3000")
	viper.Set(kinolkAllowedMethodsEnv, "GET, POST, OPTIONS")
	viper.Set(kinolkAllowCredentialsEnv, "true")
	viper.Set(kinolkAllowedHeadersEnv, "Content-Type")

	req := httptest.NewRequest(http.MethodGet, "/movie/1/reviews", nil)
	rec := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	CsrfTokenMiddleware(handler).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, handlerCalled)
}

func TestMiddlewareCSRF_PostMethod_AuthRoute(t *testing.T) {
	viper.Set(kinolkHostEnv, "http://localhost:3000")
	viper.Set(kinolkAllowedMethodsEnv, "GET, POST, OPTIONS")
	viper.Set(kinolkAllowCredentialsEnv, "true")
	viper.Set(kinolkAllowedHeadersEnv, "Content-Type")

	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	rec := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	CsrfTokenMiddleware(handler).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, handlerCalled)
}

func TestMiddlewareCSRF_PostMethod_NoCSRFHeaderOrCookie(t *testing.T) {
	viper.Set(kinolkHostEnv, "http://localhost:3000")
	viper.Set(kinolkAllowedMethodsEnv, "GET, POST, OPTIONS")
	viper.Set(kinolkAllowCredentialsEnv, "true")
	viper.Set(kinolkAllowedHeadersEnv, "Content-Type")

	req := httptest.NewRequest(http.MethodPost, "/movie/1/reviews", nil)
	rec := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	CsrfTokenMiddleware(handler).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.False(t, handlerCalled)
}

func TestMiddlewareCSRF_PostReview_MissingSessionCookie(t *testing.T) {

	dummyCookie := &config.Cookie{
		SessionName: common.CSRF_TOKEN_NAME,
		HTTPOnly:    true,
		Secure:      false,
		SameSite:    http.SameSiteLaxMode,
		Path:        "/",
	}

	ctx := config.WrapCookieContext(context.Background(), dummyCookie)

	req := httptest.NewRequest(http.MethodPost, "/movie/1/reviews", strings.NewReader(`{}`))
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	CsrfTokenMiddleware(handler).ServeHTTP(rec, req)

	res := rec.Result()

	var parcedResultBody jsonutil.ErrorResponse
	json.NewDecoder(rec.Result().Body).Decode(&parcedResultBody)

	assert.Equal(t, http.StatusForbidden, res.StatusCode)
	assert.False(t, handlerCalled)
	assert.Equal(t, errs.ErrUnauthorizedShort, parcedResultBody.Error)
	assert.Equal(t, errs.ErrMsgBadCSRFToken, parcedResultBody.Message)
}

func TestMiddlewareCSRF_PostReview_MissingCSRFHeader(t *testing.T) {
	viper.Set(kinolkHostEnv, "http://localhost:3000")
	viper.Set(kinolkAllowedMethodsEnv, "GET, POST, OPTIONS")
	viper.Set(kinolkAllowCredentialsEnv, "true")
	viper.Set(kinolkAllowedHeadersEnv, "Content-Type")

	req := httptest.NewRequest(http.MethodPost, "/movie/1/reviews", nil)
	rec := httptest.NewRecorder()

	cookie := &http.Cookie{Name: common.CSRF_TOKEN_NAME, Value: "some_value"}
	req.AddCookie(cookie)

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	CsrfTokenMiddleware(handler).ServeHTTP(rec, req)

	res := rec.Result()

	var parcedResultBody jsonutil.ErrorResponse
	json.NewDecoder(rec.Result().Body).Decode(&parcedResultBody)

	assert.Equal(t, http.StatusForbidden, res.StatusCode)
	assert.False(t, handlerCalled)
	assert.Equal(t, errs.ErrUnauthorizedShort, parcedResultBody.Error)
	assert.Equal(t, errs.ErrMsgBadCSRFToken, parcedResultBody.Message)
}

func TestMiddlewareCSRF_PostReview_Success(t *testing.T) {
	viper.Set(kinolkHostEnv, "http://localhost:3000")
	viper.Set(kinolkAllowedMethodsEnv, "GET, POST, OPTIONS")
	viper.Set(kinolkAllowCredentialsEnv, "true")
	viper.Set(kinolkAllowedHeadersEnv, "Content-Type, X-CSRF-Token")

	req := httptest.NewRequest(http.MethodPost, "/movie/1/reviews", nil)
	rec := httptest.NewRecorder()

	newCSRFToken, errCSRF := csrftoken.GenerateCSRFToken()
	assert.NoError(t, errCSRF)
	req.Header.Set("X-CSRF-TOKEN", newCSRFToken)

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	CsrfTokenMiddleware(handler).ServeHTTP(rec, req)

	res := rec.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, handlerCalled)
}
