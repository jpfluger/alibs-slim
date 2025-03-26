package autils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AppPathMap is a map that associates AppPathKeys with their path values.
type AppPathMap map[AppPathKey]string

// SetPath sets the path value for a given AppPathKey.
func (apm *AppPathMap) SetPath(key AppPathKey, pathValue string) {
	(*apm)[key] = pathValue
}

// GetPath retrieves the path value associated with a given AppPathKey.
// If the key does not exist in the map, an empty string is returned.
func (apm AppPathMap) GetPath(key AppPathKey) string {
	val, ok := apm[key]
	if !ok {
		return "" // Key not found, return empty string
	}
	return val
}

// Validate checks if the paths associated with the keys in the map exist.
// It returns an error if any path does not exist or cannot be resolved.
func (apm *AppPathMap) Validate() error {
	for key, val := range *apm {
		if key.IsDir() {
			// Assume ResolveDirectory is a function that checks the existence of a directory.
			if _, err := ResolveDirectory(val); err != nil {
				return fmt.Errorf("failed to validate existence of directory %s; %v", key.String(), err)
			}
		} else {
			// Assume ResolveFile is a function that checks the existence of a file.
			if _, err := ResolveFile(val); err != nil {
				return fmt.Errorf("failed to validate existence of file %s; %v", key.String(), err)
			}
		}
	}
	return nil
}

// ValidateWithOption validates the existence of paths in the AppPathMap.
// If dirRoot is an empty string, it attempts to retrieve it from the AppPathMap using DIR_ROOT.
// Relative paths are resolved using dirRoot before checking existence.
// If a relative path is modified, it updates the map with the absolute path.
func (apm *AppPathMap) ValidateWithOption(dirRoot string) error {
	// If dirRoot is empty, attempt to retrieve it from the map
	if dirRoot == "" {
		if rootPath, exists := (*apm)[DIR_ROOT]; exists {
			dirRoot = rootPath
		}
	}

	for key, val := range *apm {
		val = strings.TrimSpace(val)
		if val == "" {
			return fmt.Errorf("invalid path for %s", key.String())
		}

		originalVal := val // Store the original value for comparison

		// If the path is relative, join it with dirRoot
		if !filepath.IsAbs(val) && dirRoot != "" {
			val = filepath.Join(dirRoot, val)
		}

		// Validate the resolved path
		if key.IsDir() {
			if _, err := ResolveDirectory(val); err != nil {
				return fmt.Errorf("failed to validate existence of directory %s: %v", key.String(), err)
			}
		} else if key.IsFile() {
			if _, err := ResolveFile(val); err != nil {
				return fmt.Errorf("failed to validate existence of file %s: %v", key.String(), err)
			}
		}

		if originalVal != val {
			(*apm)[key] = val // Update the map with the resolved absolute path
		}
	}

	return nil
}

// EnsureDirs ensures that all directories in the AppPathMap exist.
// It resolves relative paths using dirRoot before creation.
func (apm *AppPathMap) EnsureDirs(dirRoot string) error {
	if dirRoot == "" {
		if rootPath, exists := (*apm)[DIR_ROOT]; exists {
			dirRoot = rootPath
		}
	}

	for key, val := range *apm {
		if key.IsDir() {

			val = strings.TrimSpace(val)
			if val == "" {
				return fmt.Errorf("invalid path for %s", key.String())
			}

			// Resolve relative paths
			if !filepath.IsAbs(val) && dirRoot != "" {
				val = filepath.Join(dirRoot, val)
				(*apm)[key] = val // Store the resolved absolute path back in apm
			}

			// Check if directory exists before creating it
			if _, err := os.Stat(val); os.IsNotExist(err) {
				if err = os.MkdirAll(val, PATH_CHMOD_DIR_LIMIT); err != nil {
					return fmt.Errorf("failed to create directory %s: %v", val, err)
				}
			}
		}
	}
	return nil
}

// RequireWithOption ensures that all specified AppPathKeys exist in the map.
// If dirRoot is empty, it attempts to retrieve DIR_ROOT from the map.
// Returns an error if any required key is missing or has an empty value.
func (apm *AppPathMap) RequireWithOption(dirRoot string, keys ...AppPathKey) error {
	if dirRoot == "" {
		if rootPath, exists := (*apm)[DIR_ROOT]; exists {
			dirRoot = rootPath
		}
	}

	for _, key := range keys {
		if val, exists := (*apm)[key]; !exists || strings.TrimSpace(val) == "" {
			return fmt.Errorf("app path key %s does not exist", key)
		}
	}

	return apm.ValidateWithOption(dirRoot)
}

// RecreateDir removes and recreates a directory for the specified AppPathKey with the given file mode.
// If the directory exists, it is deleted before being recreated.
// Returns an error if the path is missing or if any operation fails.
func (apm *AppPathMap) RecreateDir(key AppPathKey, mode os.FileMode) error {
	val := apm.GetPath(key)
	if val == "" {
		return fmt.Errorf("path not found for %s", key.String())
	}

	if _, err := ResolveDirectory(val); err == nil {
		if err := os.RemoveAll(val); err != nil {
			return fmt.Errorf("failed to remove existing directory %s: %v", val, err)
		}
	}

	if err := os.MkdirAll(val, mode); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", val, err)
	}

	return nil
}
