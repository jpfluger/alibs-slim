package acontact

import (
	"strings"
	// Importing custom utilities package for string manipulation
	"github.com/jpfluger/alibs-slim/autils"
)

// PersonType defines a custom type for contact entities
type PersonType string

const (
	PERSONTYPE_CUSTOMER  = PersonType("customer")  // Represents a customer
	PERSONTYPE_COLLEAGUE = PersonType("colleague") // Represents a colleague
	PERSONTYPE_ASSOCIATE = PersonType("associate") // Represents an associate
	PERSONTYPE_FRIEND    = PersonType("friend")    // Represents a friend
	PERSONTYPE_RELATIVE  = PersonType("relative")  // Represents a relative
	PERSONTYPE_PERSON    = PersonType("person")    // Represents a general person contact
)

// IsEmpty checks if the PersonType is empty after trimming spaces
func (ct PersonType) IsEmpty() bool {
	return strings.TrimSpace(string(ct)) == ""
}

// TrimSpace returns a new PersonType with leading and trailing spaces removed
func (ct PersonType) TrimSpace() PersonType {
	return PersonType(strings.TrimSpace(string(ct)))
}

// HasMatch checks if the target PersonType matches the current PersonType
func (ct PersonType) HasMatch(target PersonType) bool {
	return ct == target
}

// String returns the string representation of PersonType in trimmed and lowercase format
func (ct PersonType) String() string {
	return ct.ToStringTrimLower()
}

// ToStringTrimLower returns the PersonType as a trimmed and lowercase string
func (ct PersonType) ToStringTrimLower() string {
	return autils.ToStringTrimLower(string(ct))
}

// PersonTypes defines a slice of PersonType
type PersonTypes []PersonType
