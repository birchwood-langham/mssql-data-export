package dataexport

import (
	"crypto/sha256"
	"encoding/base64"
)

// Encrypt returns a SHA256 encrypted string from the text provided
func Encrypt(text string, secret string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text + secret))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}
