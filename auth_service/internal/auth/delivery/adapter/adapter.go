package adapter

import (
	auth "github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/api/auth_api_v1/proto"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/internal/auth/service/dto"
)

// ToSrvCreateCookieFromDesc converts grpc createSessionReq to srv dto
func ToSrvCreateCookieFromDesc(createCookieReq *auth.CreateSessionRequest) *dto.SrvCreateCookie {
	if createCookieReq == nil {
		return nil
	}

	return &dto.SrvCreateCookie{
		UserID: createCookieReq.UserID,
	}
}

// ToDescCreateCookieRespFromSrv converts srv dto createSessionResp to grpc response
func ToDescCreateCookieRespFromSrv(cookie *dto.Cookie) *auth.CreateSessionResponse {
	if cookie == nil {
		return nil
	}

	return &auth.CreateSessionResponse{
		Name:   cookie.Name,
		Cookie: cookie.UserID,
		MaxAge: cookie.Expiry.Unix(),
	}
}
