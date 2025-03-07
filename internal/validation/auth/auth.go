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
