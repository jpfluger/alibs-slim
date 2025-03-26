package aapp

import (
	"strings"
)

// Predefined constants for different build types.
const (
	BUILDTYPE_RELEASE BuildType = "release"
	BUILDTYPE_DEBUG   BuildType = "debug"
	BUILDTYPE_PROFILE BuildType = "profile"
)

// BuildType represents a type of build as a string.
type BuildType string

// IsEmpty checks if the BuildType is empty after trimming whitespace.
func (bt BuildType) IsEmpty() bool {
	return strings.TrimSpace(string(bt)) == ""
}

// String returns the string representation of the BuildType.
func (bt BuildType) String() string {
	return string(bt)
}

// BuildTypes is a slice of BuildType, used to handle multiple build types.
type BuildTypes []BuildType

// Add appends a new BuildType to the slice if it's not empty and not already present.
func (bts *BuildTypes) Add(bt BuildType) {
	if !bt.IsEmpty() && !bts.Contains(bt) {
		*bts = append(*bts, bt)
	}
}

// Remove deletes a BuildType from the slice if it exists.
func (bts *BuildTypes) Remove(bt BuildType) {
	for i, b := range *bts {
		if b == bt {
			*bts = append((*bts)[:i], (*bts)[i+1:]...)
			break
		}
	}
}

// Contains checks if the BuildTypes slice contains the specified BuildType.
func (bts BuildTypes) Contains(bt BuildType) bool {
	for _, b := range bts {
		if b == bt {
			return true
		}
	}
	return false
}
