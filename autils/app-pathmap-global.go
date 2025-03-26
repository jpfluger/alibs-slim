package autils

import (
	"fmt"
	"os"
	"sync"
)

// Global AppPathMap stored in sync.Map
var globalsAppPathMap sync.Map

// SetAppPathMap initializes or updates the global AppPathMap.
func SetAppPathMap(apm AppPathMap) error {
	// Ensure DIR_ROOT is set
	dirRoot := apm.GetPath(DIR_ROOT)
	if dirRoot == "" {
		var err error
		dirRoot, err = os.Getwd() // Default to current working directory
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %v", err)
		}
		apm[DIR_ROOT] = dirRoot
	}

	// Ensure required directories exist
	if err := apm.EnsureDirs(dirRoot); err != nil {
		return fmt.Errorf("failed to create required directories: %v", err)
	}

	// Validate paths
	if err := apm.ValidateWithOption(dirRoot); err != nil {
		return fmt.Errorf("path validation failed: %v", err)
	}

	// Store validated paths in the global sync.Map
	for key, value := range apm {
		globalsAppPathMap.Store(key, value)
	}

	return nil
}

// GetAppPathMap retrieves the global AppPathMap as a copy.
func GetAppPathMap() AppPathMap {
	apm := make(AppPathMap)
	globalsAppPathMap.Range(func(key, value interface{}) bool {
		apm[key.(AppPathKey)] = value.(string)
		return true
	})
	return apm
}

// GetAppPath retrieves a single path from the global sync.Map.
func GetAppPath(key AppPathKey) (string, bool) {
	val, ok := globalsAppPathMap.Load(key)
	if !ok {
		return "", false
	}
	return val.(string), true
}

// SetAppPath updates a single key in the global sync.Map.
func SetAppPath(key AppPathKey, value string) {
	globalsAppPathMap.Store(key, value)
}
