package auser

import (
	"strings"
)

// UsernameValidityType defines a string-based type for username validation types.
type UsernameValidityType string

// Constants for different types of username validation.
const (
	USERNAMEVALIDITYTYPE_EMAIL_OR_USER     UsernameValidityType = "email-or-user"
	USERNAMEVALIDITYTYPE_EMAIL             UsernameValidityType = "email"
	USERNAMEVALIDITYTYPE_USER              UsernameValidityType = "user"
	USERNAMEVALIDITYTYPE_USER_MINL1_MAXL99 UsernameValidityType = "user-minl1-maxl99"
)

// IsEmpty returns true if the UsernameValidityType is an empty string.
func (u UsernameValidityType) IsEmpty() bool {
	return strings.TrimSpace(string(u)) == ""
}

// String returns the trimmed, lowercase string representation of the type.
func (u UsernameValidityType) String() string {
	return strings.ToLower(strings.TrimSpace(string(u)))
}
