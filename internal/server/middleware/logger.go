package middleware

import (
	"context"
	"crypto/rand"
	"net/http"

	"github.com/rs/zerolog/log"
)

type requestIDctxKey int

const (
	requestIDKey requestIDctxKey = 0
	symbols                      = "abcdefghijklmnopqrstuvwxyz1234567890"
)

// RequestWithLoggerMiddleware places logger inside request context
func RequestWithLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("Request-ID")
		if requestID == "" {
			requestID = createRequestID()
		}

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		w.Header().Set("Request-ID", requestID)

		logger := log.With().Str("request_id", requestID).Caller().Logger()
		ctxtWithLogger := logger.WithContext(ctx)

		next.ServeHTTP(w, r.WithContext(ctxtWithLogger))
	})
}

// createRequestID generates uuid4-style request id
func createRequestID() string {
	output := make([]byte, 32)
	_, err := rand.Read(output)
	if err != nil {
		return ""
	}

	for pos := range output {
		output[pos] = symbols[uint8(output[pos])%uint8(len(symbols))]
	}

	// uuid4 styled string
	return string(output[0:8]) + "-" + string(output[8:12]) + "-4" +
		string(output[13:16]) + "-" + string(output[16:20]) + "-" + string(output[20:32])
}
