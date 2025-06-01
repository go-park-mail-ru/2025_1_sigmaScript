package interfaces

import (
	"context"
	"mime/multipart"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/internal/models"
)

//go:generate mockgen -source=interface.go -destination=../mocks/mock.go
type UserServiceInterface interface {
	GetUser(ctx context.Context, login string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	Login(ctx context.Context, loginData models.LoginData) error
	DeleteUser(ctx context.Context, login string) error
	UpdateUser(ctx context.Context, login string, newUser *models.User) error
	UpdateUserAvatar(ctx context.Context, uploadDir string, hashedAvatarName string, avatarFile multipart.File, user models.User) error
	GetProfile(ctx context.Context, login string) (*models.Profile, error)
	AddFavoriteMovie(ctx context.Context, login string, movieID string) error
	AddFavoriteActor(ctx context.Context, login string, actorID string) error
	RemoveFavoriteMovie(ctx context.Context, login string, movieID string) error
	RemoveFavoriteActor(ctx context.Context, login string, actorID string) error
}
