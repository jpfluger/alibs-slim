package acrypt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	r2 "math/rand/v2" // Use math/rand/v2 for IntN (Go 1.22+).
	"strings"
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
var GenerateSecretKey = func() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret key: %v", err)
	}
	return key, nil
}

func DecodeSecretKey(encodedKey string) ([]byte, error) {
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

func EncodeSecretKey(jwtDecodedKey []byte) (string, error) {
	if len(jwtDecodedKey) == 0 {
		return "", fmt.Errorf("decoded key is empty, cannot encode")
	}
	return base64.StdEncoding.EncodeToString(jwtDecodedKey), nil
}

// InitializeSecretKey initializes or returns a decoded secret key.
// If encodedKey is empty, it generates a new key; otherwise, it decodes the provided base64 string.
func InitializeSecretKey(encodedKey string) ([]byte, error) {
	return DecodeSecretKey(encodedKey)
}

// InitializeJWTSecretKey initializes or returns a decoded JWT secret key.
// If encodedKey is empty, it generates a new key; otherwise, it decodes the provided base64 string.
func InitializeJWTSecretKey(encodedKey string) ([]byte, error) {
	return DecodeSecretKey(encodedKey)
}

// EncodeJWTSecretKey takes a decoded JWT secret key and returns its base64-encoded string representation.
func EncodeJWTSecretKey(jwtDecodedKey []byte) (string, error) {
	return EncodeSecretKey(jwtDecodedKey)
}

// RandomStrongOneWaySecret generates a secure random secret.
func RandomStrongOneWaySecret() (string, error) {
	return RandomStrongOneWayByVariableLength(0, 0)
}

// RandomStrongOneWayByVariableLength generates a secure random key of variable length between low and high (inclusive).
// It first generates a fixed 32-byte random key, base64-encodes it (producing 44 characters), then trims to a random length in [low, high].
// Defaults to low=25, high=37 if low > high or invalid.
// For FIPS compliance, build with appropriate flags (e.g., GOEXPERIMENT=systemcrypto).
func RandomStrongOneWayByVariableLength(low, high int) (string, error) {
	if low > high || low < 1 || high < 1 {
		low = 25
		high = 37
	}

	// Generate 32 random bytes.
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("failed to generate secret key: %v", err)
	}

	// Base64 encode (44 characters).
	encoded := base64.StdEncoding.EncodeToString(key)

	// Pick random length in [low, high].
	varLen := r2.IntN(high-low+1) + low

	// Trim to varLen.
	return encoded[:varLen], nil
}

// GenerateEncryptionKeyWithLength generates a secure random encryption key of the specified length in bytes.
// The length must be 16, 24, or 32 (for AES-128, AES-192, or AES-256); otherwise, it defaults to 32.
// It fills the key with high-entropy bytes from crypto/rand.
// For FIPS compliance, build with appropriate flags (e.g., GOEXPERIMENT=systemcrypto).
func GenerateEncryptionKeyWithLength(selectedLen int) ([]byte, error) {
	// Validate and default to 32 if invalid
	if selectedLen != 16 && selectedLen != 24 && selectedLen != 32 {
		selectedLen = 32
	}
	key := make([]byte, selectedLen)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %v", err)
	}
	return key, nil
}

// GenerateEncryptionKeyWithLengthBase64 generates a secure random encryption key of the specified length in bytes (16, 24, or 32),
// base64-encodes it, and prefixes the result with "base64:". Defaults to 32 bytes if invalid length.
// This format is suitable for direct use in configurations that expect prefixed base64 keys.
func GenerateEncryptionKeyWithLengthBase64(selectedLen int) (string, error) {
	key, err := GenerateEncryptionKeyWithLength(selectedLen)
	if err != nil {
		return "", err
	}
	encoded, err := EncodeSecretKey(key)
	if err != nil {
		return "", err
	}
	return "base64:" + encoded, nil
}

// DecodePrefixedBase64 decodes a base64-encoded string that may be prefixed with "base64:".
// It strips the prefix if present and decodes the remaining string using standard base64 encoding.
// Returns the decoded bytes or an error if decoding fails.
func DecodePrefixedBase64(s string) ([]byte, error) {
	encoded := strings.TrimPrefix(s, "base64:")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("invalid base64-encoded string after 'base64:' prefix: %v", err)
	}
	return decoded, nil
}
