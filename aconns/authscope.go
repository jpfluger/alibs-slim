package aconns

import (
	"fmt"
	"strings"
)

const (
	AUTHSCOPE_MASTER AuthScope = "master"
	AUTHSCOPE_DOMAIN AuthScope = "domain"
	AUTHSCOPE_MODULE AuthScope = "module"
	AUTHSCOPE_ADMIN  AuthScope = "admin"
)

// AuthScope categorizes the level of the authentication system.
type AuthScope string

// IsValid checks if the AuthScope value is recognized.
func (at AuthScope) IsValid() bool {
	switch at.TrimSpace() {
	case AUTHSCOPE_MASTER, AUTHSCOPE_DOMAIN, AUTHSCOPE_MODULE, AUTHSCOPE_ADMIN:
		return true
	default:
		return false
	}
}

// IsEmpty returns true if the AuthScope is blank.
func (at AuthScope) IsEmpty() bool {
	return strings.TrimSpace(string(at)) == ""
}

// TrimSpace trims extra whitespace.
func (at AuthScope) TrimSpace() AuthScope {
	return AuthScope(strings.TrimSpace(string(at)))
}

// String returns a cleaned-up string.
func (at AuthScope) String() string {
	return string(at.TrimSpace())
}

// AuthScopes is a collection of AuthScope values.
type AuthScopes []AuthScope

// Has checks if the slice contains a specific tier.
func (ats AuthScopes) Has(target AuthScope) bool {
	for _, t := range ats {
		if t == target {
			return true
		}
	}
	return false
}

// Validate ensures all values are from known AuthScopes.
func (ats AuthScopes) Validate() error {
	for i, t := range ats {
		if !t.IsValid() {
			return fmt.Errorf("auth tier at index %d is invalid: %q", i, t)
		}
	}
	return nil
}
