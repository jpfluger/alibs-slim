package acrypt

import (
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	RSAKeySizeMin    = 256 // Example: Minimum size for RSA public keys
	Ed25519KeySize   = 32  // Size for Ed25519 public keys
	ECDSAP256KeySize = 64  // Size for P-256 ECDSA public keys
)

// CryptKeyBase64 represents a key where base64 is expected.
type CryptKeyBase64 string

func NewCryptKeyBase64(key []byte) CryptKeyBase64 {
	if !IsBase64Encoded(key) {
		return CryptKeyBase64(base64.StdEncoding.EncodeToString(key))
	}
	return CryptKeyBase64(key)
}

// IsEmpty checks if the PublicKey is empty.
func (pk CryptKeyBase64) IsEmpty() bool {
	return strings.TrimSpace(string(pk)) == ""
}

func (pk CryptKeyBase64) Encoded() string {
	return string(pk)
}

func (pk CryptKeyBase64) MustDecode() []byte {
	decoded, err := base64.StdEncoding.DecodeString(string(pk))
	if err != nil {
		return []byte{}
	}
	return decoded
}

func (pk CryptKeyBase64) Decoded() ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(pk))
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %v", err)
	}
	return decoded, nil
}

func (pk CryptKeyBase64) IsBase64() bool {
	if pk.IsEmpty() {
		return false
	}
	_, err := base64.StdEncoding.DecodeString(string(pk))
	return err == nil
}

// Validate checks if the key is valid base64 and matches the expected size.
func (pk CryptKeyBase64) Validate(expectedSize int) error {
	if pk.IsEmpty() {
		return fmt.Errorf("key is empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(string(pk))
	if err != nil {
		return fmt.Errorf("key is not valid base64: %v", err)
	}

	if len(decoded) != expectedSize {
		return fmt.Errorf("invalid key size: expected %d bytes, got %d bytes", expectedSize, len(decoded))
	}

	return nil
}
