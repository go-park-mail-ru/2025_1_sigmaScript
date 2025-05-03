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
}

type UserService struct {
	repo UserRepositoryInterface
}

func NewUserService(repo UserRepositoryInterface) *UserService {
	return &UserService{
		repo: repo,
	}
}
