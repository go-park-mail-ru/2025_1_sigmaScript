package interceptors

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type requestIDctxKey int

const (
	requestIDKey requestIDctxKey = 0
)

func LoggerInterceptor(
	ctx context.Context,
	request interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Info().Msg("no metadata found in context")
	}

	requestID := md.Get("request_id")
	if len(requestID) > 0 {
		ctx = context.WithValue(ctx, requestIDKey, requestID[0])
		logger := log.With().Str("request_id", requestID[0]).Caller().Logger()
		ctx = logger.WithContext(ctx)
	} else {
		log.Info().Msg("empty request_id in metadata")
	}

	return handler(ctx, request)
}

// AccessLogInterceptor logs the request
func AccessLogInterceptor(
	ctx context.Context,
	request interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()
	logger := log.Ctx(ctx)

	defer func() {
		if r := recover(); r != nil {
			duration := time.Since(start)

			logger.Info().
				Str("method", info.FullMethod).
				Interface("request_body", request).
				Interface("response_body", resp).
				Str("start_time", start.Format(time.RFC3339)).
				Str("duration_human_readable", duration.String()).
				Int64("duration_ms", duration.Milliseconds()).
				Interface("panic", r).
				Msg("panic occurred")

			err = status.Errorf(codes.Internal, "panic occurred: %v", r)
		}
	}()

	resp, err = handler(ctx, request)

	duration := time.Since(start)

	var msg string
	if err != nil {
		msg = "Request completed with error"
	} else {
		msg = "Request completed successfully"
	}

	logger.Info().
		Str("method", info.FullMethod).
		Interface("request_body", request).
		Interface("response_body", resp).
		Str("start_time", start.Format(time.RFC3339)).
		Str("duration_human_readable", duration.String()).
		Int64("duration_ms", duration.Milliseconds()).
		Msg(msg)

	return resp, err
}
