package client

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	userpb "github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/api/user_api_v1/proto"
)

//go:generate mockgen -source=auth.go -destination=../auth/service/mocks/mock.go
type UserClientInterface interface {
	GetUser(ctx context.Context, login string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	Login(ctx context.Context, loginData models.LoginData) error
	DeleteUser(ctx context.Context, login string) error
	UpdateUser(ctx context.Context, login string, newUser *models.User) error
	UpdateUserAvatar(ctx context.Context, uploadDir, hashedAvatarName string, avatarFile multipart.File, user models.User) error
	GetProfile(ctx context.Context, login string) (*models.Profile, error)
	AddFavoriteMovie(ctx context.Context, login, movieID string) error
	AddFavoriteActor(ctx context.Context, login, actorID string) error
	RemoveFavoriteMovie(ctx context.Context, login, movieID string) error
	RemoveFavoriteActor(ctx context.Context, login, actorID string) error
}

type UserClient struct {
	userMicroService userpb.UserServiceClient
}

func NewUserClient(userMS userpb.UserServiceClient) *UserClient {
	return &UserClient{
		userMicroService: userMS,
	}
}

func (cl *UserClient) GetUser(ctx context.Context, login string) (*models.User, error) {
	resp, err := cl.userMicroService.GetUser(ctx, &userpb.GetUserRequest{Login: login})
	if err != nil {
		return nil, err
	}

	pb := resp.GetUser()
	return &models.User{
		ID:             pb.GetId(),
		Username:       pb.GetUsername(),
		Avatar:         pb.GetAvatar(),
		CreatedAt:      pb.GetCreatedAt(),
		UpdatedAt:      pb.GetUpdatedAt(),
		HashedPassword: pb.GetHashedPassword(),
	}, nil
}

func (cl *UserClient) CreateUser(ctx context.Context, u *models.User) error {
	_, err := cl.userMicroService.CreateUser(ctx, &userpb.CreateUserRequest{
		User: &userpb.User{
			Id:             u.ID,
			Username:       u.Username,
			Avatar:         u.Avatar,
			CreatedAt:      u.CreatedAt,
			UpdatedAt:      u.UpdatedAt,
			HashedPassword: u.HashedPassword,
		},
	})

	return err
}

func (cl *UserClient) Login(ctx context.Context, loginData models.LoginData) error {
	_, err := cl.userMicroService.Login(ctx, &userpb.LoginRequest{
		LoginData: &userpb.LoginData{
			Username: loginData.Username,
			Password: loginData.Password,
		},
	})

	return err
}

func (cl *UserClient) DeleteUser(ctx context.Context, login string) error {
	_, err := cl.userMicroService.DeleteUser(ctx, &userpb.DeleteUserRequest{Login: login})

	return err
}

func (cl *UserClient) UpdateUser(ctx context.Context, login string, u *models.User) error {
	_, err := cl.userMicroService.UpdateUser(ctx, &userpb.UpdateUserRequest{
		Login: login,
		NewUser: &userpb.User{
			Id:             u.ID,
			Username:       u.Username,
			Avatar:         u.Avatar,
			CreatedAt:      u.CreatedAt,
			UpdatedAt:      u.UpdatedAt,
			HashedPassword: u.HashedPassword,
		},
	})

	return err
}

func (cl *UserClient) UpdateUserAvatar(ctx context.Context, uploadDir, hashedAvatarName string, avatarFile multipart.File, user models.User) error {
	data, err := io.ReadAll(avatarFile)
	if err != nil {
		return fmt.Errorf("failed to read avatar file: %w", err)
	}

	_, err = cl.userMicroService.UpdateUserAvatar(ctx, &userpb.UpdateUserAvatarRequest{
		UploadDir:        uploadDir,
		HashedAvatarName: hashedAvatarName,
		AvatarFile:       data,
		User: &userpb.User{
			Id:             user.ID,
			Username:       user.Username,
			Avatar:         user.Avatar,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			HashedPassword: user.HashedPassword,
		},
	})

	return err
}

func (cl *UserClient) GetProfile(ctx context.Context, login string) (*models.Profile, error) {
	resp, err := cl.userMicroService.GetProfile(ctx, &userpb.GetProfileRequest{Login: login})
	if err != nil {
		return nil, err
	}
	pb := resp.GetProfile()

	var movies []mocks.Movie
	for _, m := range pb.GetMovieCollection() {
		movies = append(movies, mocks.Movie{
			ID:          int(m.GetId()),
			Title:       m.GetTitle(),
			PreviewURL:  m.GetPreviewUrl(),
			Duration:    m.GetDuration(),
			ReleaseDate: m.GetReleaseDate(),
			Rating:      m.GetRating(),
		})
	}

	var actors []mocks.PersonJSON
	for _, a := range pb.GetActors() {
		actors = append(actors, mocks.PersonJSON{
			ID:         int(a.GetId()),
			FullName:   a.GetFullName(),
			EnFullName: a.GetEnFullName(),
			Photo:      a.GetPhoto(),
			About:      a.GetAbout(),
			Sex:        a.GetSex(),
			Growth:     a.GetGrowth(),
			Birthday:   a.GetBirthday(),
			Death:      a.GetDeath(),
			Career:     a.GetCareer(),
			Genres:     a.GetGenres(),
			TotalFilms: a.GetTotalFilms(),
		})
	}

	var reviews []mocks.ReviewJSON
	for _, r := range pb.GetReviews() {
		reviews = append(reviews, mocks.ReviewJSON{
			ID:         int(r.GetId()),
			Score:      r.GetScore(),
			ReviewText: r.GetReviewText(),
			CreatedAt:  r.GetCreatedAt(),
			User: mocks.ReviewUserDataJSON{
				Login:  r.GetUser().GetLogin(),
				Avatar: r.GetUser().GetAvatar(),
			},
		})
	}

	return &models.Profile{
		Username:        pb.GetUsername(),
		Avatar:          pb.GetAvatar(),
		CreatedAt:       pb.GetCreatedAt(),
		UpdatedAt:       pb.GetUpdatedAt(),
		MovieCollection: movies,
		Actors:          actors,
		Reviews:         reviews,
	}, nil
}

func (cl *UserClient) AddFavoriteMovie(ctx context.Context, login, movieID string) error {
	_, err := cl.userMicroService.AddFavoriteMovie(ctx, &userpb.FavoriteMovieRequest{
		Login:   login,
		MovieId: movieID,
	})

	return err
}

func (cl *UserClient) AddFavoriteActor(ctx context.Context, login, actorID string) error {
	_, err := cl.userMicroService.AddFavoriteActor(ctx, &userpb.FavoriteActorRequest{
		Login:   login,
		ActorId: actorID,
	})

	return err
}

func (cl *UserClient) RemoveFavoriteMovie(ctx context.Context, login, movieID string) error {
	_, err := cl.userMicroService.RemoveFavoriteMovie(ctx, &userpb.FavoriteMovieRequest{
		Login:   login,
		MovieId: movieID,
	})

	return err
}

func (cl *UserClient) RemoveFavoriteActor(ctx context.Context, login, actorID string) error {
	_, err := cl.userMicroService.RemoveFavoriteActor(ctx, &userpb.FavoriteActorRequest{
		Login:   login,
		ActorId: actorID,
	})

	return err
}
