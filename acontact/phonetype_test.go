package acontact

import (
	"testing"
)

// TestPhoneTypeIsEmpty checks if the IsEmpty method correctly identifies empty PhoneTypes
func TestPhoneTypeIsEmpty(t *testing.T) {
	tests := []struct {
		pt   PhoneType
		want bool
	}{
		{PHONETYPE_MOBILE, false},
		{PHONETYPE_HOME, false},
		{PHONETYPE_WORK, false},
		{PHONETYPE_DIRECT, false},
		{PhoneType(" "), true},
		{PhoneType(""), true},
	}

	for _, tt := range tests {
		if got := tt.pt.IsEmpty(); got != tt.want {
			t.Errorf("PhoneType.IsEmpty() = %v, want %v", got, tt.want)
		}
	}
}

// TestPhoneTypeTrimSpace checks if the TrimSpace method correctly trims PhoneTypes
func TestPhoneTypeTrimSpace(t *testing.T) {
	tests := []struct {
		pt   PhoneType
		want PhoneType
	}{
		{PhoneType(" mobile "), PHONETYPE_MOBILE},
		{PhoneType(" home "), PHONETYPE_HOME},
		{PhoneType(" work "), PHONETYPE_WORK},
		{PhoneType(" direct "), PHONETYPE_DIRECT},
		{PhoneType(" "), PhoneType("")},
	}

	for _, tt := range tests {
		if got := tt.pt.TrimSpace(); got != tt.want {
			t.Errorf("PhoneType.TrimSpace() = %v, want %v", got, tt.want)
		}
	}
}

// TestPhoneTypeString checks if the String method correctly converts PhoneTypes to strings
func TestPhoneTypeString(t *testing.T) {
	if PHONETYPE_MOBILE.String() != "mobile" {
		t.Errorf("PhoneType.String() failed, got %v, want %v", PHONETYPE_MOBILE.String(), "mobile")
	}
}

// TestPhoneTypeToStringTrimLower checks if the ToStringTrimLower method correctly converts PhoneTypes to trimmed, lowercase strings
func TestPhoneTypeToStringTrimLower(t *testing.T) {
	if PHONETYPE_MOBILE.ToStringTrimLower() != "mobile" {
		t.Errorf("PhoneType.ToStringTrimLower() failed, got %v, want %v", PHONETYPE_MOBILE.ToStringTrimLower(), "mobile")
	}
}
