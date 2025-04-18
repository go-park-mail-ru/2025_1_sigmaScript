package config

import (
	"context"
)

type ContextServerKey struct{}
type ContextCookieKey struct{}

func WrapServerContext(ctx context.Context, data interface{}) context.Context {
	return context.WithValue(ctx, ContextServerKey{}, data)
}

func FromServerContext(ctx context.Context) *Server {
	srv, ok := ctx.Value(ContextServerKey{}).(*Server)
	if !ok {
		return nil
	}
	return srv
}

func WrapCookieContext(ctx context.Context, data interface{}) context.Context {
	return context.WithValue(ctx, ContextCookieKey{}, data)
}

func FromCookieContext(ctx context.Context) *Cookie {
	cookie, ok := ctx.Value(ContextCookieKey{}).(*Cookie)
	if !ok {
		return nil
	}
	return cookie
}
