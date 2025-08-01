package acrypt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateSecretKey generates a secure 256-bit (32-byte) random key suitable for FIPS-approved cryptographic uses,
// such as HMAC-SHA256 for JWT signing.
// It uses crypto/rand for high-entropy generation.
//
// Compatibility note:
//   - In Go 1.23 and earlier, this function may return an error if randomness generation fails.
//   - In Go 1.24 and later, crypto/rand.Read is guaranteed to succeed or panic (it never returns a non-nil error).
//     Thus, the error return is always nil in Go 1.24+, and the error-handling branch is dead code but harmless.
//   - For FIPS 140-3 compliance, build with appropriate flags (e.g., GOEXPERIMENT=systemcrypto) to use certified modules.
func GenerateSecretKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret key: %v", err)
	}
	return key, nil
}

// InitializeJWTSecretKey initializes or returns a decoded JWT secret key.
// If encodedKey is empty, it generates a new key; otherwise, it decodes the provided base64 string.
func InitializeJWTSecretKey(encodedKey string) ([]byte, error) {
	if encodedKey == "" {
		secretKey, err := GenerateSecretKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate secret key for JWT: %v", err)
		}
		return secretKey, nil
	}
	secretKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT secret key: %v", err)
	}
	return secretKey, nil
}

// EncodeJWTSecretKey takes a decoded JWT secret key and returns its base64-encoded string representation.
func EncodeJWTSecretKey(jwtDecodedKey []byte) (string, error) {
	if len(jwtDecodedKey) == 0 {
		return "", fmt.Errorf("decoded key is empty, cannot encode")
	}
	return base64.StdEncoding.EncodeToString(jwtDecodedKey), nil
}
