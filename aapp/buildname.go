package aapp

import (
	"strings"
)

// BuildName is the name given to the binary when built.
type BuildName string

// IsEmpty checks if the BuildName is empty after trimming whitespace.
func (bn BuildName) IsEmpty() bool {
	return strings.TrimSpace(string(bn)) == ""
}

// String returns the string representation of the BuildName.
func (bn BuildName) String() string {
	return string(bn)
}
