package aapp

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Predefined constants for common build types.
const (
	BUILDTYPE_RELEASE BuildType = "release"
	BUILDTYPE_DEBUG   BuildType = "debug"
	BUILDTYPE_PROFILE BuildType = "profile"
)

// knownBuildTypes is a map of recognized build types for efficient lookup.
var knownBuildTypes = map[BuildType]struct{}{
	BUILDTYPE_RELEASE: {},
	BUILDTYPE_DEBUG:   {},
	BUILDTYPE_PROFILE: {},
}

// defaultBuildTypePriorityOrder defines the priority for selecting a default build type.
var defaultBuildTypePriorityOrder = []BuildType{
	BUILDTYPE_RELEASE,
	BUILDTYPE_PROFILE,
	BUILDTYPE_DEBUG,
}

// validBuildTypeRegex ensures BuildType contains only alphanumeric characters, underscores, and ":debug" suffix.
var validBuildTypeRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+(:debug)?$`)

// BuildType represents a type of build as a string.
type BuildType string

// UnmarshalJSON implements json.Unmarshaler to validate the build type during unmarshaling.
// Returns an error if the unmarshaled string is invalid.
func (bt *BuildType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("failed to unmarshal build type: %w", err)
	}
	if !validBuildTypeRegex.MatchString(s) {
		return fmt.Errorf("invalid build type %q: must contain only alphanumeric characters or underscores with optional :debug suffix", s)
	}
	*bt = BuildType(s)
	return nil
}

// IsEmpty checks if the BuildType is empty after trimming whitespace.
// Returns true if the string is empty or contains only whitespace.
func (bt BuildType) IsEmpty() bool {
	return strings.TrimSpace(string(bt)) == ""
}

// String returns the string representation of the BuildType.
func (bt BuildType) String() string {
	return string(bt)
}

// IsDebug checks if the BuildType has a ":debug" suffix.
// Returns true if the suffix is present.
func (bt BuildType) IsDebug() bool {
	return strings.HasSuffix(string(bt), ":debug")
}

// SetDebug adds or removes the ":debug" suffix based on the isDebug flag.
// Returns an error if the BuildType is empty or invalid.
func (bt BuildType) SetDebug(isDebug bool) (BuildType, error) {
	if bt.IsEmpty() {
		return "", fmt.Errorf("build type cannot be empty")
	}
	if !validBuildTypeRegex.MatchString(string(bt)) && !bt.IsDebug() {
		return "", fmt.Errorf("invalid build type %q: must contain only alphanumeric characters or underscores", bt)
	}
	base := strings.TrimSuffix(string(bt), ":debug")
	if isDebug {
		return BuildType(base + ":debug"), nil
	}
	return BuildType(base), nil
}

// IsValid checks if the BuildType is valid (alphanumeric, underscores, optional ":debug" suffix).
// Returns true if valid, false otherwise.
func (bt BuildType) IsValid() bool {
	return validBuildTypeRegex.MatchString(string(bt))
}

// BuildTypes is a slice of BuildType, used to manage multiple build types.
type BuildTypes []BuildType

// Add appends a new BuildType to the slice if it’s valid, non-empty, and not already present.
// Returns an error if the BuildType is invalid.
func (bts *BuildTypes) Add(bt BuildType) error {
	if bt.IsEmpty() {
		return fmt.Errorf("cannot add empty build type")
	}
	if !bt.IsValid() {
		return fmt.Errorf("invalid build type %q: must contain only alphanumeric characters or underscores", bt)
	}
	if !bts.Contains(bt) {
		*bts = append(*bts, bt)
	}
	return nil
}

// Remove deletes a BuildType from the slice if it exists.
func (bts *BuildTypes) Remove(bt BuildType) {
	newSlice := make(BuildTypes, 0, len(*bts))
	for _, b := range *bts {
		if b != bt {
			newSlice = append(newSlice, b)
		}
	}
	*bts = newSlice
}

// Contains checks if the BuildTypes slice contains the specified BuildType.
// Returns true if the BuildType is found.
func (bts BuildTypes) Contains(bt BuildType) bool {
	for _, b := range bts {
		if b == bt {
			return true
		}
	}
	return false
}

// IsKnownType checks if the given BuildType is one of the recognized types.
// Types with a ":debug" suffix are considered known if their base type is recognized.
// Returns true if the type is known.
func (bts BuildTypes) IsKnownType(bt BuildType) bool {
	base := strings.TrimSuffix(string(bt), ":debug")
	_, exists := knownBuildTypes[BuildType(base)]
	return exists
}

// SelectPreferredDefault returns the preferred BuildType based on the priority order:
// release > profile > debug. Returns BUILDTYPE_DEBUG if no known types are found.
func (bts BuildTypes) SelectPreferredDefault() BuildType {
	for _, candidate := range defaultBuildTypePriorityOrder {
		if bts.Contains(candidate) {
			return candidate
		}
	}
	return BUILDTYPE_DEBUG
}
