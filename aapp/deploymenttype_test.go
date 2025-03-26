package aapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDeploymentType_IsEmpty checks if IsEmpty correctly identifies empty deployment types.
func TestDeploymentType_IsEmpty(t *testing.T) {
	var dt DeploymentType
	assert.True(t, dt.IsEmpty(), "Expected DeploymentType to be empty")

	dt = "dev"
	assert.False(t, dt.IsEmpty(), "Expected DeploymentType to be non-empty")
}

// TestDeploymentType_String checks if String returns the correct string representation of DeploymentType.
func TestDeploymentType_String(t *testing.T) {
	dt := DeploymentType("dev")
	assert.Equal(t, "dev", dt.String(), "Expected DeploymentType string to match")
}

// TestDeploymentType_IsDemo checks if IsDemo correctly identifies demo deployment types.
func TestDeploymentType_IsDemo(t *testing.T) {
	dt := DeploymentType("dev:demo")
	assert.True(t, dt.IsDemo(), "Expected DeploymentType to be a demo")

	dt = "dev"
	assert.False(t, dt.IsDemo(), "Expected DeploymentType to not be a demo")
}

// TestDeploymentType_SetDemo checks if SetDemo correctly adds or removes the ":demo" suffix.
func TestDeploymentType_SetDemo(t *testing.T) {
	dt := DeploymentType("dev")
	dt = dt.SetDemo(true)
	assert.Equal(t, "dev:demo", dt.String(), "Expected DeploymentType to have ':demo' suffix")

	dt = dt.SetDemo(false)
	assert.Equal(t, "dev", dt.String(), "Expected DeploymentType to not have ':demo' suffix")
}

// TestDeploymentTypes_Add checks if Add correctly appends a new DeploymentType to the slice.
func TestDeploymentTypes_Add(t *testing.T) {
	dts := DeploymentTypes{}
	dt := DeploymentType("dev")

	dts.Add(dt)
	assert.Contains(t, dts, dt, "Expected DeploymentTypes to contain the added DeploymentType")

	// Test adding a duplicate DeploymentType
	dts.Add(dt)
	assert.Len(t, dts, 1, "Expected DeploymentTypes to not add duplicate DeploymentType")
}

// TestDeploymentTypes_Remove checks if Remove correctly deletes a DeploymentType from the slice.
func TestDeploymentTypes_Remove(t *testing.T) {
	dt1 := DeploymentType("dev")
	dt2 := DeploymentType("qa")
	dts := DeploymentTypes{dt1, dt2}

	dts.Remove(dt1)
	assert.NotContains(t, dts, dt1, "Expected DeploymentTypes to not contain the removed DeploymentType")
}

// TestDeploymentTypes_Contains checks if Contains correctly identifies if a DeploymentType is in the slice.
func TestDeploymentTypes_Contains(t *testing.T) {
	dt := DeploymentType("dev")
	dts := DeploymentTypes{dt}

	assert.True(t, dts.Contains(dt), "Expected DeploymentTypes to contain the DeploymentType")
}
