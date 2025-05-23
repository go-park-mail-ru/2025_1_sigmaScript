package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	assert.NotEmpty(t, rec.Header().Get("Request-ID"))
}

func TestRequestWithLoggerMiddleware_UsesRequestIDHeader(t *testing.T) {
	customID := "custom-req-id-42"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(requestIDKey)
		assert.Equal(t, customID, requestID)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Request-ID", customID)
	rec := httptest.NewRecorder()

	RequestWithLoggerMiddleware(handler).ServeHTTP(rec, req)

	assert.Equal(t, customID, rec.Header().Get("Request-ID"))
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
