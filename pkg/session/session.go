package session

import (
  "crypto/rand"
  "encoding/base64"
  "fmt"
)

func GenerateSessionID(length int) (string, error) {
  session := make([]byte, length)
  if _, err := rand.Read(session); err != nil {
    return "", fmt.Errorf("failed to generate session id: %w", err)
  }
  return base64.URLEncoding.EncodeToString(session), nil
}
