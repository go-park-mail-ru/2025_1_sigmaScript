package hashsaltfilename

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"path/filepath"
	"strconv"
	"time"
)

func HashSaltFilename(filename string) string {
	rand.NewSource(time.Now().UnixNano())

	ext := filepath.Ext(filename)
	nameWithoutExt := filename[:len(filename)-len(ext)]
	salt := rand.Intn(1 << 8)

	hashBytes := sha256.Sum256([]byte(nameWithoutExt))
	hashedName := strconv.Itoa(salt) + hex.EncodeToString(hashBytes[:]) + ext

	return hashedName
}
