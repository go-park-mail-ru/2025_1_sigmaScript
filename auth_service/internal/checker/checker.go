package checker

import (
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/internal/errors"
)

// ValidateCookie validates cookie presence
func ValidateCookie(cookie string) error {
	if cookie == "" {
		return errs.ErrInvalidCookie
	}

	return nil
}

// ValidateUserID validates userID presense
func ValidateUserID(userID string) error {
	if userID == "" {
		return errs.ErrInvalidUserID
	}

	return nil
}
