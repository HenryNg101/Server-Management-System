package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// Decrypt reverses Encrypt()
//
// Steps:
// - base64 decode
// - extract nonce
// - decrypt using same key
// - verify integrity automatically
func Decrypt(cipherText string, key []byte) (string, error) {

	// 1. Decode base64 back to raw bytes
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	// 2. Recreate AES + GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 3. Extract nonce
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", err
	}

	nonce := data[:nonceSize]
	cipherData := data[nonceSize:]

	// 4. Decrypt + verify integrity
	// If data was tampered → this FAILS
	plain, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
