// Package anotes provides a structure and methods for managing a collection of notes.
package anotes

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/areflect"
	"reflect"
	"time"
)

// INote is an interface for note-related actions.
type INote interface {
	GetType() NoteType
	GetDate() time.Time
	GetText() string
	GetUserId() string
}

// INotes is a slice of INote objects.
type INotes []INote

// UnmarshalJSON custom unmarshals INotes from JSON.
func (ns *INotes) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	for _, r := range raw {
		var rawmap map[string]interface{}
		err := json.Unmarshal(r, &rawmap)
		if err != nil {
			return fmt.Errorf("cannot unmarshal json INotes because rawmap could not be unmarshaled; %v", err)
		}

		rawType := ""
		if t, ok := rawmap["type"].(string); !ok {
			rawType = NOTETYPE_TEXT.String()
		} else {
			rawType = NoteType(t).String()
		}

		rtype, err := areflect.TypeManager().FindReflectType(TYPEMANAGER_NOTE, rawType)
		if err != nil {
			return fmt.Errorf("cannot find type struct '%s'; %v", rawType, err)
		}

		obj, ok := reflect.New(rtype).Interface().(INote)
		if !ok {
			return fmt.Errorf("interface is not of type INote")
		}

		if err := json.Unmarshal(r, obj); err != nil {
			return fmt.Errorf(`cannot json unmarshal INote for type '%s'; %v`, rawType, err)
		}

		*ns = append(*ns, obj)
	}
	return nil
}

// FindByType locates a note by its type.
func (ns INotes) FindByType(noteType NoteType) INote {
	if len(ns) == 0 || noteType.IsEmpty() {
		return nil
	}
	for _, note := range ns {
		if note.GetType() == noteType {
			return note
		}
	}
	return nil
}

// HasByType checks if a note of a certain type exists.
func (ns INotes) HasByType(noteType NoteType) bool {
	return ns.FindByType(noteType) != nil
}

// SetByType adds or replaces a note in the collection.
func (ns *INotes) SetByType(note INote) {
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

// Remove deletes a note of a specific type.
func (ns *INotes) Remove(ntype NoteType) {
	if ntype.IsEmpty() {
		return
	}
	var arr INotes
	for _, n := range *ns {
		if n.GetType() != ntype {
			arr = append(arr, n)
		}
	}
	*ns = arr
}
