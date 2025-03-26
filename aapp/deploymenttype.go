package aapp

import (
	"fmt"
	"strings"
)

// Predefined constants for different deployment types.
const (
	DEPLOYMENTTYPE_DEV  DeploymentType = "dev"
	DEPLOYMENTTYPE_QA   DeploymentType = "qa"
	DEPLOYMENTTYPE_PROD DeploymentType = "prod"
)

// DeploymentType represents a type of deployment as a string.
type DeploymentType string

// IsEmpty checks if the DeploymentType is empty after trimming whitespace.
func (dt DeploymentType) IsEmpty() bool {
	return strings.TrimSpace(string(dt)) == ""
}

// String returns the string representation of the DeploymentType.
func (dt DeploymentType) String() string {
	return string(dt)
}

// IsDemo checks if the DeploymentType has a ":demo" suffix.
func (dt DeploymentType) IsDemo() bool {
	return strings.HasSuffix(string(dt), ":demo")
}

// SetDemo adds or removes the ":demo" suffix based on the isDemo flag.
func (dt DeploymentType) SetDemo(isDemo bool) DeploymentType {
	if dt.IsDemo() {
		if !isDemo {
			return DeploymentType(strings.TrimSuffix(string(dt), ":demo"))
		}
		return dt
	}
	if !isDemo {
		return dt
	}
	return DeploymentType(fmt.Sprintf("%s:demo", strings.TrimSpace(string(dt))))
}

// DeploymentTypes is a slice of DeploymentType, used to handle multiple deployment types.
type DeploymentTypes []DeploymentType

// Add appends a new DeploymentType to the slice if it's not empty and not already present.
func (dts *DeploymentTypes) Add(dt DeploymentType) {
	if !dt.IsEmpty() && !dts.Contains(dt) {
		*dts = append(*dts, dt)
	}
}

// Remove deletes a DeploymentType from the slice if it exists.
func (dts *DeploymentTypes) Remove(dt DeploymentType) {
	for i, d := range *dts {
		if d == dt {
			*dts = append((*dts)[:i], (*dts)[i+1:]...)
			break
		}
	}
}

// Contains checks if the DeploymentTypes slice contains the specified DeploymentType.
func (dts DeploymentTypes) Contains(dt DeploymentType) bool {
	for _, d := range dts {
		if d == dt {
			return true
		}
	}
	return false
}
