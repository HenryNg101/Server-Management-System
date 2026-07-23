package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// Encrypt takes plaintext (API key) and encrypts it using AES-GCM.
//
// AES-GCM gives you:
// - confidentiality (data is encrypted)
// - integrity (tampering will be detected automatically)
//
// Output format:
// base64( nonce || ciphertext )
//
// where:
// - nonce = random bytes (unique per encryption)
// - ciphertext = encrypted data
func Encrypt(plainText string, key []byte) (string, error) {

	// 1. Create AES cipher block using key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 2. Wrap it with GCM mode (adds authentication + integrity)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 3. Generate a random nonce (required for GCM)
	// nonce MUST be unique per encryption
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 4. Encrypt the plaintext
	// Seal() returns: nonce + encrypted data + auth tag
	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)

	// 5. Encode to base64 so it can be stored as string in JSON
	return base64.StdEncoding.EncodeToString(cipherText), nil
}
