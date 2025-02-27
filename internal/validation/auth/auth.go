package auth

import (
  "strings"
  "unicode"
  "unicode/utf8"

  errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
  "github.com/pkg/errors"
)

const (
  MinPasswordLength = 6
)

func IsValidPassword(password string) error {
  if utf8.RuneCountInString(password) < MinPasswordLength {
    return errors.New(errs.ErrPasswordTooShort)
  }
  if strings.TrimSpace(password) == "" {
    return errors.New(errs.ErrEmptyPassword)
  }
  var lower, upper, digit bool
  for _, c := range password {
    switch {
    case unicode.IsNumber(c):
      digit = true
    case unicode.IsUpper(c):
      upper = true
    case unicode.IsLower(c):
      lower = true
    }
  }
  if lower && upper && digit {
    return nil
  }
  return errors.New(errs.ErrEasyPassword)
}
