package autils

import (
	"github.com/Masterminds/semver/v3"
	"testing"
)

// TestSemver_NewSemver tests the NewSemver function for creating a new Semver instance.
func TestSemver_NewSemver(t *testing.T) {
	version, err := semver.NewVersion("1.0.0")
	if err != nil {
		t.Fatalf("Failed to create semver version: %v", err)
	}
	s := NewSemver(version, "dir", "v1.0.0")
	if s.Version.String() != "1.0.0" {
		t.Errorf("NewSemver() Version = %v, want %v", s.Version.String(), "1.0.0")
	}
	if s.Directory != "dir" {
		t.Errorf("NewSemver() Directory = %v, want %v", s.Directory, "dir")
	}
	if s.Tag != "v1.0.0" {
		t.Errorf("NewSemver() Tag = %v, want %v", s.Tag, "v1.0.0")
	}
}

// TestSemver_MustNewVersion tests the MustNewVersion function for creating a new semver.Version.
func TestSemver_MustNewVersion(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustNewVersion() did not panic on invalid version")
		}
	}()
	_ = MustNewVersion("invalid")
}

// TestSemver_NewVersionFromString tests the NewVersionFromString function for parsing a semver from string.
func TestSemver_NewVersionFromString(t *testing.T) {
	v := NewVersionFromString("1.0.0", false)
	if v == nil {
		t.Errorf("NewVersionFromString() = nil, want non-nil")
	}
	if v.String() != "1.0.0" {
		t.Errorf("NewVersionFromString() = %v, want %v", v.String(), "1.0.0")
	}

	v = NewVersionFromString("invalid", true)
	if v != nil {
		t.Errorf("NewVersionFromString() = %v, want nil for invalid version", v)
	}
}

// TestSemver_IsSemverVersionInvalid tests the IsSemverVersionInvalid function for checking version validity.
func TestSemver_IsSemverVersionInvalid(t *testing.T) {
	v := MustNewVersion("1.0.0")
	if err := IsSemverVersionInvalid(&v); err != nil {
		t.Errorf("IsSemverVersionInvalid() = %v, want nil", err)
	}

	v = semver.Version{}
	if err := IsSemverVersionInvalid(&v); err == nil {
		t.Errorf("IsSemverVersionInvalid() = nil, want error for empty version")
	}
}
