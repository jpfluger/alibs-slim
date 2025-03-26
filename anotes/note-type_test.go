package anotes

import (
	"testing"
)

// TestIsEmpty checks if the IsEmpty method correctly identifies empty NoteTypes.
func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		noteType NoteType
		want     bool
	}{
		{"EmptyString", "", true},
		{"WhitespaceOnly", "   ", true},
		{"NonEmptyText", NOTETYPE_TEXT, false},
		{"NonEmptyFlag", NOTETYPE_FLAG, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.noteType.IsEmpty(); got != tt.want {
				t.Errorf("NoteType.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTrimSpace checks if the TrimSpace method correctly trims the NoteType.
func TestTrimSpace(t *testing.T) {
	tests := []struct {
		name     string
		noteType NoteType
		want     NoteType
	}{
		{"NoTrimNeeded", NOTETYPE_TEXT, "text"},
		{"TrimSpaces", "  text  ", "text"},
		{"TrimFlagType", "  flag  ", "flag"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.noteType.TrimSpace(); got != tt.want {
				t.Errorf("NoteType.TrimSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestString checks if the String method correctly converts NoteType to a string.
func TestString(t *testing.T) {
	tests := []struct {
		name     string
		noteType NoteType
		want     string
	}{
		{"Empty", "", ""},
		{"TextType", NOTETYPE_TEXT, "text"},
		{"FlagType", NOTETYPE_FLAG, "flag"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.noteType.String(); got != tt.want {
				t.Errorf("NoteType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToStringTrimLower checks if ToStringTrimLower method correctly processes the NoteType.
func TestToStringTrimLower(t *testing.T) {
	tests := []struct {
		name     string
		noteType NoteType
		want     string
	}{
		{"AllCaps", "TEXT", "text"},
		{"MixedCase", "TeXt", "text"},
		{"WithSpaces", "  TeXt  ", "text"},
		{"ImageFlagType", "  IMAGE_FLAG  ", "image_flag"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.noteType.ToStringTrimLower(); got != tt.want {
				t.Errorf("NoteType.ToStringTrimLower() = %v, want %v", got, tt.want)
			}
		})
	}
}
