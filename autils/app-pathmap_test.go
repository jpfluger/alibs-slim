package autils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"path/filepath"
	"testing"
)

// TestAppPathMap_SetPath checks the SetPath method for AppPathMap.
func TestAppPathMap_SetPath(t *testing.T) {
	apm := make(AppPathMap)
	apm.SetPath("DIR_TEST", "/test/dir")
	if apm["DIR_TEST"] != "/test/dir" {
		t.Errorf("SetPath() failed, expected /test/dir, got %s", apm["DIR_TEST"])
	}
}

// TestAppPathMap_GetPath checks the GetPath method for AppPathMap.
func TestAppPathMap_GetPath(t *testing.T) {
	apm := AppPathMap{
		"DIR_TEST": "/test/dir",
	}
	if path := apm.GetPath("DIR_TEST"); path != "/test/dir" {
		t.Errorf("GetPath() failed, expected /test/dir, got %s", path)
	}
	if path := apm.GetPath("DIR_UNKNOWN"); path != "" {
		t.Errorf("GetPath() failed, expected empty string for unknown key, got %s", path)
	}
}

// TestAppPathMap_Validate checks the ValidateExistence method for AppPathMap.
func TestAppPathMap_TestValidate(t *testing.T) {

	dir, err := CreateTempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	if err := os.WriteFile(path.Join(dir, "file"), []byte("test"), 0666); err != nil {
		t.Error(err)
		return
	}

	apm := AppPathMap{
		"DIR_VALID":  dir,
		"FILE_VALID": path.Join(dir, "file"),
	}

	// Test with valid paths
	if err := apm.Validate(); err != nil {
		t.Errorf("ValidateExistence() failed, expected no error, got %v", err)
	}

	// Test with an invalid directory
	apm["DIR_INVALID"] = "/invalid/dir"
	if err := apm.Validate(); err == nil {
		t.Errorf("ValidateExistence() failed, expected error for invalid directory, got nil")
	}

	// Test with an invalid file
	apm["FILE_INVALID"] = "/invalid/file"
	if err := apm.Validate(); err == nil {
		t.Errorf("ValidateExistence() failed, expected error for invalid file, got nil")
	}
}

// TestAppPathMap_TestValidateWithOption checks the ValidateWithOption method for AppPathMap.
func TestAppPathMap_TestValidateWithOption(t *testing.T) {
	// Create a temporary root directory
	dirRoot, err := CreateTempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirRoot) // Cleanup after test

	// Initialize an AppPathMap with relative paths
	apm := AppPathMap{
		DIR_ROOT: dirRoot, // Absolute root path
		DIR_ETC:  "etc",   // Relative path -> should resolve to dirRoot/etc
		DIR_DATA: "data",  // Relative path -> should resolve to dirRoot/data
		DIR_LOGS: "logs",  // Relative path -> should resolve to dirRoot/logs
	}

	// Ensure directories are created before validation
	err = apm.EnsureDirs(dirRoot)
	if err != nil {
		t.Fatalf("EnsureDirs failed: %v", err)
	}

	// Verify directories exist
	for _, key := range []AppPathKey{DIR_ETC, DIR_DATA, DIR_LOGS} {
		if _, err := os.Stat(apm[key]); os.IsNotExist(err) {
			t.Errorf("Expected directory %s to be created, but it was not", apm[key])
		}
	}

	// Run validation after directory creation
	err = apm.ValidateWithOption("")
	if err != nil {
		t.Fatalf("ValidateWithOption failed after EnsureDirs: %v", err)
	}

	// Add a file entry and test its validation
	testFileKey := AppPathKey("FILE_CONFIG")
	apm[testFileKey] = "config.yaml"

	// Create the test file
	testFilePath := filepath.Join(dirRoot, "config.yaml")
	file, err := os.Create(testFilePath)
	if err != nil {
		t.Fatalf("Failed to create test file %s: %v", testFilePath, err)
	}
	file.Close()

	// Run validation again
	err = apm.ValidateWithOption("")
	if err != nil {
		t.Fatalf("ValidateWithOption failed after EnsureDirs: %v", err)
	}

	// Ensure relative paths were updated in apm
	expectedPaths := map[AppPathKey]string{
		DIR_ETC:     filepath.Join(dirRoot, "etc"),
		DIR_DATA:    filepath.Join(dirRoot, "data"),
		DIR_LOGS:    filepath.Join(dirRoot, "logs"),
		testFileKey: filepath.Join(dirRoot, "config.yaml"),
	}

	for key, expected := range expectedPaths {
		actual := apm[key]
		if actual != expected {
			t.Errorf("Path resolution failed for %s: expected %s, got %s", key, expected, actual)
		}
	}

	// Add a non-existent directory to test failure handling
	apm["DIR_MISSING"] = "missing_dir"

	err = apm.ValidateWithOption("")
	if err == nil {
		t.Error("Expected error for missing directory, but got nil")
	}
}

// TestRequireWithOption_Success tests that RequireWithOption succeeds when all required keys exist.
func TestRequireWithOption_Success(t *testing.T) {
	// Create a temporary root directory
	dirRoot, err := CreateTempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirRoot) // Cleanup after test

	// Initialize an AppPathMap with valid paths
	apm := AppPathMap{
		DIR_ROOT: dirRoot,
		DIR_ETC:  "etc",
		DIR_DATA: "data",
		DIR_LOGS: "logs",
	}

	// Ensure directories exist before validation
	err = apm.EnsureDirs(dirRoot)
	if err != nil {
		t.Fatalf("EnsureDirs failed: %v", err)
	}

	// Define required keys
	requiredKeys := AppPathKeys{DIR_ETC, DIR_DATA, DIR_LOGS}

	// Run RequireWithOption (should pass)
	err = apm.RequireWithOption(dirRoot, requiredKeys...)
	if err != nil {
		t.Fatalf("RequireWithOption failed unexpectedly: %v", err)
	}
}

