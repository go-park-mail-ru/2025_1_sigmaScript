package repository

import (
	"context"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewSessionRepository(t *testing.T) {
	r := NewSessionRepository()
	assert.NotNil(t, r)
}

func TestSessionRepository_GetSession(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(r *SessionRepository)
		sessionID     string
		login         string
		expectedError error
		expectedLogin string
	}{
		{
			name: "get existing session",
			setupFunc: func(r *SessionRepository) {
				err := r.StoreSession(context.Background(), "session", "user")
				assert.Nil(t, err)
			},
			sessionID:     "session",
			login:         "user",
			expectedError: nil,
			expectedLogin: "user",
		},
		{
			name:          "get non existing session",
			setupFunc:     func(r *SessionRepository) {},
			sessionID:     "session",
			login:         "user",
			expectedError: errs.ErrSessionNotExists,
			expectedLogin: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := NewSessionRepository()
			if tt.setupFunc != nil {
				tt.setupFunc(r)
			}

			login, err := r.GetSession(context.Background(), tt.sessionID)
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

func TestSessionRepository_DeleteSession(t *testing.T) {
	tests := []struct {
		name          string
		sessionID     string
		setupFunc     func(r *SessionRepository)
		expectedError error
	}{
		{
			name:      "delete existing session",
			sessionID: "session",
			setupFunc: func(r *SessionRepository) {
				err := r.StoreSession(context.Background(), "session", "user")
				assert.Nil(t, err)
			},
			expectedError: nil,
		},
		{
			name:          "delete non existing session",
			sessionID:     "other session",
			setupFunc:     func(r *SessionRepository) {},
			expectedError: errs.ErrSessionNotExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := NewSessionRepository()
			if tt.setupFunc != nil {
				tt.setupFunc(r)
			}

			err := r.DeleteSession(context.Background(), tt.sessionID)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSessionRepository_StoreSession(t *testing.T) {
	tests := []struct {
		name      string
		sessionID string
		login     string
	}{
		{
			name:      "store valid session",
			sessionID: "session",
			login:     "user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := NewSessionRepository()
			err := r.StoreSession(context.Background(), tt.sessionID, tt.login)
			assert.NoError(t, err)
		})
	}
}
