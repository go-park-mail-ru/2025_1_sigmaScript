package service

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (s *UserService) UpdateUserAvatar(
	ctx context.Context,
	uploadDir string,
	hashedAvatarName string,
	avatarFile multipart.File,
	user models.User,
) error {
	logger := log.Ctx(ctx)
	logger.Info().Msg("creating user avatar file")

	// Create the destination file
	filePath := filepath.Join(uploadDir, hashedAvatarName)

	dst, err := os.Create(filePath)
	if err != nil {
		wrapped := errors.Wrap(err, "error creating avatar file")
		logger.Error().Err(wrapped).Msg(wrapped.Error())
		return wrapped
	}
	defer func() {
		errClose := dst.Close()
		if errClose != nil {
			wrapped := errors.Wrap(err, "error closing avatar file")
			logger.Error().Err(wrapped).Msg(wrapped.Error())
			return
		}
	}()

	// Copy the uploaded file to the destination file
	_, err = io.Copy(dst, avatarFile)
	if err != nil {
		wrapped := errors.Wrap(err, "error creating avatar file")
		logger.Error().Err(wrapped).Msg(wrapped.Error())
		return wrapped
	}

	if len(user.Avatar) > 0 && !strings.Contains(user.Avatar, "avatars/avatar_default_picture.svg") && !strings.Contains(user.Avatar, "img/avatar_placeholder.png") {
		err = os.Remove(uploadDir + getFilenameInUserAvatarField(user.Avatar))
		if err != nil {
			wrapped := errors.Wrap(err, "error removing old avatar file")
			logger.Error().Err(wrapped).Msg(wrapped.Error())
			return wrapped
		}
	}
	logger.Info().Msg("successfully created user avatar file")

	return nil
}

func getFilenameInUserAvatarField(s string) string {
	lastSlashIndex := strings.LastIndex(s, "/")
	if lastSlashIndex == -1 {
		return s
	}
	return s[lastSlashIndex+1:]
}
