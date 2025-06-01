package service

import (
	"context"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	mockRepo "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mockRepo.NewMockUserRepositoryInterface(ctrl)
	s := NewUserService(r)

	assert.NotNil(t, s)
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		user          *models.User
		mockSetupFunc func(*testing.T, *mockRepo.MockUserRepositoryInterface)
		expectedError error
	}{
		{
			name: "valid user",
			user: &models.User{
				Username:       "test",
				HashedPassword: "test",
				Avatar:         "test/url.png",
			},
			mockSetupFunc: func(t *testing.T, r *mockRepo.MockUserRepositoryInterface) {
				r.EXPECT().CreateUserPostgres(gomock.Any(), &models.User{
					Username:       "test",
					HashedPassword: "test",
					Avatar:         "test/url.png",
				}).
					Return(nil).Times(1)
			},
			expectedError: nil,
		},
		{
			name: "already exist user",
			user: &models.User{
				Username:       "test",
				HashedPassword: "test",
				Avatar:         "test/url.png",
			},
			mockSetupFunc: func(t *testing.T, r *mockRepo.MockUserRepositoryInterface) {
				r.EXPECT().CreateUserPostgres(gomock.Any(), &models.User{
					Username:       "test",
					HashedPassword: "test",
					Avatar:         "test/url.png",
				}).
					Return(errors.New(errs.ErrAlreadyExists)).Times(1)
			},
			expectedError: errors.New(errs.ErrAlreadyExists),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mockRepo.NewMockUserRepositoryInterface(ctrl)
			if tt.mockSetupFunc != nil {
				tt.mockSetupFunc(t, r)
			}

			s := NewUserService(r)
			err := s.CreateUser(context.Background(), tt.user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		mockSetupFunc func(*testing.T, *mockRepo.MockUserRepositoryInterface)
		expectedError error
	}{
		{
			name:     "valid user",
			username: "test",
			mockSetupFunc: func(t *testing.T, r *mockRepo.MockUserRepositoryInterface) {
				r.EXPECT().DeleteUserPostgres(gomock.Any(), "test").
					Return(nil).Times(1)
			},
			expectedError: nil,
		},
		{
			name:     "not exist user",
			username: "test",
			mockSetupFunc: func(t *testing.T, r *mockRepo.MockUserRepositoryInterface) {
				r.EXPECT().DeleteUserPostgres(gomock.Any(), "test").
					Return(errors.New(errs.ErrIncorrectLogin)).Times(1)
			},
			expectedError: errors.New(errs.ErrIncorrectLogin),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mockRepo.NewMockUserRepositoryInterface(ctrl)
			if tt.mockSetupFunc != nil {
				tt.mockSetupFunc(t, r)
			}

			s := NewUserService(r)
			err := s.DeleteUser(context.Background(), tt.username)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		mockSetupFunc func(*testing.T, *mockRepo.MockUserRepositoryInterface)
		expectedError error
		expectedUser  *models.User
	}{
		{
			name:     "valid user",
			username: "valid user",
			mockSetupFunc: func(t *testing.T, r *mockRepo.MockUserRepositoryInterface) {
				r.EXPECT().GetUserPostgres(gomock.Any(), "valid user").
					Return(&models.User{
						Username:       "test",
						HashedPassword: "test",
						Avatar:         "test/url.png",
					}, nil).Times(1)
			},
			expectedError: nil,
			expectedUser: &models.User{
				Username:       "test",
				HashedPassword: "test",
				Avatar:         "test/url.png",
			},
		},
		{
			name:     "invalid user",
			username: "invalid user",
			mockSetupFunc: func(t *testing.T, r *mockRepo.MockUserRepositoryInterface) {
				r.EXPECT().GetUserPostgres(gomock.Any(), "invalid user").
					Return(nil, errors.New(errs.ErrIncorrectLogin)).Times(1)
			},
			expectedError: errors.New(errs.ErrIncorrectLogin),
			expectedUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mockRepo.NewMockUserRepositoryInterface(ctrl)
			if tt.mockSetupFunc != nil {
				tt.mockSetupFunc(t, r)
			}

			s := NewUserService(r)
			user, err := s.GetUser(context.Background(), tt.username)

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

func TestUserService_Login(t *testing.T) {
	tests := []struct {
		name          string
		loginData     models.LoginData
		mockSetupFunc func(*testing.T, *mockRepo.MockUserRepositoryInterface)
		expectedError error
	}{
		{
			name: "valid user",
			loginData: models.LoginData{
				Username: "valid user",
				Password: "test password",
			},
			mockSetupFunc: func(t *testing.T, r *mockRepo.MockUserRepositoryInterface) {
				hashedPass, err := bcrypt.GenerateFromPassword([]byte("test password"), bcrypt.DefaultCost)
				assert.NoError(t, err)
				r.EXPECT().GetUserPostgres(gomock.Any(), "valid user").
					Return(&models.User{
						Username:       "valid user",
						HashedPassword: string(hashedPass),
						Avatar:         "test/url.png",
					}, nil).Times(1)
			},
			expectedError: nil,
		},
		{
			name: "invalid user",
			loginData: models.LoginData{
				Username: "invalid user",
				Password: "test password",
			},
			mockSetupFunc: func(t *testing.T, r *mockRepo.MockUserRepositoryInterface) {
				r.EXPECT().GetUserPostgres(gomock.Any(), "invalid user").
					Return(nil, errors.New(errs.ErrIncorrectLogin)).Times(1)
			},
			expectedError: errors.New(errs.ErrIncorrectLogin),
		},
		{
			name: "invalid user",
			loginData: models.LoginData{
				Username: "invalid user",
				Password: "incorrect password",
			},
			mockSetupFunc: func(t *testing.T, r *mockRepo.MockUserRepositoryInterface) {
				r.EXPECT().GetUserPostgres(gomock.Any(), "invalid user").
					Return(&models.User{
						Username:       "invalid user",
						HashedPassword: "correct password",
						Avatar:         "test/url.png",
					}, nil).Times(1)
			},
			expectedError: errors.New(errs.ErrIncorrectPassword),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mockRepo.NewMockUserRepositoryInterface(ctrl)
			if tt.mockSetupFunc != nil {
				tt.mockSetupFunc(t, r)
			}

			s := NewUserService(r)
			err := s.Login(context.Background(), tt.loginData)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		newUser       *models.User
		mockSetupFunc func(*testing.T, *mockRepo.MockUserRepositoryInterface)
		expectedError error
	}{
		{
			name:  "valid user",
			login: "valid user",
			newUser: &models.User{
				Username:       "valid user",
				HashedPassword: "test password",
				Avatar:         "test/url.png",
			},
			mockSetupFunc: func(t *testing.T, r *mockRepo.MockUserRepositoryInterface) {
				r.EXPECT().UpdateUserPostgres(gomock.Any(), "valid user", &models.User{
					Username:       "valid user",
					HashedPassword: "test password",
					Avatar:         "test/url.png",
				}).
					Return(&models.User{
						Username:       "valid user",
						HashedPassword: "test password",
						Avatar:         "test/url.png",
					}, nil)
			},
			expectedError: nil,
		},
		{
			name:  "invalid user",
			login: "incorrect user",
			newUser: &models.User{
				Username:       "incorrect user",
				HashedPassword: "test password",
				Avatar:         "test/url.png",
			},
			mockSetupFunc: func(t *testing.T, r *mockRepo.MockUserRepositoryInterface) {
				r.EXPECT().UpdateUserPostgres(gomock.Any(), "incorrect user", &models.User{
					Username:       "incorrect user",
					HashedPassword: "test password",
					Avatar:         "test/url.png",
				}).
					Return(nil, errors.New(errs.ErrIncorrectLogin))
			},
			expectedError: errors.New(errs.ErrIncorrectLogin),
		},
		// {
		// 	name:  "invalid user",
		// 	login: "correct login",
		// 	newUser: &models.User{
		// 		Username:       "correct login",
		// 		HashedPassword: "test password",
		// 		Avatar:         "test/url.png",
		// 	},
		// 	mockSetupFunc: func(t *testing.T, r *mockRepo.MockUserRepositoryInterface) {
		// 		r.EXPECT().DeleteUserPostgres(gomock.Any(), "correct login").
		// 			Return(nil).Times(1)
		// 		r.EXPECT().CreateUserPostgres(gomock.Any(), &models.User{
		// 			Username:       "correct login",
		// 			HashedPassword: "test password",
		// 			Avatar:         "test/url.png",
		// 		}).Return(errors.New(errs.ErrAlreadyExists)).Times(1)
		// 	},
		// 	expectedError: errors.New(errs.ErrAlreadyExists),
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mockRepo.NewMockUserRepositoryInterface(ctrl)
			if tt.mockSetupFunc != nil {
				tt.mockSetupFunc(t, r)
			}

			s := NewUserService(r)
			err := s.UpdateUser(context.Background(), tt.login, tt.newUser)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
