package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type UserRepositoryInterface interface {
	GetUser(ctx context.Context, login string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, login string) error
}

type UserService struct {
	repo UserRepositoryInterface
}

func NewUserService(repo UserRepositoryInterface) *UserService {
	return &UserService{
		repo: repo,
	}
}
