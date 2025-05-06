package config

import (
	"context"
)

type ContextCookieKey struct{}

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
