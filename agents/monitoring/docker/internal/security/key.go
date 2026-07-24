package security

import (
	"crypto/sha256"
	"os"
	"strings"
)

// getMachineID returns system hostname
func getMachineID() string {
	data, err := os.ReadFile("/host/machine-id")
	if err != nil {
		return "fallback"
	}
	return strings.TrimSpace(string(data))
}

// DeriveKey creates a 32-byte AES key from some stable input (e.g. instanceID).
// SHA256 ensures:
// - fixed length (32 bytes → AES-256)
// - deterministic (same input → same key)
// - not reversible (you can't get original input from hash)
func DeriveKey(instanceID string) []byte {
	material := instanceID + ":" + getMachineID()
	hash := sha256.Sum256([]byte(material))
	return hash[:] // convert [32]byte → []byte
}
