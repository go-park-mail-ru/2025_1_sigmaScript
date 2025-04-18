package csrftoken

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/pkg/errors"
)

const (
	CSRF_TOKEN_LENGTH = 32
)

func GenerateCSRFToken() (string, error) {
	tokenCSRF := make([]byte, CSRF_TOKEN_LENGTH)
	if _, err := rand.Read(tokenCSRF); err != nil {
		return "", errors.Wrap(err, errs.ErrMsgGenerateCSRFToken)
	}
	hash := sha256.Sum256(tokenCSRF)
	return hex.EncodeToString(hash[:]), nil
}
