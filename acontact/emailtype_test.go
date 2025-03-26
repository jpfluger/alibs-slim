package acontact

import (
	"testing"
)

// TestEmailType_IsEmpty tests the IsEmpty method of the EmailType type
func TestEmailType_IsEmpty(t *testing.T) {
	tests := []struct {
		et   EmailType
		want bool
	}{
		{EMAILTYPE_NO_REPLY, false},
		{EmailType(" "), true},
		{EmailType(""), true},
	}

	for _, tt := range tests {
		if got := tt.et.IsEmpty(); got != tt.want {
			t.Errorf("EmailType.IsEmpty() = %v, want %v", got, tt.want)
		}
	}
}

// TestEmailType_TrimSpace tests the TrimSpace method of the EmailType type
func TestEmailType_TrimSpace(t *testing.T) {
	tests := []struct {
		et   EmailType
		want EmailType
	}{
		{EmailType(" noreply "), EMAILTYPE_NO_REPLY},
		{EmailType(" "), EmailType("")},
	}

	for _, tt := range tests {
		if got := tt.et.TrimSpace(); got != tt.want {
			t.Errorf("EmailType.TrimSpace() = %v, want %v", got, tt.want)
		}
	}
}

// TestEmailType_String tests the String method of the EmailType type
func TestEmailType_String(t *testing.T) {
	if EMAILTYPE_NO_REPLY.String() != "noreply" {
		t.Errorf("EmailType.String() failed, got %v, want %v", EMAILTYPE_NO_REPLY.String(), "noreply")
	}
}

// TestEmailType_ToStringTrimLower tests the ToStringTrimLower method of the EmailType type
func TestEmailType_ToStringTrimLower(t *testing.T) {
	if EMAILTYPE_NO_REPLY.ToStringTrimLower() != "noreply" {
		t.Errorf("EmailType.ToStringTrimLower() failed, got %v, want %v", EMAILTYPE_NO_REPLY.ToStringTrimLower(), "noreply")
	}
}

// TestEmailType_GetType tests the GetType method of the EmailType type
func TestEmailType_GetType(t *testing.T) {
	et := EmailType("work:primary")
	if et.GetType() != "work" {
		t.Errorf("EmailType.GetType() failed, got %v, want %v", et.GetType(), "work")
	}
}

// TestEmailType_GetPart tests the GetPart method of the EmailType type
func TestEmailType_GetPart(t *testing.T) {
	et := EmailType("home:secondary")
	if et.GetPart() != "secondary" {
		t.Errorf("EmailType.GetPart() failed, got %v, want %v", et.GetPart(), "secondary")
	}
}
