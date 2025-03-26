package atemplates

import (
	"github.com/jpfluger/alibs-slim/autils" // Importing custom utility package
	"strings"                                              // Importing strings package for string manipulation
)

// ParseDirectiveType defines a new type based on string to work with directive parsing.
type ParseDirectiveType string

// IsEmpty checks if the ParseDirectiveType is empty after trimming space.
func (pdType ParseDirectiveType) IsEmpty() bool {
	trimmed := strings.TrimSpace(string(pdType)) // Trim space from pdType and convert to string
	return trimmed == ""                         // Return true if trimmed string is empty
}

// TrimSpace trims the spaces from ParseDirectiveType and returns a new ParseDirectiveType.
func (pdType ParseDirectiveType) TrimSpace() ParseDirectiveType {
	trimmed := strings.TrimSpace(string(pdType)) // Trim space from pdType and convert to string
	return ParseDirectiveType(trimmed)           // Return new ParseDirectiveType with trimmed spaces
}

// String converts ParseDirectiveType to string.
func (pdType ParseDirectiveType) String() string {
	return string(pdType) // Convert pdType to string and return
}

// ToStringTrimLower trims spaces from ParseDirectiveType, converts it to lower case, and returns as string.
func (pdType ParseDirectiveType) ToStringTrimLower() string {
	return autils.ToStringTrimLower(string(pdType)) // Use custom utility function to trim, lower, and convert to string
}
