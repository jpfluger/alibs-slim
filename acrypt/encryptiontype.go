package acrypt

import "strings"

// Supported EncryptionTypes
const (
	ENCRYPTIONTYPE_AES128 EncryptionType = "aes128"
	ENCRYPTIONTYPE_AES256 EncryptionType = "aes256"
)

// EncryptionType represents the type of encryption used.
type EncryptionType string

// IsEmpty checks if EncryptionType is empty after trimming whitespace.
func (et EncryptionType) IsEmpty() bool {
	return strings.TrimSpace(string(et)) == ""
}
