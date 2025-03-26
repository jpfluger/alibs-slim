package acontact

import (
	"strings"
	// autils is assumed to be a custom package containing utility functions like ToStringTrimLower
	"github.com/jpfluger/alibs-slim/autils"
)

// PhoneType defines a custom type for categorizing phone numbers.
type PhoneType string

// Constants for PhoneType values.
const (
	PHONETYPE_MOBILE   = PhoneType("mobile")   // Mobile phone number
	PHONETYPE_HOME     = PhoneType("home")     // Home phone number
	PHONETYPE_WORK     = PhoneType("work")     // Work phone number
	PHONETYPE_DIRECT   = PhoneType("direct")   // Direct line phone number
	PHONETYPE_FAX      = PhoneType("fax")      // Fax number
	PHONETYPE_OTHER    = PhoneType("other")    // Other types of phone numbers
	PHONETYPE_PERSONAL = PhoneType("personal") // Personal phone number
)

// IsEmpty checks if the PhoneType is empty after trimming spaces.
func (pt PhoneType) IsEmpty() bool {
	return strings.TrimSpace(string(pt)) == ""
}

// TrimSpace returns a new PhoneType with leading and trailing spaces removed.
func (pt PhoneType) TrimSpace() PhoneType {
	return PhoneType(strings.TrimSpace(string(pt)))
}

// String converts the PhoneType to a string.
func (pt PhoneType) String() string {
	return string(pt)
}

// ToStringTrimLower returns the PhoneType as a trimmed and lowercase string.
func (pt PhoneType) ToStringTrimLower() string {
	return autils.ToStringTrimLower(string(pt))
}
