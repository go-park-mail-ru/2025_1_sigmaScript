package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mockSessionRepo "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/service/mocks"
)

func dummyCtxWithSessionLength(length int) context.Context {
	cfg := &config.Cookie{
		SessionName:   "session_id",
		HTTPOnly:      true,
		Secure:        false,
		SameSite:      http.SameSiteLaxMode,
		Path:          "/",
		SessionLength: length,
	}
	return config.WrapCookieContext(context.Background(), cfg)
}

func TestNewSessionService(t *testing.T) {
	ctx := dummyCtxWithSessionLength(64)
	mockRepo := mockSessionRepo.NewMockSessionRepositoryInterface(gomock.NewController(t))
	svc := NewSessionService(ctx, mockRepo)
	assert.NotNil(t, svc)
}

func TestSessionService_CreateSession(t *testing.T) {
	ctx := dummyCtxWithSessionLength(64)

	tests := []struct {
		name          string
		username      string
		mockSetupFunc func(t *testing.T, repo *mockSessionRepo.MockSessionRepositoryInterface)
		expectedError error
	}{
		{
			name:     "success create session",
			username: "user",
			mockSetupFunc: func(t *testing.T, repo *mockSessionRepo.MockSessionRepositoryInterface) {
				repo.EXPECT().StoreSession(gomock.Any(), gomock.Any(), "user").Return(nil).Times(1)
			},
			expectedError: nil,
		},
		{
			name:     "error generate session",
			username: "user",
			mockSetupFunc: func(t *testing.T, repo *mockSessionRepo.MockSessionRepositoryInterface) {
				repo.EXPECT().StoreSession(gomock.Any(), gomock.Any(), "user").
					Return(errs.ErrGenerateSession).Times(1)
			},
			expectedError: errs.ErrGenerateSession,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockSessionRepo.NewMockSessionRepositoryInterface(ctrl)
			if tt.mockSetupFunc != nil {
				tt.mockSetupFunc(t, mockRepo)
			}

			svc := NewSessionService(ctx, mockRepo)
			_, err := svc.CreateSession(context.Background(), tt.username)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSessionService_DeleteSession(t *testing.T) {
	ctx := dummyCtxWithSessionLength(64)

	tests := []struct {
		name          string
		sessionID     string
		mockSetupFunc func(t *testing.T, repo *mockSessionRepo.MockSessionRepositoryInterface)
		expectedError error
	}{
		{
			name:      "delete existing session",
			sessionID: "session",
			mockSetupFunc: func(t *testing.T, repo *mockSessionRepo.MockSessionRepositoryInterface) {
				repo.EXPECT().DeleteSession(gomock.Any(), "session").Return(nil).Times(1)
			},
			expectedError: nil,
		},
		{
			name:      "delete non existing session",
			sessionID: "session",
			mockSetupFunc: func(t *testing.T, repo *mockSessionRepo.MockSessionRepositoryInterface) {
				repo.EXPECT().DeleteSession(gomock.Any(), "session").Return(errs.ErrSessionNotExists).Times(1)
			},
			expectedError: errs.ErrSessionNotExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockSessionRepo.NewMockSessionRepositoryInterface(ctrl)
			if tt.mockSetupFunc != nil {
				tt.mockSetupFunc(t, mockRepo)
			}

			svc := NewSessionService(ctx, mockRepo)
			err := svc.DeleteSession(context.Background(), tt.sessionID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSessionService_GetSession(t *testing.T) {
	ctx := dummyCtxWithSessionLength(64)

	tests := []struct {
		name          string
		sessionID     string
		mockSetupFunc func(t *testing.T, repo *mockSessionRepo.MockSessionRepositoryInterface)
		expectedError error
		expectedLogin string
	}{
		{
			name:      "get existing session",
			sessionID: "session",
			mockSetupFunc: func(t *testing.T, repo *mockSessionRepo.MockSessionRepositoryInterface) {
				repo.EXPECT().GetSession(gomock.Any(), "session").
					Return("user", nil).Times(1)
			},
			expectedError: nil,
			expectedLogin: "user",
		},
		{
			name:      "get non existing session",
			sessionID: "session",
			mockSetupFunc: func(t *testing.T, repo *mockSessionRepo.MockSessionRepositoryInterface) {
				repo.EXPECT().GetSession(gomock.Any(), "session").
					Return(noData, errs.ErrSessionNotExists).Times(1)
			},
			expectedError: errs.ErrSessionNotExists,
			expectedLogin: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockSessionRepo.NewMockSessionRepositoryInterface(ctrl)
			if tt.mockSetupFunc != nil {
				tt.mockSetupFunc(t, mockRepo)
			}

			svc := NewSessionService(ctx, mockRepo)
			login, err := svc.GetSession(context.Background(), tt.sessionID)
			assert.Equal(t, tt.expectedLogin, login)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
