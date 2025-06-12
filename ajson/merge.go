package ajson

import (
	"dario.cat/mergo"
	"encoding/json"
	"errors"
	"github.com/hjson/hjson-go/v4"
	"github.com/jpfluger/alibs-slim/autils"
	"os"
	"regexp"
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
