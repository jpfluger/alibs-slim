package aclient_smtp

import (
	"errors"
	"strings"
	"unicode"
)

const (
	DIALMODE_UNKNOWN  DialMode = "unknown"   // Auto-detect the connection mode
	DIALMODE_NOTLS    DialMode = "no-tls"    // No encryption
	DIALMODE_TLS      DialMode = "tls"       // Implicit TLS
	DIALMODE_STARTTLS DialMode = "start-tls" // Explicit STARTTLS
)

// DialMode represents an authentication type
type DialMode string

// IsEmpty checks if the DialMode is empty
func (at DialMode) IsEmpty() bool {
	return at == ""
}

// TrimSpace returns the trimmed string representation of the DialMode
func (at DialMode) TrimSpace() DialMode {
	return DialMode(strings.TrimSpace(string(at)))
}

// String returns the string representation of the DialMode
func (at DialMode) String() string {
	return string(at)
}

// Matches checks if the DialMode matches the given string
func (at DialMode) Matches(s string) bool {
	return string(at) == s
}

// ToStringTrimLower returns the trimmed and lowercased string representation of the DialMode
func (at DialMode) ToStringTrimLower() string {
	return strings.ToLower(strings.TrimSpace(string(at)))
}

// Validate checks if the DialMode is valid
func (at DialMode) Validate() error {
	trimmedLower := at.ToStringTrimLower()
	if trimmedLower == "" {
		return errors.New("DialMode is empty")
	}

	for _, r := range trimmedLower {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return errors.New("DialMode contains invalid characters")
		}
	}

	return nil
}

// DialModes represents a slice of DialMode
type DialModes []DialMode

// IsEmpty checks if the DialModes slice is empty
func (ats DialModes) IsEmpty() bool {
	return len(ats) == 0
}

// String returns the string representation of the DialModes slice
func (ats DialModes) String() string {
	return strings.Join(ats.ToStringArray(), ", ")
}

// ToStringArray returns an array of DialModes as strings
func (ats DialModes) ToStringArray() []string {
	strArray := make([]string, len(ats))
	for i, at := range ats {
		strArray[i] = at.String()
	}
	return strArray
}

// Find returns the DialMode if found, otherwise an empty DialMode
func (ats DialModes) Find(at DialMode) DialMode {
	for _, v := range ats {
		if v == at {
			return v
		}
	}
	return ""
}

// HasKey checks if the DialModes slice contains the given DialMode
func (ats DialModes) HasKey(s DialMode) bool {
	return ats.Find(s) != ""
}

// Matches checks if any DialMode in the DialModes slice matches the given string
func (ats DialModes) Matches(s string) bool {
	for _, v := range ats {
		if v.Matches(s) {
			return true
		}
	}
	return false
}
