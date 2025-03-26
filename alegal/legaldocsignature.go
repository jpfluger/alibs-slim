package alegal

import (
	"time"
)

// LegalDocSignature represents the signature information for a legal document.
type LegalDocSignature struct {
	// Key uniquely identifies the legal document.
	Key LegalDocSignatureKey `json:"key"`

	// AcceptDate is the date when the document was accepted.
	AcceptDate *time.Time `json:"acceptDate,omitempty"`

	// RejectDate is the date when the document was rejected.
	RejectDate *time.Time `json:"rejectDate,omitempty"`

	// EffectiveDate is the date when the document becomes effective.
	EffectiveDate *time.Time `json:"effectiveDate,omitempty"`

	// TerminationDate is the date when the document was terminated, if applicable.
	TerminationDate *time.Time `json:"terminationDate,omitempty"`

	// History captures the history of previously signed documents, if applicable.
	History LegalDocSignatures `json:"history,omitempty"`
}

// IsValidAcceptDate checks if the AcceptDate is valid (non-nil and non-zero).
func (lda *LegalDocSignature) IsValidAcceptDate() bool {
	return lda.AcceptDate != nil && !lda.AcceptDate.IsZero()
}

// IsValidRejectDate checks if the RejectDate is valid (non-nil and non-zero).
func (lda *LegalDocSignature) IsValidRejectDate() bool {
	return lda.RejectDate != nil && !lda.RejectDate.IsZero()
}

// IsValidEffectiveDate checks if the EffectiveDate is valid (non-nil and non-zero).
func (lda *LegalDocSignature) IsValidEffectiveDate() bool {
	return lda.EffectiveDate != nil && !lda.EffectiveDate.IsZero()
}

// IsValidTerminationDate checks if the TerminationDate is valid (non-nil and non-zero).
func (lda *LegalDocSignature) IsValidTerminationDate() bool {
	return lda.TerminationDate != nil && !lda.TerminationDate.IsZero()
}

// IsAfterDate checks if the target date is after the requiredAfterDate.
func (ld *LegalDocSignature) IsAfterDate(target *time.Time, requiredAfterDate time.Time) bool {
	if target == nil || target.IsZero() {
		return false
	}
	return target.After(requiredAfterDate)
}
