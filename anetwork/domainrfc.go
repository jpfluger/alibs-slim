package anetwork

import (
	"strings"

	"golang.org/x/net/idna"
)

// DomainRFC represents a domain string with extended validation and helper methods.
type DomainRFC string

// IsEmpty checks if the domain is empty or consists of only whitespace.
func (d DomainRFC) IsEmpty() bool {
	return strings.TrimSpace(string(d)) == ""
}

// String returns the string representation of the domain.
func (d DomainRFC) String() string {
	return string(d)
}

// IsValid checks if the domain is valid according to external DNS rules.
// This means no IPs are allowed and no single-label hosts.
// Wildcard prefixes are ignored during validation.
func (d DomainRFC) IsValid() (bool, error) {
	return IsValidDomainWithError(d.String(), false)
}

// IsValidWithOptions checks if the domain is valid according to DNS rules.
// If allowIPs is true, then IPs are allowed and single-label hosts.
// Wildcard prefixes are ignored during validation.
func (d DomainRFC) IsValidWithOptions(allowIPs bool) (bool, error) {
	return IsValidDomainWithError(d.String(), allowIPs)
}

// ToUnicode converts a Punycode domain to its Unicode representation.
func (d DomainRFC) ToUnicode() (string, error) {
	return idna.ToUnicode(d.String())
}

// ToASCII converts a Unicode domain to its Punycode (ASCII) representation.
func (d DomainRFC) ToASCII() (string, error) {
	return idna.ToASCII(d.String())
}

// Normalize trims whitespace and converts the domain to its ASCII representation.
func (d DomainRFC) Normalize() (DomainRFC, error) {
	asciiDomain, err := d.ToASCII()
	if err != nil {
		return "", err
	}
	return DomainRFC(strings.TrimSpace(strings.ToLower(asciiDomain))), nil
}

// GetSlugReverse generates a slug by reversing the domain parts and joining with "/".
// For example, "hipaa.usa.gov" becomes "gov/usa/hipaa".
// Returns an empty string if the domain is empty.
func (d DomainRFC) GetSlugReverse() string {
	cleaned := strings.TrimSpace(string(d))
	if cleaned == "" {
		return ""
	}
	parts := strings.Split(cleaned, ".")
	// Reverse the parts slice
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, "/")
}

// DomainRFCs represents a slice of DomainRFC.
type DomainRFCs []DomainRFC

// FilterInvalid removes invalid domains from the slice.
func (ds DomainRFCs) FilterInvalid() (valid DomainRFCs, invalid map[DomainRFC]error) {
	return ds.FilterInvalidWithErrors(false)
}

// FilterInvalidWithErrors removes invalid domains and provides a map of invalid domains with errors.
func (ds DomainRFCs) FilterInvalidWithErrors(allowIPs bool) (valid DomainRFCs, invalid map[DomainRFC]error) {
	invalid = make(map[DomainRFC]error)
	for _, domain := range ds {
		if ok, err := domain.IsValidWithOptions(allowIPs); ok {
			valid = append(valid, domain)
		} else {
			invalid[domain] = err
		}
	}
	return valid, invalid
}
