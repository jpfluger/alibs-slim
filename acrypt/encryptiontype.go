package acrypt

import "strings"

// Supported EncryptionTypes (all FIPS 140-3 compliant; default to AES-256 for quantum resistance)
const (
	ENCRYPTIONTYPE_AES128 EncryptionType = "aes128" // 128-bit key; quantum-effective: 64 bits (use cautiously)
	ENCRYPTIONTYPE_AES192 EncryptionType = "aes192" // 192-bit key; quantum-effective: 96 bits
	ENCRYPTIONTYPE_AES256 EncryptionType = "aes256" // 256-bit key; quantum-effective: 128 bits (recommended default)
)

// EncryptionType represents the type of encryption used.
type EncryptionType string

// IsEmpty checks if EncryptionType is empty after trimming whitespace.
func (et EncryptionType) IsEmpty() bool {
	return strings.TrimSpace(string(et)) == ""
}

// KeySize returns the key size in bytes for the EncryptionType.
func (et EncryptionType) KeySize() int {
	switch et {
	case ENCRYPTIONTYPE_AES128:
		return 16
	case ENCRYPTIONTYPE_AES192:
		return 24
	case ENCRYPTIONTYPE_AES256:
		return 32
	default:
		return 0 // Invalid
	}
}
