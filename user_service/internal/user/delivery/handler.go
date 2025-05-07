package delivery

import (
	"context"

	user "github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/api/user_api_v1/proto"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/internal/user/delivery/adapter"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/internal/user/delivery/interfaces"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServiceGRPCHandler struct {
	user.UnimplementedUserServiceServer
	userService interfaces.UserServiceInterface
}

func NewUserServiceGRPCHandler(userService interfaces.UserServiceInterface) *UserServiceGRPCHandler {
	return &UserServiceGRPCHandler{
		userService: userService,
	}
}

func (h *UserServiceGRPCHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	logger := log.Ctx(ctx)

	usr, err := h.userService.GetUser(ctx, req.Login)
	if err != nil {
		logger.Error().Err(err).Msg("GetUser: error getting user")
		return nil, err
	}

	return adapter.ToDescGetUserRespFromSrv(usr), nil
}

func (h *UserServiceGRPCHandler) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*emptypb.Empty, error) {
	logger := log.Ctx(ctx)

	domainUser := adapter.ToSrvUserFromDesc(req)
	if err := h.userService.CreateUser(ctx, domainUser); err != nil {
		logger.Error().Err(err).Msg("CreateUser: error creating user")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *UserServiceGRPCHandler) Login(ctx context.Context, req *user.LoginRequest) (*emptypb.Empty, error) {
	logger := log.Ctx(ctx)

	loginData := adapter.ToSrvLoginDataFromDesc(req)
	if err := h.userService.Login(ctx, loginData); err != nil {
		logger.Error().Err(err).Msg("Login: error authenticating user")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *UserServiceGRPCHandler) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*emptypb.Empty, error) {
	logger := log.Ctx(ctx)

	if err := h.userService.DeleteUser(ctx, req.Login); err != nil {
		logger.Error().Err(err).Msg("DeleteUser: error deleting user")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *UserServiceGRPCHandler) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*emptypb.Empty, error) {
	logger := log.Ctx(ctx)

	newUserData := adapter.ToSrvUpdateUserFromDesc(req)

	if err := h.userService.UpdateUser(ctx, req.Login, newUserData); err != nil {
		logger.Error().Err(err).Msg("UpdateUser: error updating user")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *UserServiceGRPCHandler) UpdateUserAvatar(ctx context.Context, req *user.UpdateUserAvatarRequest) (*emptypb.Empty, error) {
	logger := log.Ctx(ctx)

	uploadDir, hashedName, avatarFile, userData, err := adapter.ToSrvAvatarDataFromDesc(req)
	if err != nil {
		logger.Error().Err(err).Msg("UpdateUserAvatar: error converting avatar data")
		return nil, err
	}

	if err := h.userService.UpdateUserAvatar(ctx, uploadDir, hashedName, avatarFile, userData); err != nil {
		logger.Error().Err(err).Msg("UpdateUserAvatar: error updating avatar")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *UserServiceGRPCHandler) GetProfile(ctx context.Context, req *user.GetProfileRequest) (*user.GetProfileResponse, error) {
	logger := log.Ctx(ctx)

	profile, err := h.userService.GetProfile(ctx, req.Login)
	if err != nil {
		logger.Error().Err(err).Msg("GetProfile: error getting profile")
		return nil, err
	}

	return adapter.ToDescGetProfileRespFromSrv(profile), nil
}

func (h *UserServiceGRPCHandler) AddFavoriteMovie(ctx context.Context, req *user.FavoriteMovieRequest) (*emptypb.Empty, error) {
	logger := log.Ctx(ctx)

	if err := h.userService.AddFavoriteMovie(ctx, req.Login, req.MovieId); err != nil {
		logger.Error().Err(err).Msg("AddFavoriteMovie: error adding favorite movie")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *UserServiceGRPCHandler) AddFavoriteActor(ctx context.Context, req *user.FavoriteActorRequest) (*emptypb.Empty, error) {
	logger := log.Ctx(ctx)

	if err := h.userService.AddFavoriteActor(ctx, req.Login, req.ActorId); err != nil {
		logger.Error().Err(err).Msg("AddFavoriteActor: error adding favorite actor")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *UserServiceGRPCHandler) RemoveFavoriteMovie(ctx context.Context, req *user.FavoriteMovieRequest) (*emptypb.Empty, error) {
	logger := log.Ctx(ctx)

	if err := h.userService.RemoveFavoriteMovie(ctx, req.Login, req.MovieId); err != nil {
		logger.Error().Err(err).Msg("RemoveFavoriteMovie: error removing favorite movie")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *UserServiceGRPCHandler) RemoveFavoriteActor(ctx context.Context, req *user.FavoriteActorRequest) (*emptypb.Empty, error) {
	logger := log.Ctx(ctx)

	if err := h.userService.RemoveFavoriteActor(ctx, req.Login, req.ActorId); err != nil {
		logger.Error().Err(err).Msg("RemoveFavoriteActor: error removing favorite actor")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
