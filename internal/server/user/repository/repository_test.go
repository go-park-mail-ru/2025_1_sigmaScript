package repository

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewUserRepository(t *testing.T) {
	r := NewUserRepository()
	assert.NotNil(t, r)
}

func TestUserRepository_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(r *UserRepository)
		user          *models.User
		expectedError error
	}{
		{
			name:      "valid user",
			setupFunc: func(r *UserRepository) {},
			user: &models.User{
				Username:       "user",
				HashedPassword: "password",
				Avatar:         "avatar/url.png",
			},
			expectedError: nil,
		},
		{
			name: "duplicate user",
			setupFunc: func(r *UserRepository) {
				_ = r.CreateUser(context.Background(), &models.User{
					Username:       "user",
					HashedPassword: "password",
					Avatar:         "avatar/url.png",
				})
			},
			user: &models.User{
				Username:       "user",
				HashedPassword: "password",
				Avatar:         "avatar/url.png",
			},
			expectedError: errors.New(errs.ErrAlreadyExists),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := NewUserRepository()
			if tt.setupFunc != nil {
				tt.setupFunc(r)
			}

			err := r.CreateUser(context.Background(), tt.user)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRepository_DeleteUser(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(r *UserRepository)
		login         string
		expectedError error
	}{
		{
			name: "delete existing user",
			setupFunc: func(r *UserRepository) {
				_ = r.CreateUser(context.Background(), &models.User{
					Username:       "user",
					HashedPassword: "password",
					Avatar:         "avatar/url.png",
				})
			},
			login:         "user",
			expectedError: nil,
		},
		{
			name:          "delete non existing user",
			setupFunc:     func(r *UserRepository) {},
			login:         "other user",
			expectedError: errors.New(errs.ErrIncorrectLogin),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := NewUserRepository()
			if tt.setupFunc != nil {
				tt.setupFunc(r)
			}

			err := r.DeleteUser(context.Background(), tt.login)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		setupFunc     func(r *UserRepository)
		login         string
		expectedUser  *models.User
		expectedError error
	}{
		{
			name: "get existing user",
			setupFunc: func(r *UserRepository) {
				_ = r.CreateUser(ctx, &models.User{
					Username:       "user",
					HashedPassword: "password",
					Avatar:         "avatar/url.png",
				})
			},
			login: "user",
			expectedUser: &models.User{
				Username:       "user",
				HashedPassword: "password",
				Avatar:         "avatar/url.png",
			},
			expectedError: nil,
		},
		{
			name:          "get non-existent user",
			setupFunc:     func(r *UserRepository) {},
			login:         "other user",
			expectedUser:  nil,
			expectedError: errors.New(errs.ErrIncorrectLogin),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := NewUserRepository()
			if tt.setupFunc != nil {
				tt.setupFunc(r)
			}

			user, err := r.GetUser(ctx, tt.login)
			assert.Equal(t, tt.expectedUser, user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