// TestRequireWithOption_MissingKey tests that RequireWithOption fails when a required key is missing.
func TestRequireWithOption_MissingKey(t *testing.T) {
	// Create a temporary root directory
	dirRoot, err := CreateTempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirRoot) // Cleanup after test

	// Initialize an AppPathMap with a missing key (DIR_LOGS is not included)
	apm := AppPathMap{
		DIR_ROOT: dirRoot,
		DIR_ETC:  "etc",
		DIR_DATA: "data",
	}

	// Ensure directories exist before validation
	err = apm.EnsureDirs(dirRoot)
	if err != nil {
		t.Fatalf("EnsureDirs failed: %v", err)
	}

	// Define required keys (including a missing one)
	requiredKeys := AppPathKeys{DIR_ETC, DIR_DATA, DIR_LOGS}

	// Run RequireWithOption (should fail)
	err = apm.RequireWithOption(dirRoot, requiredKeys...)
	if err == nil {
		t.Error("Expected RequireWithOption to fail for missing key DIR_LOGS, but it passed")
	} else if err.Error() != "app path key DIR_LOGS does not exist" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// TestRequireWithOption_EmptyValue tests that RequireWithOption fails when a required key has an empty value.
func TestRequireWithOption_EmptyValue(t *testing.T) {
	// Create a temporary root directory
	dirRoot, err := CreateTempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirRoot) // Cleanup after test

	// Initialize an AppPathMap with an empty value for DIR_DATA
	apm := AppPathMap{
		DIR_ROOT: dirRoot,
		DIR_ETC:  "etc",
		DIR_DATA: "", // Empty value
		DIR_LOGS: "logs",
	}

	// Ensure directories exist before validation
	err = apm.EnsureDirs(dirRoot)
	assert.Error(t, err, "Expected RequireWithOption to fail for an empty key, but it passed")

	// Define required keys
	requiredKeys := AppPathKeys{DIR_ETC, DIR_DATA, DIR_LOGS}

	// Run RequireWithOption (should fail)
	err = apm.RequireWithOption(dirRoot, requiredKeys...)
	if err == nil {
		t.Error("Expected RequireWithOption to fail for empty DIR_DATA, but it passed")
	} else if err.Error() != "app path key DIR_DATA does not exist" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// TestRequireWithOption_ValidatesPaths tests that RequireWithOption validates paths after ensuring they exist.
func TestRequireWithOption_ValidatesPaths(t *testing.T) {
	// Create a temporary root directory
	dirRoot, err := CreateTempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirRoot) // Cleanup after test

	// Initialize an AppPathMap with valid paths
	apm := AppPathMap{
		DIR_ROOT: dirRoot,
		DIR_ETC:  "etc",
		DIR_DATA: "data",
		DIR_LOGS: "logs",
	}

	// Ensure directories exist before validation
	err = apm.EnsureDirs(dirRoot)
	if err != nil {
		t.Fatalf("EnsureDirs failed: %v", err)
	}

	// Define required keys
	requiredKeys := AppPathKeys{DIR_ETC, DIR_DATA, DIR_LOGS}

	// Run RequireWithOption (should pass)
	err = apm.RequireWithOption(dirRoot, requiredKeys...)
	if err != nil {
		t.Fatalf("RequireWithOption failed unexpectedly: %v", err)
	}

	// Verify that paths were resolved correctly
	expectedPaths := map[AppPathKey]string{
		DIR_ETC:  filepath.Join(dirRoot, "etc"),
		DIR_DATA: filepath.Join(dirRoot, "data"),
		DIR_LOGS: filepath.Join(dirRoot, "logs"),
	}

	for key, expected := range expectedPaths {
		actual, exists := apm[key]
		if !exists {
			t.Errorf("Expected key %s to exist in AppPathMap but was not found", key)
		} else if actual != expected {
			t.Errorf("Path resolution failed for %s: expected %s, got %s", key, expected, actual)
		}
	}
}

// TestRequireWithOption_HandlesRelativePaths tests that RequireWithOption correctly handles relative paths.
func TestRequireWithOption_HandlesRelativePaths(t *testing.T) {
	// Create a temporary root directory
	dirRoot, err := CreateTempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirRoot) // Cleanup after test

	// Initialize an AppPathMap with relative paths
	apm := AppPathMap{
		DIR_ROOT: dirRoot,
		DIR_ETC:  "etc",
		DIR_DATA: "data",
		DIR_LOGS: "logs",
	}

	// Ensure directories exist
	err = apm.EnsureDirs(dirRoot)
	if err != nil {
		t.Fatalf("EnsureDirs failed: %v", err)
	}

	// Define required keys
	requiredKeys := AppPathKeys{DIR_ETC, DIR_DATA, DIR_LOGS}

	// Run RequireWithOption (should pass)
	err = apm.RequireWithOption(dirRoot, requiredKeys...)
	if err != nil {
		t.Fatalf("RequireWithOption failed unexpectedly: %v", err)
	}

	// Verify that relative paths were resolved correctly
	for key, path := range apm {
		if !filepath.IsAbs(path) {
			t.Errorf("Path %s is not absolute for key %s", path, key)
		}
	}
}
