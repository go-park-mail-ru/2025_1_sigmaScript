package auth

import (
  "strings"
  "unicode"
  "unicode/utf8"

  errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
  "github.com/pkg/errors"
)

func IsValidPassword(password string) error {
  if utf8.RuneCountInString(password) < 6 {
    return errors.New(errs.ErrPasswordTooShort)
  }
  if strings.TrimSpace(password) == "" {
    return errors.New(errs.ErrEmptyPassword)
  }
  lower := false
  upper := false
  digit := false
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
  if !lower || !upper || !digit {
    return errors.New(errs.ErrEasyPassword)
  }
  return nil
}
