package anotes

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// NOTETYPE_TEXT is a constant representing a text note type.
const NOTETYPE_TEXT NoteType = "text"

// NoteText struct represents a note with text content.
type NoteText struct {
	Type   NoteType  `json:"type,omitempty"`   // The type of the note.
	Date   time.Time `json:"date,omitempty"`   // The date associated with the note.
	Text   string    `json:"text,omitempty"`   // The text content of the note.
	UserId string    `json:"userId,omitempty"` // The ID of the user who created the note.

	mu sync.RWMutex // Mutex to protect concurrent access to the note.
}

// FixIntegrity checks for valid Type and Date, while Text and UserId are optional.
func (n *NoteText) FixIntegrity() bool {
	if n == nil {
		return false
	}
	if n.Type.IsEmpty() {
		n.Type = NOTETYPE_TEXT
	}
	if err := n.ValidateWithOptions(true, false, false); err != nil {
		return false
	}
	return true
}

// Validate checks the note for validity.
func (n *NoteText) Validate() error {
	return n.ValidateWithOptions(true, true, true)
}

// ValidateWithOptions performs validation checks on the note based on the provided options.
func (n *NoteText) ValidateWithOptions(checkType, checkText, checkUserId bool) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if checkType && n.Type.IsEmpty() {
		return fmt.Errorf("type is empty")
	}
	if n.Date.IsZero() {
		return fmt.Errorf("date is zero")
	}
	n.Text = strings.TrimSpace(n.Text)
	if checkText && n.Text == "" {
		return fmt.Errorf("text is empty")
	}
	n.UserId = strings.TrimSpace(n.UserId)
	if checkUserId && n.UserId == "" {
		return fmt.Errorf("userId is empty")
	}
	return nil
}

// GetType safely retrieves the note's type.
func (n *NoteText) GetType() NoteType {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Type
}

// GetDate safely retrieves the note's date.
func (n *NoteText) GetDate() time.Time {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Date
}

// GetText safely retrieves the note's text.
func (n *NoteText) GetText() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Text
}

// GetUserId safely retrieves the note's user ID.
func (n *NoteText) GetUserId() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.UserId
}

// NoteTexts is a slice of pointers to NoteText objects.
type NoteTexts []*NoteText

// FindByType locates a note by its type.
func (ns NoteTexts) FindByType(noteType NoteType) *NoteText {
	if len(ns) == 0 || noteType.IsEmpty() {
		return nil
	}
	for _, note := range ns {
		if note.Type == noteType {
			return note
		}
	}
	return nil
}

// HasByType checks if a note of a certain type exists.
func (ns NoteTexts) HasByType(noteType NoteType) bool {
	return ns.FindByType(noteType) != nil
}

// SetByType adds or replaces a note in the collection based on its type.
func (ns *NoteTexts) SetByType(note *NoteText) {
	if note == nil || note.GetType().IsEmpty() {
		return
	}
	for i, n := range *ns {
		if n.GetType() == note.GetType() {
			(*ns)[i] = note
			return
		}
	}
	*ns = append(*ns, note)
}

// RemoveByType deletes a note of a specific type from the collection.
func (ns *NoteTexts) RemoveByType(ntype NoteType) {
	if ntype.IsEmpty() {
		return
	}
	var arr NoteTexts
	for _, n := range *ns {
		if n.GetType() != ntype {
			arr = append(arr, n)
		}
	}
	*ns = arr
}

// Validate checks all notes in the collection for validity.
func (ns NoteTexts) Validate() error {
	return ns.ValidateWithOptions(true, true, true)
}

// ValidateWithOptions performs validation checks on all notes in the collection based on the provided options.
func (ns NoteTexts) ValidateWithOptions(checkType, checkText, checkUserId bool) error {
	if len(ns) == 0 {
		return nil
	}
	for _, n := range ns {
		if err := n.Validate(); err != nil {
			return err
		}
	}
	return nil
}
