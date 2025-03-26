package acontact

import (
	"strings"
	// Assuming autils is a custom package containing utility functions like ToStringTrimLower
	"github.com/jpfluger/alibs-slim/autils"
)

// EmailType defines a custom type for categorizing email addresses.
type EmailType string

// Constants for EmailType values.
const (
	EMAILTYPE_PERSONAL  = EmailType("personal")  // Personal email address
	EMAILTYPE_WORK      = EmailType("work")      // Work email address
	EMAILTYPE_BILLING   = EmailType("billing")   // Billing email address
	EMAILTYPE_SUPPORT   = EmailType("support")   // Support email address
	EMAILTYPE_INFO      = EmailType("info")      // General information email address
	EMAILTYPE_SALES     = EmailType("sales")     // Sales-related email address
	EMAILTYPE_MARKETING = EmailType("marketing") // Marketing email address
	EMAILTYPE_NO_REPLY  = EmailType("noreply")   // No-reply email address
)

// IsEmpty checks if the EmailType is empty after trimming spaces.
func (et EmailType) IsEmpty() bool {
	return strings.TrimSpace(string(et)) == ""
}

// TrimSpace returns a new EmailType with leading and trailing spaces removed.
func (et EmailType) TrimSpace() EmailType {
	return EmailType(strings.TrimSpace(string(et)))
}

// String converts the EmailType to a string.
func (et EmailType) String() string {
	return string(et)
}

// ToStringTrimLower returns the EmailType as a trimmed and lowercase string.
func (et EmailType) ToStringTrimLower() string {
	return autils.ToStringTrimLower(string(et))
}

// GetType extracts the type part before the colon in the EmailType.
func (et EmailType) GetType() string {
	if strings.Contains(string(et), ":") {
		parts := strings.Split(string(et), ":")
		return parts[0]
	}
	return string(et)
}

// GetPart extracts the part after the colon in the EmailType, if it exists.
func (et EmailType) GetPart() string {
	if strings.Contains(string(et), ":") {
		parts := strings.Split(string(et), ":")
		if len(parts) != 2 {
			return ""
		}
		return parts[1]
	}
	return ""
}
