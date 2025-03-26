package acrypt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/chacha20"
)

// GenerateChaCha20Key generates a secure 256-bit secret key using ChaCha20.
// The RandGenerate32 function generates a 32-character random string, which can be used as a secret key. However, there are some differences between generating a random string and generating a secure byte array for cryptographic purposes.
//
// Differences:
//   - Character Set: The RandGenerate32 function generates a string using a specific set of characters, which may include letters and digits. This is suitable for passwords but might not be as secure as a byte array generated using crypto/rand for cryptographic keys.
//   - Entropy: A byte array generated using crypto/rand has higher entropy and is more suitable for cryptographic purposes.
func GenerateChaCha20Key() ([]byte, error) {
	key := make([]byte, chacha20.KeySize) // ChaCha20 key size is 32 bytes (256 bits)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ChaCha20 key: %v", err)
	}
	return key, nil
}

// InitializeJWTSecretKey initializes or returns a decoded JWT secret key.
func InitializeJWTSecretKey(encodedKey string) (jwtDecodedKey []byte, err error) {
	var secretKey []byte
	if encodedKey == "" {
		secretKey, err = GenerateChaCha20Key()
		if err != nil {
			return nil, fmt.Errorf("failed to generate ChaCha20 key for jwtDecodedKey: %v", err)
		}
	} else {
		secretKey, err = base64.StdEncoding.DecodeString(encodedKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decode jwtSecretKey: %v", err)
		}
	}
	jwtDecodedKey = secretKey
	return
}

// EncodeJWTSecretKey takes a decoded JWT secret key and returns an encoded string.
func EncodeJWTSecretKey(jwtDecodedKey []byte) (string, error) {
	// Check if the decoded key is valid
	if len(jwtDecodedKey) == 0 {
		return "", fmt.Errorf("decoded key is empty, cannot encode")
	}

	// Encode the key to a Base64 string
	encodedKey := base64.StdEncoding.EncodeToString(jwtDecodedKey)
	return encodedKey, nil
}
