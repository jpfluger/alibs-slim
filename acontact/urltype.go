package acontact

import (
	"strings"
	// autils is a custom package for additional utility functions
	"github.com/jpfluger/alibs-slim/autils"
)

// UrlType defines a custom type for URL categorization
type UrlType string

// Constants for UrlType values
const (
	URLTYPE_HOME         = UrlType("home")         // URL associated with home
	URLTYPE_PERSONAL     = UrlType("personal")     // URL associated with an individual
	URLTYPE_BUSINESS     = UrlType("business")     // URL associated with a business
	URLTYPE_SOCIAL       = UrlType("social")       // Social media profile URL
	URLTYPE_PROFESSIONAL = UrlType("professional") // Professional networking profile URL
	URLTYPE_BLOG         = UrlType("blog")         // Personal or professional blog URL
	URLTYPE_PORTFOLIO    = UrlType("portfolio")    // Portfolio URL for showcasing work
	URLTYPE_COMPANY      = UrlType("company")      // Company or organization's main
	URLTYPE_LEGAL        = UrlType("legal")        // Legal or compliance documents (e.g. ToS, privacy)
)

// IsEmpty checks if the UrlType is empty after trimming spaces
func (ut UrlType) IsEmpty() bool {
	return strings.TrimSpace(string(ut)) == ""
}

// TrimSpace returns a new UrlType with leading and trailing spaces removed
func (ut UrlType) TrimSpace() UrlType {
	return UrlType(strings.TrimSpace(string(ut)))
}

// String returns the string representation of UrlType
func (ut UrlType) String() string {
	return string(ut)
}

// ToStringTrimLower returns the UrlType as a trimmed and lowercase string
func (ut UrlType) ToStringTrimLower() string {
	// Utilizes the autils package's ToStringTrimLower function
	return autils.ToStringTrimLower(string(ut))
}
