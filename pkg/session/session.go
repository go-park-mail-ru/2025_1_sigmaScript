package session

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/pkg/errors"
)

const (
	MinSessionIDLength = 8
	MaxSessionIDLength = 512
)

func GenerateSessionID(length int) (string, error) {
	if length < 0 {
		return "", errors.New(errs.ErrMsgNegativeSessionIDLength)
	}
	if length < MinSessionIDLength {
		return "", errors.New(errs.ErrMsgLengthTooShort)
	}
	if length > MaxSessionIDLength {
		return "", errors.New(errs.ErrMsgLengthTooLong)
	}
	session := make([]byte, length)
	if _, err := rand.Read(session); err != nil {
		return "", errors.Wrap(err, errs.ErrMsgGenerateSession)
	}
	hash := sha256.Sum256(session)
	return hex.EncodeToString(hash[:]), nil
}
