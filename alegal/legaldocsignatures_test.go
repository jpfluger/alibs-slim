package alegal

import (
	"testing"
	"time"
)

// TestFind checks if the Find method correctly finds a LegalDocSignature by key.
func TestFind(t *testing.T) {
	// Create a slice of LegalDocSignatures with different keys
	signatures := LegalDocSignatures{
		{Key: LEGALKEY_TERMS},
		{Key: "privacy"},
	}

	// Test finding an existing key
	found := signatures.Find(LEGALKEY_TERMS)
	if found == nil || found.Key != LEGALKEY_TERMS {
		t.Errorf("Find() - expected to find key '%s', got '%v'", LEGALKEY_TERMS, found)
	}

	// Test finding a non-existing key
	notFound := signatures.Find("nonexistent")
	if notFound != nil {
		t.Errorf("Find() - expected to not find key 'nonexistent', got '%v'", notFound)
	}
}

// TestHasAcceptDate checks if the HasAcceptDate method correctly identifies documents with an accept date after the required date.
func TestHasAcceptDate(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	signature := &LegalDocSignature{Key: LEGALKEY_TERMS, AcceptDate: &now}
	signatures := LegalDocSignatures{signature}

	if !signatures.HasAcceptDate(LEGALKEY_TERMS, past) {
		t.Errorf("HasAcceptDate() - expected true, got false")
	}
}

// TestHasRejectDate checks if the HasRejectDate method correctly identifies documents with a reject date after the required date.
func TestHasRejectDate(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	signature := &LegalDocSignature{Key: LEGALKEY_TERMS, RejectDate: &now}
	signatures := LegalDocSignatures{signature}

	if !signatures.HasRejectDate(LEGALKEY_TERMS, past) {
		t.Errorf("HasRejectDate() - expected true, got false")
	}
}

// Add more tests as needed...

// TestHasEffectiveDate checks if the HasEffectiveDate method correctly identifies documents with an effective date after the required date.
func TestHasEffectiveDate(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	signature := &LegalDocSignature{Key: LEGALKEY_TERMS, EffectiveDate: &now}
	signatures := LegalDocSignatures{signature}

	if !signatures.HasEffectiveDate(LEGALKEY_TERMS, past) {
		t.Errorf("HasEffectiveDate() - expected true, got false")
	}
}

// TestIsAccepted checks if the IsAccepted method correctly identifies documents that have been accepted after the required date.
func TestIsAccepted(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	signature := &LegalDocSignature{Key: LEGALKEY_TERMS, AcceptDate: &now}
	signatures := LegalDocSignatures{signature}

	if !signatures.IsAccepted(LEGALKEY_TERMS, past) {
		t.Errorf("IsAccepted() - expected true, got false")
	}
}

// TestIsAcceptedAndEffective checks if the IsAcceptedAndEffective method correctly identifies documents that have been accepted and are effective after the required date.
func TestIsAcceptedAndEffective(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	signature := &LegalDocSignature{Key: LEGALKEY_TERMS, AcceptDate: &now, EffectiveDate: &now}
	signatures := LegalDocSignatures{signature}

	if !signatures.IsAcceptedAndEffective(LEGALKEY_TERMS, past) {
		t.Errorf("IsAcceptedAndEffective() - expected true, got false")
	}
}

// TestHasTerminationDate checks if the HasTerminationDate method correctly identifies documents with a termination date set.
func TestHasTerminationDate(t *testing.T) {
	now := time.Now()
	signature := &LegalDocSignature{Key: LEGALKEY_TERMS, TerminationDate: &now}
	signatures := LegalDocSignatures{signature}

	if !signatures.HasTerminationDate(LEGALKEY_TERMS) {
		t.Errorf("HasTerminationDate() - expected true, got false")
	}
}
