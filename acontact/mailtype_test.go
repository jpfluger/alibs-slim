package acontact

import (
	"testing"
)

// TestMailType_IsEmpty tests the IsEmpty method of the MailType type
func TestMailType_IsEmpty(t *testing.T) {
	tests := []struct {
		mt   MailType
		want bool
	}{
		{MAILTYPE_HOME, false},
		{MailType(" "), true},
		{MailType(""), true},
	}

	for _, tt := range tests {
		if got := tt.mt.IsEmpty(); got != tt.want {
			t.Errorf("MailType.IsEmpty() = %v, want %v", got, tt.want)
		}
	}
}

// TestMailType_TrimSpace tests the TrimSpace method of the MailType type
func TestMailType_TrimSpace(t *testing.T) {
	tests := []struct {
		mt   MailType
		want MailType
	}{
		{MailType(" home "), MAILTYPE_HOME},
		{MailType(" "), MailType("")},
	}

	for _, tt := range tests {
		if got := tt.mt.TrimSpace(); got != tt.want {
			t.Errorf("MailType.TrimSpace() = %v, want %v", got, tt.want)
		}
	}
}

// TestMailType_String tests the String method of the MailType type
func TestMailType_String(t *testing.T) {
	if MAILTYPE_HOME.String() != "home" {
		t.Errorf("MailType.String() failed, got %v, want %v", MAILTYPE_HOME.String(), "home")
	}
}

// TestMailType_ToStringTrimLower tests the ToStringTrimLower method of the MailType type
func TestMailType_ToStringTrimLower(t *testing.T) {
	if MAILTYPE_HOME.ToStringTrimLower() != "home" {
		t.Errorf("MailType.ToStringTrimLower() failed, got %v, want %v", MAILTYPE_HOME.ToStringTrimLower(), "home")
	}
}
