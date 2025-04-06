package middleware

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type requestIDctxKey int

const (
	requestIDKey  requestIDctxKey = 0
	symbols                       = "abcdefghijklmnopqrstuvwxyz1234567890"
	logMiddleware                 = "request logger middleware"
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

		customResponseWriter := NewResponseWriterWithStatus(w, r.URL.Path)

		requestStartTime := time.Now()

		next.ServeHTTP(customResponseWriter, r.WithContext(ctxtWithLogger))
		status := customResponseWriter.Status
		logRequestData(r, requestStartTime, logMiddleware, requestID, status, requestURLPath(w, r))
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

type responseWriterWithStatus struct {
	http.ResponseWriter
	Status  int
	URLPath string
}

// this WriteHeader method captures the status code and calls the original WriteHeader.
func (rws *responseWriterWithStatus) WriteHeader(statusCode int) {
	rws.Status = statusCode
	rws.ResponseWriter.WriteHeader(statusCode)
}

// wrap response
func NewResponseWriterWithStatus(w http.ResponseWriter, path string) *responseWriterWithStatus {
	return &responseWriterWithStatus{
		ResponseWriter: w,
		Status:         http.StatusOK,
		URLPath:        path,
	}
}

func requestURLPath(w http.ResponseWriter, r *http.Request) string {
	urlPath := mux.CurrentRoute(r)
	if urlPath == nil {
		http.Error(w, "Route not found", http.StatusNotFound)
		return ""
	}

	return urlPath.GetName()
}

func logRequestData(r *http.Request, start time.Time, msg string, requestID string, status int, path string) {
	var bodyCopy bytes.Buffer
	duration := time.Since(start)

	tee := io.TeeReader(r.Body, &bodyCopy)
	r.Body = io.NopCloser(&bodyCopy)
	bodyBytes, err := io.ReadAll(tee)
	if err != nil {
		log.Error().Err(err).Msg(errs.ErrBadPayload)
	}

	log.Info().
		Str("method", r.Method).
		Str("remote_addr", r.RemoteAddr).
		Str("url", path).
		Str("request_id", requestID).
		Bytes("body", bodyBytes).
		Dur("work_time", duration).
		Int("status", status).
		Str("user_agent", r.UserAgent()).
		Str("host", r.Host).
		Str("real_ip", getRealIPAddr(r)).
		Int64("content_length", r.ContentLength).
		Str("start_time", start.Format(time.RFC3339)).
		Str("duration_human_readable", duration.String()).
		Int64("duration_ms", duration.Milliseconds()).
		Msg(msg)
}

func getRealIPAddr(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		parts := strings.Split(ip, ",")
		if len(parts) > 0 {
			realIP := strings.TrimSpace(parts[0])
			if net.ParseIP(realIP) != nil {
				return realIP
			}
		}
	}

	hostIPAddr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return hostIPAddr
}
