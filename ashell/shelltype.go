package ashell

import (
	"strings" // Importing the strings package for string manipulation.
)

const (
	SHELLTYPE_BASH     ShellType = "bash"
	SHELLTYPE_SH       ShellType = "sh"
	SHELLTYPE_WIN      ShellType = "win"
	SHELLTYPE_PWS      ShellType = "pws"
	SHELLTYPE_OSX      ShellType = "osx"
	SHELLTYPE_PYTHON   ShellType = "python"
	SHELLTYPE_PYTHON3  ShellType = "python3"
	SHELLTYPE_DOCKER   ShellType = "docker"
	SHELLTYPE_DCOMPOSE ShellType = "docker-compose"
)

// ShellType defines a custom type for shell types.
type ShellType string

// IsEmpty checks if the ShellType is empty after trimming spaces.
func (st ShellType) IsEmpty() bool {
	// Trim spaces from the ShellType and check if the result is an empty string.
	return strings.TrimSpace(string(st)) == ""
}

// TrimSpace trims spaces from the ShellType and returns a new ShellType.
func (st ShellType) TrimSpace() ShellType {
	// Trim spaces from the ShellType and return the result as a new ShellType.
	return ShellType(strings.TrimSpace(string(st)))
}

// String converts ShellType to a string.
func (st ShellType) String() string {
	// Convert the ShellType to a string and return it.
	return string(st)
}

// ToStringTrimLower converts ShellType to a string, trims spaces, and makes it lowercase.
func (st ShellType) ToStringTrimLower() string {
	// Convert the ShellType to a string, trim spaces, convert to lowercase, and return the result.
	return strings.ToLower(st.TrimSpace().String())
}

// ShellTypes defines a slice of ShellType.
type ShellTypes []ShellType

// Contains checks if the ShellTypes slice contains a specific ShellType.
func (sts ShellTypes) Contains(st ShellType) bool {
	// Iterate over the ShellTypes slice.
	for _, t := range sts {
		// Check if the current ShellType matches the specified ShellType.
		if t == st {
			return true // Return true if a match is found.
		}
	}
	return false // Return false if no match is found.
}

// Add appends a new ShellType to the ShellTypes slice if it's not already present.
func (sts *ShellTypes) Add(st ShellType) {
	// Check if the ShellType is not already in the slice.
	if !sts.Contains(st) {
		*sts = append(*sts, st) // Append the new ShellType to the slice.
	}
}

// Remove deletes a ShellType from the ShellTypes slice.
func (sts *ShellTypes) Remove(st ShellType) {
	// Create a new slice to store the result.
	var result ShellTypes
	// Iterate over the ShellTypes slice.
	for _, t := range *sts {
		// If the current ShellType does not match the specified ShellType, add it to the result slice.
		if t != st {
			result = append(result, t)
		}
	}
	*sts = result // Set the ShellTypes slice to the result.
}
