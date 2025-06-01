package delivery

import (
	"context"
	"fmt"

	auth "github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/api/auth_api_v1/proto"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/internal/auth/delivery/adapter"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/internal/auth/delivery/interfaces"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/internal/auth/service/dto"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/internal/checker"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/internal/errors"
	"github.com/rs/zerolog/log"
)

type AuthServiceGRPCHandler struct {
	auth.UnimplementedSessionRPCServer
	sessionService interfaces.SessionServiceInterface
}

func NewAuthServiceGRPCHandler(sessionService interfaces.SessionServiceInterface) *AuthServiceGRPCHandler {
	return &AuthServiceGRPCHandler{
		sessionService: sessionService,
	}
}

// CreateSession create_session grpc handler
func (am *AuthServiceGRPCHandler) CreateSession(ctx context.Context, createCookieReq *auth.CreateSessionRequest) (*auth.CreateSessionResponse, error) {
	logger := log.Ctx(ctx)
	err := checker.ValidateUserID(createCookieReq.UserID)
	if err != nil {
		err = fmt.Errorf("create_session: %w", err)
		logger.Error().Err(err).Msg("failed_to_create_session")

		return nil, err
	}

	srvData := adapter.ToSrvCreateCookieFromDesc(createCookieReq)
	if srvData == nil {
		return nil, errs.ErrBadRequest
	}

	srvResp, err := am.sessionService.CreateSession(ctx, string(srvData.UserID))
	resp := adapter.ToDescCreateCookieRespFromSrv(&dto.Cookie{
		Name:   "session_id",
		UserID: srvResp,
	})
	if err != nil {
		logger.Error().Interface("createSessionError", err).Msg("failed_to_create_session")
		return nil, err
	}

	logger.Info().Interface("createSessionSuccess", resp).Msg("successfully_create_session")
	return resp, nil
}

// DestroySession destroy_session grpc handler
func (am *AuthServiceGRPCHandler) DeleteSession(ctx context.Context, deleteCookieReq *auth.DestroySessionRequest) (*auth.Nothing, error) {
	logger := log.Ctx(ctx)
	err := checker.ValidateCookie(deleteCookieReq.Cookie)
	if err != nil {
		err = fmt.Errorf("destroy_session: %w", err)
		logger.Error().Err(err).Msg("failed_to_destroy_session")

		return nil, err
	}

	svcResp, err := am.sessionService.DeleteSession(ctx, deleteCookieReq.Cookie)
	if err != nil {
		logger.Error().Interface("destroySessionError", err).Msg("failed_to_destroy_session")
		return nil, err
	}

	logger.Info().Msg(svcResp)
	return &auth.Nothing{Dummy: true}, nil
}

// Session get_session_data grpc handler
func (am *AuthServiceGRPCHandler) GetSession(ctx context.Context, checkCookieReq *auth.GetSessionRequest) (*auth.GetSessionResponse, error) {
	logger := log.Ctx(ctx)
	err := checker.ValidateCookie(checkCookieReq.Cookie)
	if err != nil {
		err = fmt.Errorf("check_session: %w", err)
		logger.Error().Err(err).Msg("failed_to_get_session_data")

		return nil, err
	}

	srvResp, err := am.sessionService.GetSession(ctx, checkCookieReq.Cookie)
	if err != nil {
		logger.Error().Interface("getSessionDataError", err).Msg("failed_to_get_session_data")
		return nil, err
	}

	logger.Info().Interface("getSessionDataSuccess", srvResp).Msg("successfully_get_session_data")
	return &auth.GetSessionResponse{UserID: srvResp}, nil
}
