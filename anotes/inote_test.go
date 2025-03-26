package anotes

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/areflect"
	"reflect"
	"testing"
	"time"
)

// MockNoteText is a concrete type that implements the INote interface for testing.
type MockNoteText struct {
	Type   NoteType
	Date   time.Time
	Text   string
	UserId string
}

func (m MockNoteText) GetType() NoteType  { return m.Type }
func (m MockNoteText) GetDate() time.Time { return m.Date }
func (m MockNoteText) GetText() string    { return m.Text }
func (m MockNoteText) GetUserId() string  { return m.UserId }

// init function registers the note types with the type manager upon package initialization.
func init() {
	// Mocking the registration of types with the type manager.
	_ = areflect.TypeManager().Register(TYPEMANAGER_NOTE, "anotes-mock", returnTypeManagerINoteMock)
}

// returnTypeManagerINoteMock is a mock function that returns the reflect.Type corresponding to the provided typeName.
func returnTypeManagerINoteMock(typeName string) (reflect.Type, error) {
	switch NoteType(typeName) {
	case NOTETYPE_TEXT:
		return reflect.TypeOf(MockNoteText{}), nil
	// Add cases for other types as needed.
	default:
		return nil, fmt.Errorf("type not found")
	}
}

// TestINotes_UnmarshalJSON tests the UnmarshalJSON method of INotes.
func TestINotes_UnmarshalJSON(t *testing.T) {
	jsonData := `[{"type":"text","date":"2020-01-01T00:00:00Z","text":"Sample text","userId":"user1"}]`
	var notes INotes
	err := json.Unmarshal([]byte(jsonData), &notes)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if len(notes) != 1 {
		t.Fatalf("Expected 1 note, got %d", len(notes))
	}

	note := notes[0]
	if note.GetType() != "text" || note.GetText() != "Sample text" {
		t.Errorf("Unmarshaled note does not match expected values")
	}
}

// TestINotes_Methods tests the Find, Has, Set, and Remove methods of INotes.
func TestINotes_Methods(t *testing.T) {
	note1 := MockNoteText{Type: "text", Date: time.Now(), Text: "Note 1", UserId: "user1"}
	note2 := MockNoteText{Type: "text", Date: time.Now(), Text: "Note 2", UserId: "user2"}

	notes := INotes{note1}

	// Test Find
	found := notes.FindByType("text")
	if found == nil || found.GetText() != "Note 1" {
		t.Errorf("Find did not return the correct note")
	}

	// Test Has
	if !notes.HasByType("text") {
		t.Errorf("Has should return true for existing note type")
	}

	// Test Set
	notes.SetByType(note2)
	if len(notes) != 1 || notes[0].GetText() != "Note 2" {
		t.Errorf("Set did not add the new note correctly")
	}

	// Test Remove
	notes.Remove("text")
	if len(notes) != 0 {
		t.Errorf("Remove did not remove the note correctly")
	}
}
