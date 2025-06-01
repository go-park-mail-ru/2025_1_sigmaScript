package delivery

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	user "github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/api/user_api_v1/proto"
)

type UserHandlerInterface interface {
	GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error)
	CreateUser(ctx context.Context, req *user.CreateUserRequest) (*emptypb.Empty, error)
	Login(ctx context.Context, req *user.LoginRequest) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*emptypb.Empty, error)
	UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*emptypb.Empty, error)
	UpdateUserAvatar(ctx context.Context, req *user.UpdateUserAvatarRequest) (*emptypb.Empty, error)
	GetProfile(ctx context.Context, req *user.GetProfileRequest) (*user.GetProfileResponse, error)
	AddFavoriteMovie(ctx context.Context, req *user.FavoriteMovieRequest) (*emptypb.Empty, error)
	AddFavoriteActor(ctx context.Context, req *user.FavoriteActorRequest) (*emptypb.Empty, error)
	RemoveFavoriteMovie(ctx context.Context, req *user.FavoriteMovieRequest) (*emptypb.Empty, error)
	RemoveFavoriteActor(ctx context.Context, req *user.FavoriteActorRequest) (*emptypb.Empty, error)
}
