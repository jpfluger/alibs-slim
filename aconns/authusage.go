package aconns

import (
	"fmt"
	"strings"
)

const (
	AUTHUSAGE_PRIMARY AuthUsage = "primary"
	AUTHUSAGE_MFA     AuthUsage = "mfa"
	AUTHUSAGE_SSPR    AuthUsage = "sspr"
)

// AuthUsage describes the purpose for which an authentication method is used.
type AuthUsage string

// IsValid returns true if the value is recognized.
func (ut AuthUsage) IsValid() bool {
	switch ut.TrimSpace() {
	case AUTHUSAGE_PRIMARY, AUTHUSAGE_MFA, AUTHUSAGE_SSPR:
		return true
	default:
		return false
	}
}

// IsEmpty checks for a blank value.
func (ut AuthUsage) IsEmpty() bool {
	return strings.TrimSpace(string(ut)) == ""
}

// TrimSpace removes excess space.
func (ut AuthUsage) TrimSpace() AuthUsage {
	return AuthUsage(strings.TrimSpace(string(ut)))
}

// String returns the string form.
func (ut AuthUsage) String() string {
	return string(ut.TrimSpace())
}

// AuthUsages is a list of usage types.
type AuthUsages []AuthUsage

// Has checks for existence in the slice.
func (uts AuthUsages) Has(target AuthUsage) bool {
	for _, t := range uts {
		if t == target {
			return true
		}
	}
	return false
}

// Validate confirms all entries are valid.
func (uts AuthUsages) Validate() error {
	for i, t := range uts {
		if !t.IsValid() {
			return fmt.Errorf("auth usage type at index %d is invalid: %q", i, t)
		}
	}
	return nil
}
