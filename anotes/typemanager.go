// Package anotes provides utilities for note management with reflection.
package anotes

import (
	// Importing necessary packages for reflection.
	"github.com/jpfluger/alibs-slim/areflect" // Custom reflection utility package.
	"reflect"                                                // Standard library package for reflection.
)

// Constant for the note type manager.
const TYPEMANAGER_NOTE = "note"

// init registers the note types with the type manager upon package initialization.
func init() {
	// Ignoring the error on purpose here. In production code, handle the error appropriately.
	_ = areflect.TypeManager().Register(TYPEMANAGER_NOTE, "anotes", returnTypeManagerNotes)
}

// returnTypeManagerNotes returns the reflect.Type corresponding to the provided typeName.
func returnTypeManagerNotes(typeName string) (reflect.Type, error) {
	var rtype reflect.Type // nil is the zero value for pointers, maps, slices, channels, and function types, interfaces, and other compound types.
	switch NoteType(typeName) {
	case NOTETYPE_FLAG:
		// Return the type of NoteFlag if typeName is "flag".
		rtype = reflect.TypeOf(NoteFlag{})
	case NOTETYPE_IMAGE:
		// Return the type of NoteFlag if typeName is "flag".
		rtype = reflect.TypeOf(NoteImage{})
	case NOTETYPE_IMAGE_FLAG:
		// Return the type of NoteFlag if typeName is "flag".
		rtype = reflect.TypeOf(NoteImageFlag{})
	case NOTETYPE_TEXT:
		// Return the type of Note if typeName is "text" or empty.
		rtype = reflect.TypeOf(NoteText{})
	}
	// Return the determined reflect.Type and no error.
	return rtype, nil
}
