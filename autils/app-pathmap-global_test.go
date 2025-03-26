package autils

import (
	"os"
	"path/filepath"
	"testing"
)

// TestAppPathMapGlobal_SetAndGetAppPathMap tests setting and retrieving the global AppPathMap.
func TestAppPathMapGlobal_SetAndGetAppPathMap(t *testing.T) {
	// Create a temporary root directory
	dirRoot, err := CreateTempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirRoot)

	// Initialize an AppPathMap with relative paths
	apm := AppPathMap{
		DIR_ROOT: dirRoot, // Absolute root path
		DIR_ETC:  "etc",   // Relative path
		DIR_DATA: "data",  // Relative path
		DIR_LOGS: "logs",  // Relative path
	}

	// Set the global map
	err = SetAppPathMap(apm)
	if err != nil {
		t.Fatalf("SetAppPathMap failed: %v", err)
	}

	// Retrieve the stored map
	storedMap := GetAppPathMap()

	// Check if all paths are stored and resolved properly
	expectedPaths := map[AppPathKey]string{
		DIR_ETC:  filepath.Join(dirRoot, "etc"),
		DIR_DATA: filepath.Join(dirRoot, "data"),
		DIR_LOGS: filepath.Join(dirRoot, "logs"),
	}

	for key, expected := range expectedPaths {
		actual, exists := storedMap[key]
		if !exists {
			t.Errorf("Expected key %s not found in stored map", key)
		} else if actual != expected {
			t.Errorf("Path mismatch for key %s: expected %s, got %s", key, expected, actual)
		}
	}
}

// TestAppPathMapGlobal_GetAppPath tests retrieving a single path from the global AppPathMap.
func TestAppPathMapGlobal_GetAppPath(t *testing.T) {
	// Create a temporary root directory
	dirRoot, err := CreateTempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirRoot)

	// Set up paths
	apm := AppPathMap{
		DIR_ROOT: dirRoot,
		DIR_ETC:  "etc",
	}
	err = SetAppPathMap(apm)
	if err != nil {
		t.Fatalf("SetAppPathMap failed: %v", err)
	}

	// Test retrieving an existing path
	expectedEtcPath := filepath.Join(dirRoot, "etc")
	actualEtcPath, exists := GetAppPath(DIR_ETC)
	if !exists {
		t.Errorf("Expected DIR_ETC to exist in AppPathMap but was not found")
	}
	if actualEtcPath != expectedEtcPath {
		t.Errorf("Mismatch for DIR_ETC: expected %s, got %s", expectedEtcPath, actualEtcPath)
	}

	// Test retrieving a non-existent key
	_, exists = GetAppPath(AppPathKey("NON_EXISTENT_KEY"))
	if exists {
		t.Errorf("Expected NON_EXISTENT_KEY to not exist, but it was found")
	}
}

// TestAppPathMapGlobal_SetAppPath tests updating a single path in the global AppPathMap.
func TestAppPathMapGlobal_SetAppPath(t *testing.T) {
	// Create a temporary root directory
	dirRoot, err := CreateTempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirRoot)

	// Set up paths
	apm := AppPathMap{
		DIR_ROOT: dirRoot,
		DIR_ETC:  "etc",
	}
	err = SetAppPathMap(apm)
	if err != nil {
		t.Fatalf("SetAppPathMap failed: %v", err)
	}

	// Update the DIR_ETC path
	newEtcPath := filepath.Join(dirRoot, "config")
	SetAppPath(DIR_ETC, newEtcPath)

	// Verify that the new path is stored
	actualEtcPath, exists := GetAppPath(DIR_ETC)
	if !exists {
		t.Errorf("DIR_ETC key not found after update")
	}
	if actualEtcPath != newEtcPath {
		t.Errorf("DIR_ETC path mismatch after update: expected %s, got %s", newEtcPath, actualEtcPath)
	}
}
