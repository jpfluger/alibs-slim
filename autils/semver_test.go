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

func TestAreSemverVersionsValid(t *testing.T) {
	v := func(s string) *semver.Version {
		parsed, err := semver.NewVersion(s)
		if err != nil {
			t.Fatalf("failed to parse version %q: %v", s, err)
		}
		return parsed
	}

	nilVersion := (*semver.Version)(nil)
	zeroVersion := v("0.0.0")
	validVersion := v("1.2.3")

	cases := []struct {
		name     string
		versions []*semver.Version
		wantErr  bool
	}{
		{"all valid", []*semver.Version{validVersion, v("2.0.0")}, false},
		{"one is 0.0.0", []*semver.Version{validVersion, zeroVersion}, true},
		{"nil version", []*semver.Version{validVersion, nilVersion}, true},
		{"single nil", []*semver.Version{nilVersion}, true},
		{"only 0.0.0", []*semver.Version{zeroVersion}, true},
		{"empty input", []*semver.Version{}, false}, // this is up to your intent
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := AreSemverVersionsValid(tc.versions...)
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestHasAnyValidSemverVersion(t *testing.T) {
	v := func(s string) *semver.Version {
		parsed, err := semver.NewVersion(s)
		if err != nil {
			t.Fatalf("failed to parse version %q: %v", s, err)
		}
		return parsed
	}

	cases := []struct {
		name     string
		versions []*semver.Version
		want     bool
	}{
		{"all nil", []*semver.Version{nil, nil}, false},
		{"all empty", []*semver.Version{v("0.0.0"), v("0.0.0")}, false},
		{"mixed nil and 0.0.0", []*semver.Version{nil, v("0.0.0")}, false},
		{"one valid", []*semver.Version{nil, v("1.2.3")}, true},
		{"multiple valid", []*semver.Version{v("1.0.0"), v("2.0.0")}, true},
		{"only one valid among bad", []*semver.Version{v("0.0.0"), v("1.0.0"), nil}, true},
		{"empty input", []*semver.Version{}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := HasAnyValidSemverVersion(tc.versions...)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
