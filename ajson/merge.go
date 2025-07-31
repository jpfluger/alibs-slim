package ajson

import (
	"dario.cat/mergo"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hjson/hjson-go/v4"
	"github.com/jpfluger/alibs-slim/autils"
	"os"
	"regexp"
	"strings"
)

// MergeOptions defines how the JSON configs should be processed.
type MergeOptions struct {
	Files         []string `json:"files,omitempty"`         // List of file paths to merge, in order
	UseHJSON      bool     `json:"useHJSON,omitempty"`      // Whether files are in HJSON format
	StripComments bool     `json:"stripComments,omitempty"` // If true, strip comments from JSON
}

// StripComments removes both // and /* */ comments from JSON text.
func StripComments(input []byte) []byte {
	re := regexp.MustCompile(`(?m)//.*$|/\*[\s\S]*?\*/`)
	return re.ReplaceAll(input, []byte{})
}

// loadFileToMerge loads a single config file into a map, respecting options.
func loadFileToMerge(path string, useHjson bool, stripComments bool) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if useHjson {
		err = hjson.Unmarshal(data, &result)
	} else {
		if stripComments {
			data = StripComments(data)
		}
		err = json.Unmarshal(data, &result)
	}
	return result, err
}

// MergeConfigs merges multiple config files based on the provided options.
func MergeConfigs(opts MergeOptions) (map[string]interface{}, error) {
	if len(opts.Files) == 0 {
		return nil, errors.New("no config files provided")
	}

	final := make(map[string]interface{})

	for _, file := range opts.Files {
		current, err := loadFileToMerge(file, opts.UseHJSON, opts.StripComments)
		if err != nil {
			return nil, errors.New("failed to load config file: " + file + ": " + err.Error())
		}

		err = mergo.Merge(&final, current, mergo.WithOverride)
		if err != nil {
			return nil, errors.New("failed to merge config file: " + file + ": " + err.Error())
		}
	}

	return final, nil
}

// MergeConfigsInto loads and merges config files into a typed target (e.g., struct pointer).
// It calls MergeConfigs, then decodes the result into the target interface.
func MergeConfigsInto(target interface{}, opts MergeOptions) error {
	if target == nil {
		return errors.New("target cannot be nil")
	}

	mergedMap, err := MergeConfigs(opts)
	if err != nil {
		return err
	}

	mergedBytes, err := json.Marshal(mergedMap)
	if err != nil {
		return errors.New("failed to marshal merged map: " + err.Error())
	}

	if err = json.Unmarshal(mergedBytes, target); err != nil {
		return errors.New("failed to unmarshal into target: " + err.Error())
	}

	return nil
}

// MergeConfigsSaveAs merges config files and writes the result to a JSON file at saveFile.
func MergeConfigsSaveAs(saveFile string, opts MergeOptions) error {
	merged, err := MergeConfigs(opts)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(saveFile, data, autils.PATH_CHMOD_FILE)
}

// MergeFilesIntoWithSave merges multiple configuration files into the given target struct.
// The resulting structure is marshaled into a formatted JSON string and returned.
// If writeToFile is non-empty, the JSON will also be saved to the specified file path.
//
// Parameters:
//   - files: list of configuration file paths (must be non-empty)
//   - target: a pointer to the struct to unmarshal the merged config into
//   - writeToFile: optional file path to write the final merged JSON output
//
// Returns:
//   - A formatted JSON string of the merged config
//   - An error if any step in the process fails (merge, marshal, or file write)
func MergeFilesIntoWithSave(files []string, target interface{}, writeToFile string) (string, error) {
	if files == nil || len(files) == 0 {
		return "", fmt.Errorf("no files provided for merge")
	}

	mergeOpts := MergeOptions{
		Files:         files,
		UseHJSON:      false,
		StripComments: true,
	}

	// Perform the merge into the target
	if err := MergeConfigsInto(target, mergeOpts); err != nil {
		return "", fmt.Errorf("failed to merge files into target: %w", err)
	}

	// Marshal the result to string
	result, err := MarshalIndentToString(target)
	if err != nil {
		return "", fmt.Errorf("failed to marshal merged file: %w", err)
	}

	// Optionally write to file
	if strings.TrimSpace(writeToFile) != "" {
		if err = os.WriteFile(writeToFile, []byte(result), autils.PATH_CHMOD_FILE); err != nil {
			return "", fmt.Errorf("failed to save merged file to %q: %w", writeToFile, err)
		}
	}

	return result, nil
}

// MergeFilesIntoMap merges multiple configuration files into a single map[string]interface{}.
// If writeToFile is non-empty, the merged map will also be saved as formatted JSON to the given file path.
//
// Parameters:
//   - files: list of configuration file paths to merge
//   - writeToFile: optional output file path for saving the merged config
//
// Returns:
//   - The merged configuration as a map[string]interface{}
//   - An error if any part of the process fails (e.g., merge or file write)
func MergeFilesIntoMap(files []string, writeToFile string) (map[string]interface{}, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no config files provided for merge")
	}

	mergeOpts := MergeOptions{
		Files:         files,
		UseHJSON:      false,
		StripComments: true,
	}

	merged, err := MergeConfigs(mergeOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to merge config files: %w", err)
	}

	if strings.TrimSpace(writeToFile) != "" {
		data, err := json.MarshalIndent(merged, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal merged config: %w", err)
		}
		if err := os.WriteFile(writeToFile, data, autils.PATH_CHMOD_FILE); err != nil {
			return nil, fmt.Errorf("failed to save merged config to file %q: %w", writeToFile, err)
		}
	}

	return merged, nil
}
