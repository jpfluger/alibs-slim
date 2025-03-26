package acontact

import (
	"strings"
	// Importing custom utilities package for string manipulation
	"github.com/jpfluger/alibs-slim/autils"
)

type DetailsType string

// IsEmpty checks if the DetailsType is empty after trimming spaces
func (ct DetailsType) IsEmpty() bool {
	return strings.TrimSpace(string(ct)) == ""
}

// TrimSpace returns a new DetailsType with leading and trailing spaces removed
func (ct DetailsType) TrimSpace() DetailsType {
	return DetailsType(strings.TrimSpace(string(ct)))
}

// HasMatch checks if the target DetailsType matches the current DetailsType
func (ct DetailsType) HasMatch(target DetailsType) bool {
	return ct == target
}

// String returns the string representation of DetailsType in trimmed and lowercase format
func (ct DetailsType) String() string {
	return ct.ToStringTrimLower()
}

// ToStringTrimLower returns the DetailsType as a trimmed and lowercase string
func (ct DetailsType) ToStringTrimLower() string {
	return autils.ToStringTrimLower(string(ct))
}

// DetailsTypes defines a slice of DetailsType
type DetailsTypes []DetailsType
