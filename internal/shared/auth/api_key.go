package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// Generate a secure random API key
func GenerateAPIKey() (string, string, error) {
	keyBytes := make([]byte, 32) // 256-bit
	_, err := rand.Read(keyBytes)
	if err != nil {
		return "", "", err
	}

	rawKey := hex.EncodeToString(keyBytes)

	hash := sha256.Sum256([]byte(rawKey))
	hashStr := hex.EncodeToString(hash[:])

	return rawKey, hashStr, nil
}

// Hash incoming key (for comparison)
func HashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
