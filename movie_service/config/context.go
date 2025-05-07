package config

import "context"

type ContextDatabaseKey struct{}

func WrapDatabaseContext(ctx context.Context, data interface{}) context.Context {
	return context.WithValue(ctx, ContextDatabaseKey{}, data)
}

func FromDatabaseContext(ctx context.Context) *Database {
	srv, ok := ctx.Value(ContextDatabaseKey{}).(*Database)
	if !ok {
		return nil
	}
	return srv
}
