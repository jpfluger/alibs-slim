package autils

import (
	"github.com/Masterminds/semver/v3"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// TestSemverEmpty checks if a new semver.Version instance is initialized to "0.0.0" and that an empty string is not a valid version.
func TestSemverEmpty(t *testing.T) {
	v := semver.Version{}
	if v.String() != "0.0.0" {
		t.Fatalf("semver should be version 0.0.0 string but have %s", v.String())
	}

	_, err := semver.NewVersion("")
	if err == nil {
		t.Fatal("expected error when creating new version with empty string")
	}
}

// TestSemverSort checks if a slice of semver.Version instances is sorted correctly in both ascending and descending order.
func TestSemverSort(t *testing.T) {
	ascOrder := []string{"0.4.2", "1.0.0", "1.2.3", "1.3.0", "2.0.0", "3.2.0+incompatible", "10.0.0-alpha", "10.0.0-beta", "10.0.0"}
	descOrder := []string{"10.0.0", "10.0.0-beta", "10.0.0-alpha", "3.2.0+incompatible", "2.0.0", "1.3.0", "1.2.3", "1.0.0", "0.4.2"}

	raw := []string{"1.2.3", "1.0", "10.0", "1.3", "2", "10.0.0-beta", "0.4.2", "10.0.0-alpha", "3.2.0+incompatible"}
	vs := make([]*semver.Version, len(raw))
	for i, r := range raw {
		v, err := semver.NewVersion(r)
		if err != nil {
			t.Fatalf("Error parsing version: %s", err)
		}
		vs[i] = v
	}

	// Sort in ascending order and check.
	sort.Sort(semver.Collection(vs))
	for i, v := range vs {
		if v.String() != ascOrder[i] {
			t.Fatalf("ascending order mismatch: expected %s but have %s", ascOrder[i], v.String())
		}
	}

	// Sort in descending order and check.
	sort.Slice(vs, func(i, j int) bool {
		return vs[i].GreaterThan(vs[j])
	})
	for i, v := range vs {
		if v.String() != descOrder[i] {
			t.Fatalf("descending order mismatch: expected %s but have %s", descOrder[i], v.String())
		}
	}
}

// TestSemverRange checks if semver constraints correctly validate versions within a specified range.
func TestSemverRange(t *testing.T) {
	c, err := semver.NewConstraint(">= 1.2.3-0, < 1.4.0-0")
	if err != nil {
		t.Fatal(err)
	}

	// Version below range should fail.
	v, err := semver.NewVersion("1.1")
	if err != nil {
		t.Fatal(err)
	}
	if _, errs := c.Validate(v); len(errs) == 0 {
		t.Fatal("expected error for version outside of range")
	}

	// Version within range should pass.
	v, err = semver.NewVersion("1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	if _, errs := c.Validate(v); len(errs) > 0 {
		t.Fatal(errs)
	}

	// Pre-release within range should pass.
	v, err = semver.NewVersion("1.3.0-beta")
	if err != nil {
		t.Fatal(err)
	}
	if _, errs := c.Validate(v); len(errs) > 0 {
		t.Fatal(errs)
	}

	// Version above range should fail.
	v, err = semver.NewVersion("1.4.0-alpha")
	if err != nil {
		t.Fatal(err)
	}
	if _, errs := c.Validate(v); len(errs) == 0 {
		t.Fatal("expected error for version outside of range")
	}
}

// TestSemvers checks if semver instances are created correctly from strings and validates them.
func TestSemvers(t *testing.T) {
	goodVersions := []string{
		"0.0.0",  // Edge case: technically valid but often used as a placeholder.
		"v0.0.0", // Edge case with 'v' prefix.
		"1.2.3",
		"v1.2.3",
		"1.2.3-beta.1+build345",
		"v1.2.3-beta.1+build345",
		"1.2.3-nonce",
		"v1.2.3-nonce",
	}

	for _, item := range goodVersions {
		ver, err := semver.NewVersion(item)
		if err != nil {
			t.Errorf("Error parsing version %s: %v", item, err)
		}

		if item == "0.0.0" || item == "v0.0.0" {
			if err := IsSemverVersionInvalid(ver); err == nil {
				t.Errorf("Expected error for version %s but got nil", item)
			}
		} else {
			if err := IsSemverVersionInvalid(ver); err != nil {
				t.Errorf("Expected no error for version %s but got %v", item, err)
			}
		}
	}

	badVersions := []string{
		"0.0.0",  // Edge case: technically valid but often used as a placeholder.
		"v0.0.0", // Edge case with 'v' prefix.
		"0.0.0-nonce",
		"v0.0.0-nonce",
	}

	for _, item := range badVersions {
		ver, err := semver.NewVersion(item)
		if err != nil {
			t.Errorf("Error parsing version %s: %v", item, err)
		}
		if err := IsSemverVersionInvalid(ver); err == nil {
			t.Errorf("Expected error for version %s but got nil", item)
		}
	}
}

func TestNewSemverListFromStrings(t *testing.T) {
	// Test with valid version strings.
	versions := []string{"1.0.0", "2.0.0"}
	semvers, err := NewSemverListFromStrings(versions)
	if err != nil {
		t.Fatalf("NewSemverListFromStrings() returned an error: %v", err)
	}
	if len(semvers) != 2 {
		t.Errorf("NewSemverListFromStrings() returned %d semvers; want 2", len(semvers))
	}
}

func TestNewSemverListFromDirectory(t *testing.T) {
	// Create a temporary directory to simulate version directories.
	tempDir, err := os.MkdirTemp("", "semverdir-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create subdirectories to represent versions.
	os.Mkdir(filepath.Join(tempDir, "1.0.0"), 0755)
	os.Mkdir(filepath.Join(tempDir, "2.0.0"), 0755)

	// Test creating a Semvers slice from the directory.
	semvers, err := NewSemverListFromDirectory(tempDir)
	if err != nil {
		t.Fatalf("NewSemverListFromDirectory() returned an error: %v", err)
	}
	if len(semvers) != 2 {
		t.Errorf("NewSemverListFromDirectory() returned %d semvers; want 2", len(semvers))
	}
}

func TestGetHighestSemverDirectory(t *testing.T) {
	// Create a temporary directory to simulate version directories.
	tempDir, err := os.MkdirTemp("", "semverdir-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create subdirectories to represent versions.
	os.Mkdir(filepath.Join(tempDir, "1.0.0"), 0755)
	highestVersionDir := filepath.Join(tempDir, "2.0.0")
	os.Mkdir(highestVersionDir, 0755)

	// Test getting the highest version directory.
	highestDir, err := GetHighestSemverDirectory(tempDir)
	if err != nil {
		t.Fatalf("GetHighestSemverDirectory() returned an error: %v", err)
	}
	if highestDir != highestVersionDir {
		t.Errorf("GetHighestSemverDirectory() = %v, want %v", highestDir, highestVersionDir)
	}
}
