package aapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBuildType_IsEmpty checks if IsEmpty correctly identifies empty build types.
func TestBuildType_IsEmpty(t *testing.T) {
	var bt BuildType
	assert.True(t, bt.IsEmpty(), "Expected BuildType to be empty")

	bt = "release"
	assert.False(t, bt.IsEmpty(), "Expected BuildType to be non-empty")
}

// TestBuildType_String checks if String returns the correct string representation of BuildType.
func TestBuildType_String(t *testing.T) {
	bt := BuildType("release")
	assert.Equal(t, "release", bt.String(), "Expected BuildType string to match")
}

// TestBuildTypes_Add checks if Add correctly appends a new BuildType to the slice.
func TestBuildTypes_Add(t *testing.T) {
	bts := BuildTypes{}
	bt := BuildType("release")

	bts.Add(bt)
	assert.Contains(t, bts, bt, "Expected BuildTypes to contain the added BuildType")

	// Test adding a duplicate BuildType
	bts.Add(bt)
	assert.Len(t, bts, 1, "Expected BuildTypes to not add duplicate BuildType")
}

// TestBuildTypes_Remove checks if Remove correctly deletes a BuildType from the slice.
func TestBuildTypes_Remove(t *testing.T) {
	bt1 := BuildType("release")
	bt2 := BuildType("debug")
	bts := BuildTypes{bt1, bt2}

	bts.Remove(bt1)
	assert.NotContains(t, bts, bt1, "Expected BuildTypes to not contain the removed BuildType")
}

// TestBuildTypes_Contains checks if Contains correctly identifies if a BuildType is in the slice.
func TestBuildTypes_Contains(t *testing.T) {
	bt := BuildType("release")
	bts := BuildTypes{bt}

	assert.True(t, bts.Contains(bt), "Expected BuildTypes to contain the BuildType")
}
