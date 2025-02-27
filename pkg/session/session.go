package session

import (
  "crypto/rand"
  "crypto/sha256"
  "encoding/hex"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
  "github.com/pkg/errors"
)

func GenerateSessionID(length int) (string, error) {
  session := make([]byte, length)
  if _, err := rand.Read(session); err != nil {
    return "", errors.Wrap(err, errs.ErrGenerateSession)
  }
  hash := sha256.Sum256(session)
  return hex.EncodeToString(hash[:]), nil
}
