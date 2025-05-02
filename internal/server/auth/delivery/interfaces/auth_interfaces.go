package interfaces

import (
	"context"
	"mime/multipart"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
)

//go:generate mockgen -source=auth_interfaces.go -destination=../mocks/mock.go
type UserServiceInterface interface {
	GetUser(ctx context.Context, login string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	Login(ctx context.Context, loginData models.LoginData) error
	DeleteUser(ctx context.Context, login string) error
	UpdateUser(ctx context.Context, login string, newUser *models.User) error
	UpdateUserAvatar(ctx context.Context, uploadDir string, hashedAvatarName string, avatarFile multipart.File, user models.User) error
}

//go:generate mockgen -source=auth_interfaces.go -destination=../mocks/mock.go
type SessionServiceInterface interface {
	GetSession(ctx context.Context, sessionID string) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error
	CreateSession(ctx context.Context, username string) (string, error)
}
