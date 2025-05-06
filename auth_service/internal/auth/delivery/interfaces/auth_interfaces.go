package interfaces

import (
	"context"
)

//go:generate mockgen -source=auth_interfaces.go -destination=../mocks/mock.go
type SessionServiceInterface interface {
	GetSession(ctx context.Context, sessionID string) (string, error)
	DeleteSession(ctx context.Context, sessionID string) (string, error)
	CreateSession(ctx context.Context, username string) (string, error)
}
