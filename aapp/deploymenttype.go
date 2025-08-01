package aapp

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Predefined constants for common deployment types.
const (
	DEPLOYMENTTYPE_LOCAL DeploymentType = "local"
	DEPLOYMENTTYPE_DEV   DeploymentType = "dev"
	DEPLOYMENTTYPE_QA    DeploymentType = "qa"
	DEPLOYMENTTYPE_PROD  DeploymentType = "prod"
)

// knownDeploymentTypes is a map of recognized deployment types for efficient lookup.
var knownDeploymentTypes = map[DeploymentType]struct{}{
	DEPLOYMENTTYPE_LOCAL: {},
	DEPLOYMENTTYPE_DEV:   {},
	DEPLOYMENTTYPE_QA:    {},
	DEPLOYMENTTYPE_PROD:  {},
}

// defaultDeploymentTypePriorityOrder defines the priority for selecting a default deployment type.
var defaultDeploymentTypePriorityOrder = []DeploymentType{
	DEPLOYMENTTYPE_PROD,
	DEPLOYMENTTYPE_QA,
	DEPLOYMENTTYPE_DEV,
	DEPLOYMENTTYPE_LOCAL,
}

// validDeploymentTypeRegex ensures DeploymentType contains only alphanumeric characters, underscores, and ":demo" suffix.
var validDeploymentTypeRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+(:demo)?$`)

// DeploymentType represents a type of deployment as a string.
type DeploymentType string

// UnmarshalJSON implements json.Unmarshaler to validate the deployment type during unmarshaling.
// Returns an error if the unmarshaled string is invalid.
func (dt *DeploymentType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("failed to unmarshal deployment type: %w", err)
	}
	if !validDeploymentTypeRegex.MatchString(s) {
		return fmt.Errorf("invalid deployment type %q: must contain only alphanumeric characters or underscores with optional :demo suffix", s)
	}
	*dt = DeploymentType(s)
	return nil
}

// IsEmpty checks if the DeploymentType is empty after trimming whitespace.
// Returns true if the string is empty or contains only whitespace.
func (dt DeploymentType) IsEmpty() bool {
	return strings.TrimSpace(string(dt)) == ""
}

// String returns the string representation of the DeploymentType.
func (dt DeploymentType) String() string {
	return string(dt)
}

// IsDemo checks if the DeploymentType has a ":demo" suffix.
// Returns true if the suffix is present.
func (dt DeploymentType) IsDemo() bool {
	return strings.HasSuffix(string(dt), ":demo")
}

// SetDemo adds or removes the ":demo" suffix based on the isDemo flag.
// Returns an error if the DeploymentType is empty or invalid.
func (dt DeploymentType) SetDemo(isDemo bool) (DeploymentType, error) {
	if dt.IsEmpty() {
		return "", fmt.Errorf("deployment type cannot be empty")
	}
	if !validDeploymentTypeRegex.MatchString(string(dt)) && !dt.IsDemo() {
		return "", fmt.Errorf("invalid deployment type %q: must contain only alphanumeric characters or underscores", dt)
	}
	base := strings.TrimSuffix(string(dt), ":demo")
	if isDemo {
		return DeploymentType(base + ":demo"), nil
	}
	return DeploymentType(base), nil
}

// IsValid checks if the DeploymentType is valid (alphanumeric, underscores, optional ":demo" suffix).
// Returns true if valid, false otherwise.
func (dt DeploymentType) IsValid() bool {
	return validDeploymentTypeRegex.MatchString(string(dt))
}

// DeploymentTypes is a slice of DeploymentType, used to manage multiple deployment types.
type DeploymentTypes []DeploymentType

// Add appends a new DeploymentType to the slice if it’s valid, non-empty, and not already present.
// Returns an error if the DeploymentType is invalid.
func (dts *DeploymentTypes) Add(dt DeploymentType) error {
	if dt.IsEmpty() {
		return fmt.Errorf("cannot add empty deployment type")
	}
	if !dt.IsValid() {
		return fmt.Errorf("invalid deployment type %q: must contain only alphanumeric characters or underscores", dt)
	}
	if !dts.Contains(dt) {
		*dts = append(*dts, dt)
	}
	return nil
}

// Remove deletes a DeploymentType from the slice if it exists.
func (dts *DeploymentTypes) Remove(dt DeploymentType) {
	newSlice := make(DeploymentTypes, 0, len(*dts))
	for _, d := range *dts {
		if d != dt {
			newSlice = append(newSlice, d)
		}
	}
	*dts = newSlice
}

// Contains checks if the DeploymentTypes slice contains the specified DeploymentType.
// Returns true if the DeploymentType is found.
func (dts DeploymentTypes) Contains(dt DeploymentType) bool {
	for _, d := range dts {
		if d == dt {
			return true
		}
	}
	return false
}

// IsKnownType checks if the given DeploymentType is one of the recognized types.
// Types with a ":demo" suffix are considered known if their base type is recognized.
// Returns true if the type is known.
func (dts DeploymentTypes) IsKnownType(dt DeploymentType) bool {
	base := strings.TrimSuffix(string(dt), ":demo")
	_, exists := knownDeploymentTypes[DeploymentType(base)]
	return exists
}

// SelectPreferredDefault returns the preferred DeploymentType based on the priority order:
// prod > qa > dev > local. Returns DEPLOYMENTTYPE_LOCAL if no known types are found.
func (dts DeploymentTypes) SelectPreferredDefault() DeploymentType {
	for _, candidate := range defaultDeploymentTypePriorityOrder {
		if dts.Contains(candidate) {
			return candidate
		}
	}
	return DEPLOYMENTTYPE_LOCAL
}
