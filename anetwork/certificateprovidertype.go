package anetwork

import (
	"strings"
)

type CertificateProviderType string

// IsEmpty checks if the CertificateProviderType is empty after trimming spaces.
func (cpt CertificateProviderType) IsEmpty() bool {
	// Trim spaces from the CertificateProviderType and check if the result is an empty string.
	return strings.TrimSpace(string(cpt)) == ""
}

// TrimSpace trims spaces from the CertificateProviderType and returns a new CertificateProviderType.
func (cpt CertificateProviderType) TrimSpace() CertificateProviderType {
	// Trim spaces from the CertificateProviderType and return the result as a new CertificateProviderType.
	return CertificateProviderType(strings.TrimSpace(string(cpt)))
}

// String converts CertificateProviderType to a string.
func (cpt CertificateProviderType) String() string {
	// Convert the CertificateProviderType to a string and return it.
	return string(cpt)
}

// ToStringTrimLower converts CertificateProviderType to a string, trims spaces, and makes it lowercase.
func (cpt CertificateProviderType) ToStringTrimLower() string {
	// Convert the CertificateProviderType to a string, trim spaces, convert to lowercase, and return the result.
	return strings.ToLower(cpt.TrimSpace().String())
}

// CertificateProviderTypes defines a slice of CertificateProviderType.
type CertificateProviderTypes []CertificateProviderType

// Contains checks if the CertificateProviderTypes slice contains a specific CertificateProviderType.
func (cpts CertificateProviderTypes) Contains(cpt CertificateProviderType) bool {
	// Iterate over the CertificateProviderTypes slice.
	for _, t := range cpts {
		// Check if the current CertificateProviderType matches the specified CertificateProviderType.
		if t == cpt {
			return true // Return true if a match is found.
		}
	}
	return false // Return false if no match is found.
}

// Add appends a new CertificateProviderType to the CertificateProviderTypes slice if it's not already present.
func (cpts *CertificateProviderTypes) Add(cpt CertificateProviderType) {
	// Check if the CertificateProviderType is not already in the slice.
	if !cpts.Contains(cpt) {
		*cpts = append(*cpts, cpt) // Append the new CertificateProviderType to the slice.
	}
}

// Remove deletes a CertificateProviderType from the CertificateProviderTypes slice.
func (cpts *CertificateProviderTypes) Remove(cpt CertificateProviderType) {
	// Create a new slice to store the result.
	var result CertificateProviderTypes
	// Iterate over the CertificateProviderTypes slice.
	for _, t := range *cpts {
		// If the current CertificateProviderType does not match the specified CertificateProviderType, add it to the result slice.
		if t != cpt {
			result = append(result, t)
		}
	}
	*cpts = result // Set the CertificateProviderTypes slice to the result.
}
