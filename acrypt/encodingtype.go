package acrypt

import "strings"

// Supported EncodingTypes
const (
	ENCODINGTYPE_BASE64 EncodingType = "base64"
	ENCODINGTYPE_HEX    EncodingType = "hex"
	ENCODINGTYPE_PLAIN  EncodingType = "plain"
)

// EncodingType represents the type of encoding used.
type EncodingType string

// IsEmpty checks if EncodingType is empty after trimming whitespace.
func (et EncodingType) IsEmpty() bool {
	return strings.TrimSpace(string(et)) == ""
}
