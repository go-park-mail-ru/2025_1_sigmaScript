package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type UserRepositoryInterface interface {
	GetUserPostgres(ctx context.Context, login string) (*models.User, error)
	CreateUserPostgres(ctx context.Context, user *models.User) error
	DeleteUserPostgres(ctx context.Context, login string) error
	UpdateUserPostgres(ctx context.Context, login string, user *models.User) (*models.User, error)
	GetUserProfilePostgres(ctx context.Context, login string) (*models.Profile, error)
	AddFavoriteMovie(ctx context.Context, login string, movieID string) error
	AddFavoriteActor(ctx context.Context, login string, actorID string) error
	RemoveFavoriteMovie(ctx context.Context, login string, movieID string) error
	RemoveFavoriteActor(ctx context.Context, login string, actorID string) error
}

type UserService struct {
	repo UserRepositoryInterface
}

func NewUserService(repo UserRepositoryInterface) *UserService {
	return &UserService{
		repo: repo,
	}
}
