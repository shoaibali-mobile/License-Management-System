package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateAPIKey generates a random API key
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 16) // 32 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "sk-sdk-" + hex.EncodeToString(bytes), nil
}





