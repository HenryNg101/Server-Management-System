package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// Generate a secure random API key
func GenerateAPIKey() (string, string, error) {
	keyBytes := make([]byte, 32) // 256-bit
	if _, err := rand.Read(keyBytes); err != nil {
		return "", "", err
	}

	rawKey := hex.EncodeToString(keyBytes)

	return rawKey, HashAPIKey(rawKey), nil
}

// Hash incoming key (for comparison)
func HashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
