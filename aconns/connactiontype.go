package aconns

import (
	"errors"
	"strings"
	"unicode"
)

const (
	CONNACTIONTYPE_CREATE     ConnActionType = "CREATE"
	CONNACTIONTYPE_UPGRADE    ConnActionType = "UPGRADE"
	CONNACTIONTYPE_DOWNGRADE  ConnActionType = "DOWNGRADE"
	CONNACTIONTYPE_DELETE     ConnActionType = "DELETE"
	CONNACTIONTYPE_NOOP_EQUAL ConnActionType = "no-op-equal" // expect no-operation to be performed because DB is in an equal state
)

// ConnActionType represents an adapter type.
type ConnActionType string

// IsEmpty checks if the ConnActionType is empty.
func (csa ConnActionType) IsEmpty() bool {
	return csa == ""
}

// TrimSpace returns the trimmed string representation of the ConnActionType.
func (csa ConnActionType) TrimSpace() ConnActionType {
	return ConnActionType(strings.TrimSpace(string(csa)))
}

// String returns the string representation of the ConnActionType.
func (csa ConnActionType) String() string {
	return string(csa)
}

// Matches checks if the ConnActionType matches the given string.
func (csa ConnActionType) Matches(s string) bool {
	return string(csa) == s
}

// ToStringTrimLower returns the trimmed and lowercased string representation of the ConnActionType.
func (csa ConnActionType) ToStringTrimLower() string {
	return strings.ToLower(strings.TrimSpace(string(csa)))
}

// Validate checks if the ConnActionType is valid.
func (csa ConnActionType) Validate() error {
	trimmedLower := csa.ToStringTrimLower()
	if trimmedLower == "" {
		return errors.New("ConnActionType is empty")
	}

	for _, r := range trimmedLower {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return errors.New("ConnActionType contains invalid characters")
		}
	}

	return nil
}

// ConnActionTypes represents a slice of ConnActionType.
type ConnActionTypes []ConnActionType

// IsEmpty checks if the ConnActionTypes slice is empty.
func (csas ConnActionTypes) IsEmpty() bool {
	return len(csas) == 0
}

// String returns the string representation of the ConnActionTypes slice.
func (csas ConnActionTypes) String() string {
	return strings.Join(csas.ToStringArray(), ", ")
}

// ToStringArray returns an array of ConnActionTypes as strings.
func (csas ConnActionTypes) ToStringArray() []string {
	strArray := make([]string, len(csas))
	for i, csa := range csas {
		strArray[i] = csa.String()
	}
	return strArray
}

// Find returns the ConnActionType if found, otherwise an empty ConnActionType.
func (csas ConnActionTypes) Find(csa ConnActionType) ConnActionType {
	for _, v := range csas {
		if v == csa {
			return v
		}
	}
	return ""
}

// HasKey checks if the ConnActionTypes slice contains the given ConnActionType.
func (csas ConnActionTypes) HasKey(s ConnActionType) bool {
	return csas.Find(s) != ""
}

// Matches checks if any ConnActionType in the ConnActionTypes slice matches the given string.
func (csas ConnActionTypes) Matches(s string) bool {
	for _, v := range csas {
		if v.Matches(s) {
			return true
		}
	}
	return false
}
