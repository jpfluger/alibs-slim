package aconns

import (
	"fmt"
	"strings"
)

const (
	AUTHMETHOD_PRIMARY AuthMethod = "primary"
	AUTHMETHOD_MFA     AuthMethod = "mfa"
	AUTHMETHOD_SSPR    AuthMethod = "sspr"
)

// AuthMethod represents a method of authentication supported by a connection.
type AuthMethod string

// AuthMethods is a slice of AuthMethod values, used to represent all authentication methods supported by a connection.
type AuthMethods []AuthMethod

// IsValid returns true if the AuthMethod is a recognized constant.
func (am AuthMethod) IsValid() bool {
	switch am.TrimSpace() {
	case AUTHMETHOD_PRIMARY, AUTHMETHOD_MFA, AUTHMETHOD_SSPR:
		return true
	default:
		return false
	}
}

// TrimSpace returns the AuthMethod with surrounding whitespace removed.
func (am AuthMethod) TrimSpace() AuthMethod {
	return AuthMethod(strings.TrimSpace(string(am)))
}

// String returns the string representation of the AuthMethod.
func (am AuthMethod) String() string {
	return am.TrimSpace().Raw()
}

// Raw returns the underlying string of the AuthMethod without trimming.
func (am AuthMethod) Raw() string {
	return string(am)
}

// Has returns true if the AuthMethods list contains the specified method (case-sensitive).
func (ams AuthMethods) Has(target AuthMethod) bool {
	for _, m := range ams {
		if m == target {
			return true
		}
	}
	return false
}

// Validate checks that all AuthMethod entries are recognized constants.
func (ams AuthMethods) Validate() error {
	for i, m := range ams {
		if !m.IsValid() {
			return fmt.Errorf("auth method at index %d is invalid: %q", i, m)
		}
	}
	return nil
}
