package acontact

import (
	"strings"
	// Importing custom utilities package for string manipulation
	"github.com/jpfluger/alibs-slim/autils"
)

// ContactType defines a custom type for contact entities
type ContactType string

// Constants for ContactType values
const (
	CONTACTTYPE_PERSON = ContactType("person") // Represents a personal contact
	CONTACTTYPE_ENTITY = ContactType("entity") // Represents a company contact
)

// IsEmpty checks if the ContactType is empty after trimming spaces
func (ct ContactType) IsEmpty() bool {
	return strings.TrimSpace(string(ct)) == ""
}

// TrimSpace returns a new ContactType with leading and trailing spaces removed
func (ct ContactType) TrimSpace() ContactType {
	return ContactType(strings.TrimSpace(string(ct)))
}

// HasMatch checks if the target ContactType matches the current ContactType
func (ct ContactType) HasMatch(target ContactType) bool {
	return ct == target
}

// String returns the string representation of ContactType in trimmed and lowercase format
func (ct ContactType) String() string {
	return ct.ToStringTrimLower()
}

// ToStringTrimLower returns the ContactType as a trimmed and lowercase string
func (ct ContactType) ToStringTrimLower() string {
	return autils.ToStringTrimLower(string(ct))
}

// ContactTypes defines a slice of ContactType
type ContactTypes []ContactType
