package anotes

import (
	"testing"
	"time"
)

// TestFixIntegrity checks the FixIntegrity method for various scenarios.
func TestFixIntegrity(t *testing.T) {
	tests := []struct {
		name     string
		noteText *NoteText
		want     bool
	}{
		{"NilNote", nil, false},
		{"EmptyType", &NoteText{Date: time.Now()}, true},
		{"ValidNote", &NoteText{Type: NOTETYPE_TEXT, Date: time.Now()}, true},
		{"ZeroDate", &NoteText{Type: NOTETYPE_TEXT}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.noteText.FixIntegrity(); got != tt.want {
				t.Errorf("NoteText.FixIntegrity() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidate checks the Validate method for various scenarios.
func TestValidate(t *testing.T) {
	validDate := time.Now()
	tests := []struct {
		name     string
		noteText *NoteText
		wantErr  bool
	}{
		{"ValidNote", &NoteText{Type: NOTETYPE_TEXT, Date: validDate, Text: "Sample text", UserId: "user1"}, false},
		{"EmptyType", &NoteText{Date: validDate, Text: "Sample text", UserId: "user1"}, true},
		{"ZeroDate", &NoteText{Type: NOTETYPE_TEXT, Text: "Sample text", UserId: "user1"}, true},
		{"EmptyText", &NoteText{Type: NOTETYPE_TEXT, Date: validDate, UserId: "user1"}, true},
		{"EmptyUserId", &NoteText{Type: NOTETYPE_TEXT, Date: validDate, Text: "Sample text"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.noteText.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("NoteText.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetMethods checks the getter methods for the NoteText struct.
func TestGetMethods(t *testing.T) {
	note := NoteText{
		Type:   NOTETYPE_TEXT,
		Date:   time.Now(),
		Text:   "Sample text",
		UserId: "user1",
	}

	if note.GetType() != NOTETYPE_TEXT {
		t.Errorf("NoteText.GetType() = %v, want %v", note.GetType(), NOTETYPE_TEXT)
	}
	if note.GetDate().IsZero() {
		t.Error("NoteText.GetDate() is zero, want non-zero")
	}
	if note.GetText() != "Sample text" {
		t.Errorf("NoteText.GetText() = %v, want 'Sample text'", note.GetText())
	}
	if note.GetUserId() != "user1" {
		t.Errorf("NoteText.GetUserId() = %v, want 'user1'", note.GetUserId())
	}
}

// TestNoteTextsMethods checks the methods of the NoteTexts type.
func TestNoteTextsMethods(t *testing.T) {
	note1 := &NoteText{Type: NOTETYPE_TEXT, Date: time.Now(), Text: "Note 1", UserId: "user1"}
	note2 := &NoteText{Type: NOTETYPE_TEXT, Date: time.Now(), Text: "Note 2", UserId: "user2"}

	notes := NoteTexts{note1}

	// Test Find
	if found := notes.FindByType(NOTETYPE_TEXT); found == nil || found.Text != "Note 1" {
		t.Errorf("NoteTexts.Find() did not find the correct note")
	}

	// Test Has
	if !notes.HasByType(NOTETYPE_TEXT) {
		t.Errorf("NoteTexts.Has() should return true for existing note type")
	}

	// Test SetByType
	notes.SetByType(note2)
	if len(notes) != 1 {
		t.Errorf("NoteTexts.Set() should add a note when not present")
	}

	// Test Remove
	notes.RemoveByType(NOTETYPE_TEXT)
	if len(notes) != 0 {
		t.Errorf("NoteTexts.Remove() should remove notes of the specified type")
	}

	// Test Validate
	if err := notes.Validate(); err != nil {
		t.Errorf("NoteTexts.Validate() should not return an error for valid notes")
	}
}
