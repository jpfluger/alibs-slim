// Package anotes provides utilities for note handling.
package anotes

import (
	// Importing necessary packages.
	"github.com/jpfluger/alibs-slim/autils" // Custom utility package.
	"strings"                                              // Standard library package for string manipulation.
)

// NoteType defines a custom type for notes.
type NoteType string

// IsEmpty checks if the NoteType is empty after trimming spaces.
func (nt NoteType) IsEmpty() bool {
	rt := strings.TrimSpace(string(nt)) // Trim spaces from the NoteType.
	return rt == ""                     // Return true if the trimmed NoteType is empty.
}

// TrimSpace trims spaces from the NoteType and returns a new NoteType.
func (nt NoteType) TrimSpace() NoteType {
	rt := strings.TrimSpace(string(nt)) // Trim spaces from the NoteType.
	return NoteType(rt)                 // Return the trimmed NoteType.
}

// String converts NoteType to a string.
func (nt NoteType) String() string {
	return string(nt) // Return the NoteType as a string.
}

// ToStringTrimLower converts NoteType to a string, trims spaces, and makes it lowercase.
func (nt NoteType) ToStringTrimLower() string {
	// Use the custom utility function to convert, trim, and lowercase the NoteType.
	return autils.ToStringTrimLower(nt.String())
}

// NoteTypes defines a slice of NoteType.
type NoteTypes []NoteType
