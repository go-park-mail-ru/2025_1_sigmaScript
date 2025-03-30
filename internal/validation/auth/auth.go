package auth

import (
	"strings"
	"unicode/utf8"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/pkg/errors"
)

const (
	MinPasswordLength = 6
	MaxPasswordLength = 18
)

func IsValidPassword(password string) error {
	if utf8.RuneCountInString(password) < MinPasswordLength {
		return errors.New(errs.ErrPasswordTooShort)
	}
	if strings.TrimSpace(password) == "" {
		return errors.New(errs.ErrEmptyPassword)
	}
	if utf8.RuneCountInString(password) > MaxPasswordLength {
		return errors.New(errs.ErrPasswordTooLong)
	}
	return nil
}

func IsValidLogin(login string) error {

	if login = strings.TrimSpace(login); login == "" {
		return errors.New(errs.ErrEmptyLogin)
	}

	cnt := utf8.RuneCountInString(login)
	if cnt < 2 || cnt > 18 {
		return errors.New(errs.ErrLengthLogin)
	}

	allowedLocalChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
	for _, char := range login {
		if !strings.ContainsRune(allowedLocalChars, char) {
			return errors.New(errs.ErrInvalidLogin)
		}
	}
	return nil
}
