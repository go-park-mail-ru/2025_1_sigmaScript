package adapter

import (
	"bytes"
	"io"
	"mime/multipart"

	user "github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/api/user_api_v1/proto"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/internal/models"
)

type nopCloserReader struct {
	*bytes.Reader
}

func (n nopCloserReader) Close() error {
	return nil
}

func ToDescGetUserRespFromSrv(u *models.User) *user.GetUserResponse {
	return &user.GetUserResponse{
		User: convertUser(u),
	}
}

func ToSrvUserFromDesc(req *user.CreateUserRequest) *models.User {
	if req == nil || req.User == nil {
		return nil
	}

	return &models.User{
		ID:             req.User.Id,
		Username:       req.User.Username,
		HashedPassword: req.User.HashedPassword,
		Avatar:         req.User.Avatar,
		CreatedAt:      req.User.CreatedAt,
		UpdatedAt:      req.User.UpdatedAt,
	}
}

func ToSrvLoginDataFromDesc(req *user.LoginRequest) models.LoginData {
	var data models.LoginData
	if req != nil && req.LoginData != nil {
		data = models.LoginData{
			Username: req.LoginData.Username,
			Password: req.LoginData.Password,
		}
	}

	return data
}

func ToSrvAvatarDataFromDesc(req *user.UpdateUserAvatarRequest) (string, string, multipart.File, models.User, error) {
	if req == nil || req.User == nil {
		return "", "", nil, models.User{}, io.EOF
	}

	avatarFile := nopCloserReader{bytes.NewReader(req.AvatarFile)}

	userData := models.User{
		ID:             req.User.Id,
		Username:       req.User.Username,
		HashedPassword: req.User.HashedPassword,
		Avatar:         req.User.Avatar,
		CreatedAt:      req.User.CreatedAt,
		UpdatedAt:      req.User.UpdatedAt,
	}

	return req.UploadDir, req.HashedAvatarName, avatarFile, userData, nil
}

func ToDescGetProfileRespFromSrv(p *models.Profile) *user.GetProfileResponse {
	if p == nil {
		return &user.GetProfileResponse{}
	}

	pbProf := &user.Profile{
		Username:  p.Username,
		Avatar:    p.Avatar,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	for _, m := range p.MovieCollection {
		pbProf.MovieCollection = append(pbProf.MovieCollection, &user.Movie{
			Id:          int32(m.ID),
			Title:       m.Title,
			PreviewUrl:  m.PreviewURL,
			Duration:    m.Duration,
			ReleaseDate: m.ReleaseDate,
			Rating:      m.Rating,
		})
	}

	for _, a := range p.Actors {
		pbProf.Actors = append(pbProf.Actors, &user.Person{
			Id:         int32(a.ID),
			FullName:   a.FullName,
			EnFullName: a.EnFullName,
			Photo:      a.Photo,
			About:      a.About,
			Sex:        a.Sex,
			Growth:     a.Growth,
			Birthday:   a.Birthday,
			Death:      a.Death,
			Career:     a.Career,
			Genres:     a.Genres,
			TotalFilms: a.TotalFilms,
		})
	}

	for _, r := range p.Reviews {
		pbProf.Reviews = append(pbProf.Reviews, &user.Review{
			Id:         int32(r.ID),
			Score:      r.Score,
			ReviewText: r.ReviewText,
			CreatedAt:  r.CreatedAt,
			User: &user.ReviewUserData{
				Login:  r.User.Login,
				Avatar: r.User.Avatar,
			},
		})
	}

	return &user.GetProfileResponse{
		Profile: pbProf,
	}
}

func convertUser(u *models.User) *user.User {
	if u == nil {
		return nil
	}

	return &user.User{
		Id:             u.ID,
		Username:       u.Username,
		HashedPassword: u.HashedPassword,
		Avatar:         u.Avatar,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

func ToSrvUpdateUserFromDesc(u *user.UpdateUserRequest) *models.User {
	if u == nil {
		return nil
	}

	return &models.User{
		ID:             u.NewUser.Id,
		Username:       u.NewUser.Username,
		HashedPassword: u.NewUser.HashedPassword,
		Avatar:         u.NewUser.Avatar,
		CreatedAt:      u.NewUser.CreatedAt,
		UpdatedAt:      u.NewUser.UpdatedAt,
	}
}
