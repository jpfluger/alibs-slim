package autils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Masterminds/semver/v3"
)

// Semvers is a slice of pointers to Semver, representing a collection of semantic versions.
type Semvers []*Semver

// NewSemverListFromStrings creates a Semvers slice from a list of version strings.
func NewSemverListFromStrings(versions []string) (Semvers, error) {
	semvers := Semvers{}
	if len(versions) == 0 {
		return semvers, nil // Return an empty Semvers slice if the input is empty.
	}
	for _, version := range versions {
		v, err := semver.NewVersion(version)
		if err != nil {
			return nil, fmt.Errorf("invalid semantic version: %v", err) // Return an error if any version is invalid.
		}
		semvers = append(semvers, &Semver{Version: *v}) // Append the parsed version to the Semvers slice.
	}

	// Sort the Semvers slice in descending order by version.
	sort.Slice(semvers, func(i, j int) bool {
		return semvers[i].Version.GreaterThan(&semvers[j].Version)
	})

	return semvers, nil
}

// NewSemverListFromDirectory creates a Semvers slice from the directory names within the given directory.
func NewSemverListFromDirectory(dir string) (Semvers, error) {
	dirs, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err) // Return an error if reading the directory fails.
	}

	semvers := Semvers{}

	for _, d := range dirs {
		if d.IsDir() {
			v, err := semver.NewVersion(d.Name())
			if err != nil {
				return nil, fmt.Errorf("invalid directory name for semver: %v", err) // Return an error if the directory name is not a valid semver.
			}
			semv := &Semver{
				Version:   *v,
				Directory: filepath.Join(dir, d.Name()), // Set the directory path for the version.
			}
			semvers = append(semvers, semv)
		}
	}

	// Sort the Semvers slice in descending order by version.
	sort.Slice(semvers, func(i, j int) bool {
		return semvers[i].Version.GreaterThan(&semvers[j].Version)
	})

	return semvers, nil
}

// GetHighestSemverDirectory returns the directory path with the highest semantic version within the given directory.
func GetHighestSemverDirectory(dir string) (string, error) {
	semvers, err := NewSemverListFromDirectory(dir)
	if err != nil {
		return "", err // Return an error if creating the Semvers slice fails.
	}
	if len(semvers) == 0 {
		return "", nil // Return an empty string if there are no semantic versions.
	}

	return semvers[0].Directory, nil // Return the directory of the highest version.
}
