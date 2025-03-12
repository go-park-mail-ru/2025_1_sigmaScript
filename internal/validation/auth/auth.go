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

func IsValidEmail(email string) error {
  atIndex := strings.Index(email, "@")
  if atIndex == -1 || atIndex != strings.LastIndex(email, "@") {
    return errors.New(errs.ErrInvalidEmail)
  }

  local := email[:atIndex]
  domain := email[atIndex+1:]

  if len(local) == 0 || len(domain) == 0 {
    return errors.New(errs.ErrInvalidEmail)
  }

  if !strings.Contains(domain, ".") {
    return errors.New(errs.ErrInvalidEmail)
  }

  if domain[0] == '.' || domain[len(domain)-1] == '.' {
    return errors.New(errs.ErrInvalidEmail)
  }

  return nil
}
