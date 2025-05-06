package client

import (
	"context"

	auth "github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/api/auth_api_v1/proto"
)

// AuthClientInterface defines client methods to transmit to Auth microservice
//
//go:generate mockgen -source=auth.go -destination=../auth/service/mocks/mock.go
type AuthClientInterface interface {
	GetSession(ctx context.Context, sessionID string) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error
	CreateSession(ctx context.Context, username string) (string, error)
}

// AuthClient struct implements AuthClientInterface
type AuthClient struct {
	authMS auth.SessionRPCClient
}

// NewAuthClient returns an instance of AuthClientInterface
func NewAuthClient(authMS auth.SessionRPCClient) AuthClientInterface {
	return &AuthClient{
		authMS: authMS,
	}
}

// CreateSession creates new session
func (cl *AuthClient) CreateSession(ctx context.Context, username string) (string, error) {

	resp, err := cl.authMS.CreateSession(ctx, &auth.CreateSessionRequest{UserID: username})

	if err != nil {
		return "", err
	}

	return resp.Cookie, nil
}

// DestroySession destroys session
func (cl *AuthClient) DeleteSession(ctx context.Context, cookie string) error {
	_, err := cl.authMS.DeleteSession(ctx, &auth.DestroySessionRequest{Cookie: cookie})

	if err != nil {
		return err
	}

	return nil
}

// Session checks active session
func (cl *AuthClient) GetSession(ctx context.Context, cookie string) (string, error) {

	resp, err := cl.authMS.GetSession(ctx, &auth.GetSessionRequest{Cookie: cookie})

	if err != nil {
		return "", err
	}

	return resp.UserID, nil
}
