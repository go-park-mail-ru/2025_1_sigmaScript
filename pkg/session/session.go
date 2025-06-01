package session

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/pkg/errors"
)

const (
	MinSessionIDLength = 8
	MaxSessionIDLength = 512

	ErrMsgNegativeSessionIDLength = "negative session ID length"
	ErrMsgLengthTooShort          = "length too short"
	ErrMsgLengthTooLong           = "length too long"
	ErrMsgFailedToGetSession      = "failed to get session"
	ErrMsgGenerateSession         = "error generating session ID"
)

func GenerateSessionID(length int) (string, error) {
	if length < 0 {
		return "", errors.New(ErrMsgNegativeSessionIDLength)
	}
	if length < MinSessionIDLength {
		return "", errors.New(ErrMsgLengthTooShort)
	}
	if length > MaxSessionIDLength {
		return "", errors.New(ErrMsgLengthTooLong)
	}
	session := make([]byte, length)
	if _, err := rand.Read(session); err != nil {
		return "", errors.Wrap(err, ErrMsgGenerateSession)
	}
	hash := sha256.Sum256(session)
	return hex.EncodeToString(hash[:]), nil
}
