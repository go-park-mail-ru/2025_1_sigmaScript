package session

import (
  "crypto/rand"
  "crypto/sha256"
  "encoding/hex"
  "fmt"
)

func GenerateSessionID(length int) (string, error) {
  session := make([]byte, length)
  if _, err := rand.Read(session); err != nil {
    return "", fmt.Errorf("failed to generate session id: %w", err)
  }
  hash := sha256.Sum256(session)
  return hex.EncodeToString(hash[:]), nil
}
