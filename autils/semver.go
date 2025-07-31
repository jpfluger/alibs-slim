package autils

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
)

// Semver represents a semantic versioning information with optional directory and tag.
type Semver struct {
	Version   semver.Version `json:"version"`             // Version follows semantic versioning.
	Directory string         `json:"directory,omitempty"` // Directory where the version is stored (if applicable).
	Tag       string         `json:"tag,omitempty"`       // Tag associated with the version (alternative to Directory).
}

// NewSemver creates a new Semver instance with the provided version, directory, and tag.
func NewSemver(version *semver.Version, directory string, tag string) *Semver {
	smvr := &Semver{
		Version:   semver.Version{}, // Initialize with an empty version.
		Directory: directory,
		Tag:       tag,
	}
	if version != nil {
		smvr.Version = *version // Set the version if provided.
	}
	return smvr
}

// MustNewVersion creates a new semantic version from a string, panicking if the string is not a valid semver.
func MustNewVersion(version string) semver.Version {
	v, err := semver.NewVersion(version)
	if err != nil {
		panic(fmt.Sprintf("invalid version: %s", version)) // Panic instead of returning an empty version.
	}
	return *v
}

// MustNewVersionPtr creates a new semantic version from a string, panicking if the string is not a valid semver.
func MustNewVersionPtr(version string) *semver.Version {
	v, err := semver.NewVersion(version)
	if err != nil {
		panic(fmt.Sprintf("invalid version: %s", version)) // Panic instead of returning an empty version.
	}
	return v
}

// NewVersionFromString parses a semantic version from a string.
// Returns nil if parsing fails and doReturnNilElseUninitialized is true,
// otherwise returns an uninitialized version (0.0.0).
func NewVersionFromString(version string, doReturnNilElseUninitialized bool) *semver.Version {
	v, err := semver.NewVersion(version)
	if err != nil {
		if doReturnNilElseUninitialized {
			return nil
		}
		return &semver.Version{}
	}
	return v
}

// IsSemverVersionInvalid checks if the provided semantic version is invalid.
// Returns an error if version is nil, 0.0.0, or not greater than 0.0.0.
func IsSemverVersionInvalid(version *semver.Version) error {
	if version == nil {
		return fmt.Errorf("version is nil")
	}
	if version.String() == "0.0.0" {
		return fmt.Errorf("version cannot be 0.0.0")
	}
	if !version.GreaterThan(&semver.Version{}) {
		return fmt.Errorf("version must be greater than 0.0.0")
	}
	return nil
}

// IsSemverValid returns true if all provided semantic versions are valid.
func IsSemverValid(version ...*semver.Version) bool {
	return AreSemverVersionsValid(version...) == nil
}

// AreSemverVersionsValid checks one or more semver.Version pointers for validity.
// Returns the first error found, or nil if all are valid.
func AreSemverVersionsValid(versions ...*semver.Version) error {
	for i, version := range versions {
		if err := IsSemverVersionInvalid(version); err != nil {
			return fmt.Errorf("version at index %d (%v) is invalid: %w", i, version, err)
		}
	}
	return nil
}

// HasAnyValidSemverVersion returns true if at least one provided version is valid.
func HasAnyValidSemverVersion(versions ...*semver.Version) bool {
	for _, version := range versions {
		if IsSemverVersionInvalid(version) == nil {
			return true
		}
	}
	return false
}
