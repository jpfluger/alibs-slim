package ajson

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"dario.cat/mergo"
	"github.com/hjson/hjson-go/v4"
	"github.com/jpfluger/alibs-slim/autils"
)

// MergeOptions defines configuration for merging JSON or HJSON files.
type MergeOptions struct {
	Files         []string `json:"files,omitempty"`         // List of file paths to merge, in order
	UseHJSON      bool     `json:"useHJSON,omitempty"`      // Whether files are in HJSON format
	StripComments bool     `json:"stripComments,omitempty"` // If true, strip comments from JSON or HJSON
}

// commentRegex is a precompiled regular expression for removing single-line (//) and multi-line (/* */) comments.
var commentRegex = regexp.MustCompile(`(?m)//.*$|/\*[\s\S]*?\*/`)

// StripComments removes single-line (//) and multi-line (/* */) comments from JSON or HJSON text.
// It returns the cleaned byte slice with comments removed.
func StripComments(input []byte) []byte {
	return commentRegex.ReplaceAll(input, []byte{})
}

// loadFileToMerge loads a single config file into a map, respecting the provided options.
// It validates the file path, reads the file, and parses it as JSON or HJSON.
// Returns the parsed map or an error if the file cannot be read or parsed.
func loadFileToMerge(path string, useHjson bool, stripComments bool) (map[string]interface{}, error) {
	if path == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}
	if _, err := autils.ResolveFile(path); err != nil {
		return nil, fmt.Errorf("invalid file path %q: %w", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", path, err)
	}

	if stripComments {
		data = StripComments(data)
	}

	var result map[string]interface{}
	if useHjson {
		if err := hjson.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("failed to parse HJSON file %q: %w", path, err)
		}
	} else {
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("failed to parse JSON file %q: %w", path, err)
		}
	}
	return result, nil
}

// MergeConfigs merges multiple configuration files into a single map based on the provided options.
// Files are merged in order, with later files overriding earlier ones.
// Returns the merged map or an error if any file fails to load or merge.
func MergeConfigs(opts MergeOptions) (map[string]interface{}, error) {
	if opts.Files == nil || len(opts.Files) == 0 {
		return nil, fmt.Errorf("no config files provided")
	}

	final := make(map[string]interface{})
	for _, file := range opts.Files {
		current, err := loadFileToMerge(file, opts.UseHJSON, opts.StripComments)
		if err != nil {
			return nil, err
		}
		if err := mergo.Merge(&final, current, mergo.WithOverride); err != nil {
			return nil, fmt.Errorf("failed to merge config file %q: %w", file, err)
		}
	}
	return final, nil
}

// MergeConfigsInto merges multiple configuration files into a target struct or map.
// The target must be a non-nil pointer to a struct or map.
// Returns an error if the merge or target population fails.
func MergeConfigsInto(target interface{}, opts MergeOptions) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	if reflect.ValueOf(target).Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	mergedMap, err := MergeConfigs(opts)
	if err != nil {
		return err
	}

	mergedBytes, err := json.Marshal(mergedMap)
	if err != nil {
		return fmt.Errorf("failed to marshal merged map: %w", err)
	}

	if err := json.Unmarshal(mergedBytes, target); err != nil {
		return fmt.Errorf("failed to unmarshal into target: %w", err)
	}
	return nil
}

// MergeConfigsSaveAs merges configuration files and writes the result to a JSON file.
// The output is formatted with indentation for readability.
// Returns an error if the merge or file write fails.
func MergeConfigsSaveAs(saveFile string, opts MergeOptions) error {
	if saveFile == "" {
		return fmt.Errorf("save file path cannot be empty")
	}

	merged, err := MergeConfigs(opts)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal merged config: %w", err)
	}

	if err := os.WriteFile(saveFile, data, autils.PATH_CHMOD_FILE); err != nil {
		return fmt.Errorf("failed to write merged config to %q: %w", saveFile, err)
	}
	return nil
}

// MergeFilesIntoMap merges configuration files into a single map.
// If writeToFile is non-empty, the merged map is saved as formatted JSON to the specified file.
// Returns the merged map and an error if any step fails.
func MergeFilesIntoMap(files []string, writeToFile string) (map[string]interface{}, error) {
	if files == nil || len(files) == 0 {
		return nil, fmt.Errorf("no config files provided for merge")
	}

	mergeOpts := MergeOptions{
		Files:         files,
		UseHJSON:      false, // Default to JSON
		StripComments: true,  // Default to stripping comments
	}

	merged, err := MergeConfigs(mergeOpts)
	if err != nil {
		return nil, err
	}

	if writeToFile != "" {
		data, err := json.MarshalIndent(merged, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal merged config: %w", err)
		}
		if err := os.WriteFile(writeToFile, data, autils.PATH_CHMOD_FILE); err != nil {
			return nil, fmt.Errorf("failed to save merged config to %q: %w", writeToFile, err)
		}
	}

	return merged, nil
}

// LoadFileIntoElseMerge tries to load a single file into the target if it exists; otherwise, it merges multiple files.
// The target must be a non-nil pointer to a struct or map.
// Returns an error if the load, merge, or unmarshal fails.
func LoadFileIntoElseMerge(loadFile string, files []string, target interface{}) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	if reflect.ValueOf(target).Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	loadFile = strings.TrimSpace(loadFile)
	if loadFile != "" {
		if _, err := autils.ResolveFile(loadFile); err == nil {
			return UnmarshalFile(loadFile, target)
		}
	}

	if files == nil || len(files) == 0 {
		return fmt.Errorf("no files specified for merge")
	}

	mergeOpts := MergeOptions{
		Files:         files,
		UseHJSON:      true, // Default to HJSON for fallback merge
		StripComments: false,
	}

	if err := MergeConfigsInto(target, mergeOpts); err != nil {
		return fmt.Errorf("failed to merge files into target: %w", err)
	}
	return nil
}
