package autils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// AppPathMap is a map that associates AppPathKeys with their path values.
type AppPathMap map[AppPathKey]string

// NewAppPathMap creates a new AppPathMap.
func NewAppPathMap() AppPathMap {
	return make(AppPathMap)
}

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

// ValidateDir only checks those keys identified as a directory exist.
// It returns an error if any path does not exist or cannot be resolved.
func (apm *AppPathMap) ValidateDir() error {
	for key, val := range *apm {
		if key.IsDir() {
			// Assume ResolveDirectory is a function that checks the existence of a directory.
			if _, err := ResolveDirectory(val); err != nil {
				return fmt.Errorf("failed to validate existence of directory %s; %v", key.String(), err)
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
// Skips creation for paths that are not under dirRoot (if dirRoot is set).
func (apm *AppPathMap) EnsureDirs(dirRoot string) error {
	if dirRoot == "" {
		if rootPath, exists := (*apm)[DIR_ROOT]; exists {
			dirRoot = rootPath
		}
	}

	for _, key := range apm.SortedKeys() {
		val := strings.TrimSpace((*apm)[key])
		if !key.IsDir() {
			continue // Only process directories here
		}

		if val == "" {
			return fmt.Errorf("invalid path for %s", key.String())
		}

		// Resolve relative path.
		if !filepath.IsAbs(val) && dirRoot != "" {
			val = filepath.Join(dirRoot, val)
			(*apm)[key] = val
		}

		// Skip if not under dirRoot
		if dirRoot != "" {
			under, err := IsUnderRoot(dirRoot, val)
			if err != nil {
				return fmt.Errorf("failed to check if %s is under root: %w", val, err)
			}
			if !under {
				continue // Skip creation for paths outside DIR_ROOT
			}
		}

		if _, err := os.Stat(val); os.IsNotExist(err) {
			if err = os.MkdirAll(val, PATH_CHMOD_DIR_LIMIT); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", val, err)
			}
		}
	}
	return nil
}

// AutoSetupIfRootEmpty automatically sets up default paths and creates directories if DIR_ROOT is empty.
// It merges the provided defaults into the current AppPathMap (overrides win), resolves relative paths,
// and ensures all directories are created. If DIR_ROOT is not empty or does not exist, it skips setup.
// Defaults should use relative paths (e.g., "data" for DIR_DATA, "data/logs" for DIR_LOGS).
func (apm *AppPathMap) AutoSetupIfRootEmpty(defaults AppPathMap) error {
	dirRoot := apm.GetPath(DIR_ROOT)
	if dirRoot == "" {
		return fmt.Errorf("DIR_ROOT is not set")
	}

	// Check if DIR_ROOT is empty
	isEmpty, err := IsDirEmpty(dirRoot)
	if err != nil {
		return err
	}
	if !isEmpty {
		return nil // Not empty, skip auto-setup
	}

	// Merge defaults, resolving relatives to DIR_ROOT
	merged := apm.MergeAbs(dirRoot, defaults)

	// Update the current map with the merged values
	for k, v := range merged {
		(*apm)[k] = v
	}

	// Ensure all directories are created
	if err := apm.EnsureDirs(dirRoot); err != nil {
		return fmt.Errorf("failed to ensure directories during auto-setup: %w", err)
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

// IsRelative checks if the path value associated with the given AppPathKey is a non-empty relative path.
// It returns true if the path is not empty and is not an absolute path.
// Returns false if the path is empty or already absolute.
func (apm *AppPathMap) IsRelative(key AppPathKey) bool {
	val := apm.GetPath(key)
	return val != "" && !filepath.IsAbs(val)
}

// SortedKeys returns the AppPathKeys in sorted order.
func (apm *AppPathMap) SortedKeys() []AppPathKey {
	var keys []AppPathKey
	for key := range *apm {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})
	return keys
}

// Clone returns a deep copy of the AppPathMap.
// It preserves all keys and values in a new map instance.
func (m AppPathMap) Clone() AppPathMap {
	clone := make(AppPathMap, len(m))
	for k, v := range m {
		if k.IsEmpty() {
			continue
		}
		clone[k] = v
	}
	return clone
}

// Merge returns a new AppPathMap resulting from merging the current map (`base`)
// with one or more overrides. If a key exists in both, the override wins.
func (base AppPathMap) Merge(overrides ...AppPathMap) AppPathMap {
	merged := base.Clone()
	for _, override := range overrides {
		for k, v := range override {
			if k.IsEmpty() {
				continue
			}
			merged[k] = v
		}
	}
	return merged
}

// MergeAbs merges base paths with overrides, ensuring all resulting values are absolute.
// Any relative paths in the override will be joined with `baseDir`.
func (base AppPathMap) MergeAbs(baseDir string, overrides ...AppPathMap) AppPathMap {
	merged := base.Clone()
	for _, override := range overrides {
		if override == nil {
			continue
		}
		for k, v := range override {
			if k.IsEmpty() || strings.TrimSpace(v) == "" {
				continue
			}
			if !filepath.IsAbs(v) {
				v = filepath.Join(baseDir, v)
			}
			merged[k] = filepath.Clean(v)
		}
	}
	return merged
}

// ResolveAllRelatives resolves all relative paths in the AppPathMap to absolute paths
// using DIR_ROOT as the base directory. It mutates the map in place by updating
// relative path values to their absolute equivalents. If DIR_ROOT is not set, it
// returns an error. Empty or already absolute paths are skipped without changes.
func (apm *AppPathMap) ResolveAllRelatives() error {
	baseDir := apm.GetPath(DIR_ROOT)
	if baseDir == "" {
		return fmt.Errorf("DIR_ROOT is required to resolve relative paths")
	}
	for key, val := range *apm {
		// Skip if the value is empty or already absolute
		if val == "" || filepath.IsAbs(val) {
			continue
		}
		// Resolve relative path
		absPath := filepath.Join(baseDir, val)
		(*apm)[key] = filepath.Clean(absPath)
	}
	return nil
}
