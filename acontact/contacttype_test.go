package acontact

import (
	"testing"
)

// TestIsEmpty tests the IsEmpty method of ContactType
func TestIsEmpty(t *testing.T) {
	tests := []struct {
		ct   ContactType
		want bool
	}{
		{CONTACTTYPE_PERSON, false},
		{CONTACTTYPE_ENTITY, false},
		{ContactType(" "), true},
		{ContactType(""), true},
	}

	for _, tt := range tests {
		if got := tt.ct.IsEmpty(); got != tt.want {
			t.Errorf("ContactType.IsEmpty() = %v, want %v", got, tt.want)
		}
	}
}

// TestTrimSpace tests the TrimSpace method of ContactType
func TestTrimSpace(t *testing.T) {
	tests := []struct {
		ct   ContactType
		want ContactType
	}{
		{ContactType(" person "), CONTACTTYPE_PERSON},
		{ContactType(" entity "), CONTACTTYPE_ENTITY},
		{ContactType(" "), ContactType("")},
	}

	for _, tt := range tests {
		if got := tt.ct.TrimSpace(); got != tt.want {
			t.Errorf("ContactType.TrimSpace() = %v, want %v", got, tt.want)
		}
	}
}

// TestHasMatch tests the HasMatch method of ContactType
func TestHasMatch(t *testing.T) {
	if !CONTACTTYPE_PERSON.HasMatch(ContactType("person")) {
		t.Errorf("ContactType.HasMatch() failed to recognize matching ContactTypes")
	}
	if CONTACTTYPE_PERSON.HasMatch(CONTACTTYPE_ENTITY) {
		t.Errorf("ContactType.HasMatch() should not match different ContactTypes")
	}
}

// TestString tests the String method of ContactType
func TestString(t *testing.T) {
	if CONTACTTYPE_PERSON.String() != "person" {
		t.Errorf("ContactType.String() failed, got %v, want %v", CONTACTTYPE_PERSON.String(), "person")
	}
}

// TestToStringTrimLower tests the ToStringTrimLower method of ContactType
func TestToStringTrimLower(t *testing.T) {
	if CONTACTTYPE_PERSON.ToStringTrimLower() != "person" {
		t.Errorf("ContactType.ToStringTrimLower() failed, got %v, want %v", CONTACTTYPE_PERSON.ToStringTrimLower(), "person")
	}
}
