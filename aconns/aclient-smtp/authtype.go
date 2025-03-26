package aclient_smtp

import (
	"errors"
	"strings"
	"unicode"
)

// AuthType represents an authentication type
type AuthType string

// IsEmpty checks if the AuthType is empty
func (at AuthType) IsEmpty() bool {
	return at == ""
}

// TrimSpace returns the trimmed string representation of the AuthType
func (at AuthType) TrimSpace() AuthType {
	return AuthType(strings.TrimSpace(string(at)))
}

// String returns the string representation of the AuthType
func (at AuthType) String() string {
	return string(at)
}

// Matches checks if the AuthType matches the given string
func (at AuthType) Matches(s string) bool {
	return string(at) == s
}

// ToStringTrimLower returns the trimmed and lowercased string representation of the AuthType
func (at AuthType) ToStringTrimLower() string {
	return strings.ToLower(strings.TrimSpace(string(at)))
}

// Validate checks if the AuthType is valid
func (at AuthType) Validate() error {
	trimmedLower := at.ToStringTrimLower()
	if trimmedLower == "" {
		return errors.New("AuthType is empty")
	}

	for _, r := range trimmedLower {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return errors.New("AuthType contains invalid characters")
		}
	}

	return nil
}

// AuthTypes represents a slice of AuthType
type AuthTypes []AuthType

// IsEmpty checks if the AuthTypes slice is empty
func (ats AuthTypes) IsEmpty() bool {
	return len(ats) == 0
}

// String returns the string representation of the AuthTypes slice
func (ats AuthTypes) String() string {
	return strings.Join(ats.ToStringArray(), ", ")
}

// ToStringArray returns an array of AuthTypes as strings
func (ats AuthTypes) ToStringArray() []string {
	strArray := make([]string, len(ats))
	for i, at := range ats {
		strArray[i] = at.String()
	}
	return strArray
}

// Find returns the AuthType if found, otherwise an empty AuthType
func (ats AuthTypes) Find(at AuthType) AuthType {
	for _, v := range ats {
		if v == at {
			return v
		}
	}
	return ""
}

// HasKey checks if the AuthTypes slice contains the given AuthType
func (ats AuthTypes) HasKey(s AuthType) bool {
	return ats.Find(s) != ""
}

// Matches checks if any AuthType in the AuthTypes slice matches the given string
func (ats AuthTypes) Matches(s string) bool {
	for _, v := range ats {
		if v.Matches(s) {
			return true
		}
	}
	return false
}
