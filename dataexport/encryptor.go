package dataexport

import (
	"crypto/sha256"
	"encoding/hex"
)

// Encrypt returns a SHA256 encrypted string from the text provided
func Encrypt(text string, secret string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text + secret))
	return hex.EncodeToString(hasher.Sum(nil))
}
