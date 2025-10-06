package auth

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateState generates a random state string for OAuth CSRF protection
func GenerateState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
