package aapp

import (
	"strings"
)

// BuildName represents the name given to a binary when built.
type BuildName string

// IsEmpty checks if the BuildName is empty after trimming whitespace.
// Returns true if the string is empty or contains only whitespace.
func (bn BuildName) IsEmpty() bool {
	return strings.TrimSpace(string(bn)) == ""
}

// String returns the string representation of the BuildName.
func (bn BuildName) String() string {
	return string(bn)
}
