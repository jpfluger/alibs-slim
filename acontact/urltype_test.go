package acontact

import (
	"testing"
)

// TestUrlTypeIsEmpty checks if the IsEmpty method correctly identifies empty UrlTypes
func TestUrlTypeIsEmpty(t *testing.T) {
	tests := []struct {
		ut   UrlType
		want bool
	}{
		{URLTYPE_HOME, false},
		{URLTYPE_PERSONAL, false},
		{URLTYPE_BUSINESS, false},
		{UrlType(" "), true},
		{UrlType(""), true},
	}

	for _, tt := range tests {
		if got := tt.ut.IsEmpty(); got != tt.want {
			t.Errorf("UrlType.IsEmpty() = %v, want %v", got, tt.want)
		}
	}
}

// TestUrlTypeTrimSpace checks if the TrimSpace method correctly trims UrlTypes
func TestUrlTypeTrimSpace(t *testing.T) {
	tests := []struct {
		ut   UrlType
		want UrlType
	}{
		{UrlType(" home "), URLTYPE_HOME},
		{UrlType(" personal "), URLTYPE_PERSONAL},
		{UrlType(" business "), URLTYPE_BUSINESS},
		{UrlType(" "), UrlType("")},
	}

	for _, tt := range tests {
		if got := tt.ut.TrimSpace(); got != tt.want {
			t.Errorf("UrlType.TrimSpace() = %v, want %v", got, tt.want)
		}
	}
}

// TestUrlTypeString checks if the String method correctly converts UrlTypes to strings
func TestUrlTypeString(t *testing.T) {
	if URLTYPE_HOME.String() != "home" {
		t.Errorf("UrlType.String() failed, got %v, want %v", URLTYPE_HOME.String(), "home")
	}
}

// TestUrlTypeToStringTrimLower checks if the ToStringTrimLower method correctly converts UrlTypes to trimmed, lowercase strings
func TestUrlTypeToStringTrimLower(t *testing.T) {
	if URLTYPE_HOME.ToStringTrimLower() != "home" {
		t.Errorf("UrlType.ToStringTrimLower() failed, got %v, want %v", URLTYPE_HOME.ToStringTrimLower(), "home")
	}
}
