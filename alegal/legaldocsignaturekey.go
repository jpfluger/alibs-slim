package alegal

import (
	"strings"
)

// LEGALKEY_TERMS is a constant key for terms in legal documents.
const LEGALKEY_TERMS = LegalDocSignatureKey("terms")

// LegalDocSignatureKey represents a key used to identify legal document signatures.
type LegalDocSignatureKey string

// IsEmpty checks if the LegalDocSignatureKey is empty after trimming whitespace.
func (ldsk LegalDocSignatureKey) IsEmpty() bool {
	return strings.TrimSpace(string(ldsk)) == ""
}

// TrimSpace trims whitespace from the LegalDocSignatureKey and returns a new LegalDocSignatureKey.
func (ldsk LegalDocSignatureKey) TrimSpace() LegalDocSignatureKey {
	return LegalDocSignatureKey(strings.TrimSpace(string(ldsk)))
}

// String returns the LegalDocSignatureKey as a trimmed string.
func (ldsk LegalDocSignatureKey) String() string {
	return strings.TrimSpace(string(ldsk))
}

// ToStringTrimLower trims whitespace from the LegalDocSignatureKey, converts it to lowercase, and returns it as a string.
func (ldsk LegalDocSignatureKey) ToStringTrimLower() string {
	return strings.ToLower(ldsk.String())
}

// LegalDocSignatureKeys is a slice of LegalDocSignatureKey, representing multiple keys.
type LegalDocSignatureKeys []LegalDocSignatureKey
