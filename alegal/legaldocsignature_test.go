package alegal

import (
	"testing"
	"time"
)

func TestLegalDocSignature(t *testing.T) {
	ldoc := &LegalDocSignature{
		Key:             LEGALKEY_TERMS,
		AcceptDate:      nil,
		RejectDate:      nil,
		EffectiveDate:   nil,
		TerminationDate: nil,
		History:         nil,
	}

	if ldoc.IsValidAcceptDate() != false {
		t.Errorf("Expected IsValidAcceptDate to be false, got true")
	}
	if ldoc.IsValidEffectiveDate() != false {
		t.Errorf("Expected IsValidEffectiveDate to be false, got true")
	}
	if ldoc.IsValidRejectDate() != false {
		t.Errorf("Expected IsValidRejectDate to be false, got true")
	}
	if ldoc.IsValidTerminationDate() != false {
		t.Errorf("Expected IsValidTerminationDate to be false, got true")
	}
}

// TestIsValidAcceptDate tests the IsValidAcceptDate method.
func TestIsValidAcceptDate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		acceptDate *time.Time
		expected   bool
	}{
		{&now, true},
		{nil, false},
		{&time.Time{}, false},
	}

	for _, test := range tests {
		lda := LegalDocSignature{AcceptDate: test.acceptDate}
		if lda.IsValidAcceptDate() != test.expected {
			t.Errorf("IsValidAcceptDate() - expected %v, got %v", test.expected, !test.expected)
		}
	}
}

// TestIsValidRejectDate tests the IsValidRejectDate method.
func TestIsValidRejectDate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		rejectDate *time.Time
		expected   bool
	}{
		{&now, true},
		{nil, false},
		{&time.Time{}, false},
	}

	for _, test := range tests {
		lda := LegalDocSignature{RejectDate: test.rejectDate}
		if lda.IsValidRejectDate() != test.expected {
			t.Errorf("IsValidRejectDate() - expected %v, got %v", test.expected, !test.expected)
		}
	}
}

// TestIsValidEffectiveDate tests the IsValidEffectiveDate method.
func TestIsValidEffectiveDate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		effectiveDate *time.Time
		expected      bool
	}{
		{&now, true},
		{nil, false},
		{&time.Time{}, false},
	}

	for _, test := range tests {
		lda := LegalDocSignature{EffectiveDate: test.effectiveDate}
		if lda.IsValidEffectiveDate() != test.expected {
			t.Errorf("IsValidEffectiveDate() - expected %v, got %v", test.expected, !test.expected)
		}
	}
}

// TestIsValidTerminationDate tests the IsValidTerminationDate method.
func TestIsValidTerminationDate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		terminationDate *time.Time
		expected        bool
	}{
		{&now, true},
		{nil, false},
		{&time.Time{}, false},
	}

	for _, test := range tests {
		lda := LegalDocSignature{TerminationDate: test.terminationDate}
		if lda.IsValidTerminationDate() != test.expected {
			t.Errorf("IsValidTerminationDate() - expected %v, got %v", test.expected, !test.expected)
		}
	}
}

// TestIsAfterDate tests the IsAfterDate method.
func TestIsAfterDate(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tests := []struct {
		target            *time.Time
		requiredAfterDate time.Time
		expected          bool
	}{
		{&future, now, true},
		{&past, now, false},
		{nil, now, false},
		{&time.Time{}, now, false},
		{&now, time.Time{}, true},
	}

	for _, test := range tests {
		ld := LegalDocSignature{}
		if ld.IsAfterDate(test.target, test.requiredAfterDate) != test.expected {
			t.Errorf("IsAfterDate() - expected %v, got %v", test.expected, !test.expected)
		}
	}
}
