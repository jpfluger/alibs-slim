package alegal

import (
	"testing"
)

// TestIsEmpty checks if the IsEmpty method correctly identifies empty LegalDocSignatureKeys.
func TestIsEmpty(t *testing.T) {
	tests := []struct {
		key      LegalDocSignatureKey
		expected bool
	}{
		{"", true},
		{" ", true},
		{"terms", false},
	}

	for _, test := range tests {
		if test.key.IsEmpty() != test.expected {
			t.Errorf("IsEmpty() for key '%s' - expected %v, got %v", test.key, test.expected, !test.expected)
		}
	}
}

// TestTrimSpace checks if the TrimSpace method correctly trims whitespace from LegalDocSignatureKeys.
func TestTrimSpace(t *testing.T) {
	tests := []struct {
		key      LegalDocSignatureKey
		expected LegalDocSignatureKey
	}{
		{" terms ", "terms"},
		{"  terms", "terms"},
		{"terms  ", "terms"},
	}

	for _, test := range tests {
		if test.key.TrimSpace() != test.expected {
			t.Errorf("TrimSpace() for key '%s' - expected '%s', got '%s'", test.key, test.expected, test.key.TrimSpace())
		}
	}
}

// TestString checks if the String method correctly returns a trimmed string representation of LegalDocSignatureKeys.
func TestString(t *testing.T) {
	tests := []struct {
		key      LegalDocSignatureKey
		expected string
	}{
		{" terms ", "terms"},
		{"  terms", "terms"},
		{"terms  ", "terms"},
	}

	for _, test := range tests {
		if test.key.String() != test.expected {
			t.Errorf("String() for key '%s' - expected '%s', got '%s'", test.key, test.expected, test.key.String())
		}
	}
}

// TestToStringTrimLower checks if the ToStringTrimLower method correctly returns a lowercase, trimmed string.
func TestToStringTrimLower(t *testing.T) {
	tests := []struct {
		key      LegalDocSignatureKey
		expected string
	}{
		{" TERMS ", "terms"},
		{"  Terms", "terms"},
		{"terms  ", "terms"},
	}

	for _, test := range tests {
		if test.key.ToStringTrimLower() != test.expected {
			t.Errorf("ToStringTrimLower() for key '%s' - expected '%s', got '%s'", test.key, test.expected, test.key.ToStringTrimLower())
		}
	}
}

// Add more tests as needed...
