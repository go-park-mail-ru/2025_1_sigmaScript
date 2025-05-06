package delivery

import (
	"context"

	auth "github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/api/auth_api_v1/proto"
)

type AuthHandlerInterface interface {
	GetSession(ctx context.Context, checkCookieReq *auth.GetSessionRequest) (*auth.GetSessionResponse, error)
	DeleteSession(ctx context.Context, deleteCookieReq *auth.DestroySessionRequest) (*auth.Nothing, error)
	CreateSession(ctx context.Context, createCookieReq *auth.CreateSessionRequest) (*auth.CreateSessionResponse, error)
}
