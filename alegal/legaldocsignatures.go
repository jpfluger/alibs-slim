package alegal

import (
	"time"
)

// LegalDocSignatures is a slice of pointers to LegalDocSignature.
type LegalDocSignatures []*LegalDocSignature

// Find searches for a LegalDocSignature by its key within the slice.
func (lds LegalDocSignatures) Find(key LegalDocSignatureKey) *LegalDocSignature {
	if lds == nil {
		return nil // Return nil if the slice is nil
	}
	for _, ld := range lds {
		if ld.Key == key {
			return ld // Return the document if the key matches
		}
	}
	return nil // Return nil if no matching document is found
}

// HasAcceptDate checks if a document with the given key has an accept date after the required date.
func (lds LegalDocSignatures) HasAcceptDate(key LegalDocSignatureKey, requiredAfterDate time.Time) bool {
	ld := lds.Find(key)
	return ld != nil && ld.IsAfterDate(ld.AcceptDate, requiredAfterDate)
}

// HasRejectDate checks if a document with the given key has a reject date after the required date.
func (lds LegalDocSignatures) HasRejectDate(key LegalDocSignatureKey, requiredAfterDate time.Time) bool {
	ld := lds.Find(key)
	return ld != nil && ld.IsAfterDate(ld.RejectDate, requiredAfterDate)
}

// HasEffectiveDate checks if a document with the given key has an effective date after the required date.
func (lds LegalDocSignatures) HasEffectiveDate(key LegalDocSignatureKey, requiredAfterDate time.Time) bool {
	ld := lds.Find(key)
	return ld != nil && ld.IsAfterDate(ld.EffectiveDate, requiredAfterDate)
}

// IsAccepted checks if a document with the given key has been accepted after the required date.
func (lds LegalDocSignatures) IsAccepted(key LegalDocSignatureKey, requiredAfterDate time.Time) bool {
	ld := lds.Find(key)
	return ld != nil && ld.IsAfterDate(ld.AcceptDate, requiredAfterDate)
}

// IsAcceptedAndEffective checks if a document with the given key has been accepted and is effective after the required date.
func (lds LegalDocSignatures) IsAcceptedAndEffective(key LegalDocSignatureKey, requiredAfterDate time.Time) bool {
	ld := lds.Find(key)
	return ld != nil && ld.IsAfterDate(ld.AcceptDate, requiredAfterDate) && ld.IsAfterDate(ld.EffectiveDate, requiredAfterDate)
}

// HasTerminationDate checks if a document with the given key has a termination date set.
func (lds LegalDocSignatures) HasTerminationDate(key LegalDocSignatureKey) bool {
	ld := lds.Find(key)
	return ld != nil && ld.TerminationDate != nil && !ld.TerminationDate.IsZero()
}
