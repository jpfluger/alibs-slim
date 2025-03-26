package aconns

import (
	"errors"
	"strings"
	"unicode"
)

// AdapterType represents an adapter type
type AdapterType string

// IsEmpty checks if the AdapterType is empty
func (at AdapterType) IsEmpty() bool {
	return at == ""
}

// TrimSpace returns the trimmed string representation of the AdapterType
func (at AdapterType) TrimSpace() AdapterType {
	return AdapterType(strings.TrimSpace(string(at)))
}

// String returns the string representation of the AdapterType
func (at AdapterType) String() string {
	return string(at)
}

// Matches checks if the AdapterType matches the given string
func (at AdapterType) Matches(s string) bool {
	return string(at) == s
}

// ToStringTrimLower returns the trimmed and lowercased string representation of the AdapterType
func (at AdapterType) ToStringTrimLower() string {
	return strings.ToLower(strings.TrimSpace(string(at)))
}

// Validate checks if the AdapterType is valid
func (at AdapterType) Validate() error {
	trimmedLower := at.ToStringTrimLower()
	if trimmedLower == "" {
		return errors.New("AdapterType is empty")
	}

	for _, r := range trimmedLower {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return errors.New("AdapterType contains invalid characters")
		}
	}

	return nil
}

// AdapterTypes represents a slice of AdapterType
type AdapterTypes []AdapterType

// IsEmpty checks if the AdapterTypes slice is empty
func (ats AdapterTypes) IsEmpty() bool {
	return len(ats) == 0
}

// String returns the string representation of the AdapterTypes slice
func (ats AdapterTypes) String() string {
	return strings.Join(ats.ToStringArray(), ", ")
}

// ToStringArray returns an array of AdapterTypes as strings
func (ats AdapterTypes) ToStringArray() []string {
	strArray := make([]string, len(ats))
	for i, at := range ats {
		strArray[i] = at.String()
	}
	return strArray
}

// Find returns the AdapterType if found, otherwise an empty AdapterType
func (ats AdapterTypes) Find(at AdapterType) AdapterType {
	for _, v := range ats {
		if v == at {
			return v
		}
	}
	return ""
}

// HasKey checks if the AdapterTypes slice contains the given AdapterType
func (ats AdapterTypes) HasKey(s AdapterType) bool {
	return ats.Find(s) != ""
}

// Matches checks if any AdapterType in the AdapterTypes slice matches the given string
func (ats AdapterTypes) Matches(s string) bool {
	for _, v := range ats {
		if v.Matches(s) {
			return true
		}
	}
	return false
}
