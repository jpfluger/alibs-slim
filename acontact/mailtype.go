package acontact

import (
	"strings"
	// Assuming autils is a custom package containing utility functions like ToStringTrimLower
	"github.com/jpfluger/alibs-slim/autils"
)

// MailType defines a custom type for categorizing mailing addresses.
type MailType string

// Constants for MailType values.
const (
	MAILTYPE_HOME      MailType = "home"      // Home mailing address
	MAILTYPE_WORK      MailType = "work"      // Work mailing address
	MAILTYPE_BILLING   MailType = "billing"   // Billing mailing address
	MAILTYPE_SHIPPING  MailType = "shipping"  // Shipping mailing address
	MAILTYPE_POBOX     MailType = "pobox"     // Post office box mailing address
	MAILTYPE_ALTERNATE MailType = "alternate" // Alternate mailing address
	MAILTYPE_LEGAL     MailType = "legal"     // Legal mailing address
)

// IsEmpty checks if the MailType is empty after trimming spaces.
func (mt MailType) IsEmpty() bool {
	return strings.TrimSpace(string(mt)) == ""
}

// TrimSpace returns a new MailType with leading and trailing spaces removed.
func (mt MailType) TrimSpace() MailType {
	return MailType(strings.TrimSpace(string(mt)))
}

// String converts the MailType to a string.
func (mt MailType) String() string {
	return string(mt)
}

// ToStringTrimLower returns the MailType as a trimmed and lowercase string.
func (mt MailType) ToStringTrimLower() string {
	return autils.ToStringTrimLower(string(mt))
}
