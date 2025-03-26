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

// NewVersionFromString parses a semantic version from a string.
// Returns nil if parsing fails and doReturnNilElseUninitialized is true, otherwise returns an uninitialized version.
func NewVersionFromString(version string, doReturnNilElseUninitialized bool) *semver.Version {
	v, err := semver.NewVersion(version)
	if err != nil {
		if doReturnNilElseUninitialized {
			return nil // Return nil if the flag is set.
		} else {
			return &semver.Version{} // Return an uninitialized version otherwise.
		}
	}
	return v
}

// IsSemverVersionInvalid checks if the provided semantic version is invalid.
func IsSemverVersionInvalid(version *semver.Version) error {
	if version == nil {
		return fmt.Errorf("version is nil")
	}
	if !version.GreaterThan(&semver.Version{}) {
		return fmt.Errorf("version is empty")
	}
	if version.String() == "0.0.0" {
		return fmt.Errorf("version cannot be 0.0.0")
	}
	return nil
}
