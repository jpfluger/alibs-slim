package aconns

import (
	"fmt"
	"strings"
)

// Constants for ConnRole values.
const (
	CONNROLE_MASTER ConnRole = "master"
	CONNROLE_AUTH   ConnRole = "auth"
	CONNROLE_TENANT ConnRole = "tenant"
)

type ConnRole string

// IsEmpty returns true if the ConnRole is empty or only whitespace.
func (cr ConnRole) IsEmpty() bool {
	return strings.TrimSpace(string(cr)) == ""
}

// TrimSpace returns a whitespace-trimmed version of the ConnRole.
func (cr ConnRole) TrimSpace() ConnRole {
	return ConnRole(strings.TrimSpace(string(cr)))
}

// String returns the trimmed string representation of the ConnRole.
func (cr ConnRole) String() string {
	return cr.TrimSpace().StringRaw()
}

// StringRaw returns the untrimmed string representation.
func (cr ConnRole) StringRaw() string {
	return string(cr)
}

// ConnRoles is a slice of ConnRole.
type ConnRoles []ConnRole

// Validate checks whether all ConnRoles are recognized and non-empty.
// Returns an error if any role is invalid or empty.
func (crs ConnRoles) Validate() error {
	for i, role := range crs {
		trimmed := role.TrimSpace()
		if trimmed.IsEmpty() {
			return fmt.Errorf("conn role at index %d is empty or whitespace", i)
		}
		if !trimmed.IsValid() {
			return fmt.Errorf("conn role at index %d has invalid value: %q", i, role)
		}
	}
	return nil
}

// HasRole checks if the specified ConnRole exists in the list.
func (crs ConnRoles) HasRole(connRole ConnRole) bool {
	needle := connRole.TrimSpace()
	for _, r := range crs {
		if r.TrimSpace() == needle {
			return true
		}
	}
	return false
}

// IsValid returns true if the ConnRole is one of the defined constants.
func (cr ConnRole) IsValid() bool {
	switch cr.TrimSpace() {
	case CONNROLE_MASTER, CONNROLE_AUTH, CONNROLE_TENANT:
		return true
	default:
		return false
	}
}
