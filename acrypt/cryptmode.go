package acrypt

import "strings"

const (
	CRYPTMODE_ENCRYPTED CryptMode = "e"
	CRYPTMODE_DECRYPTED CryptMode = "d"
)

type CryptMode string

// IsEmpty checks if CryptMode is empty after trimming whitespace.
func (cm CryptMode) IsEmpty() bool {
	return strings.TrimSpace(string(cm)) == ""
}
