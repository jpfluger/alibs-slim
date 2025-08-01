package ajson

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/jpfluger/alibs-slim/autils"
)

// UnmarshalFile loads a JSON file into the target struct or map.
// The target must be a non-nil pointer.
// Returns an error if the file cannot be read or unmarshaled.
func UnmarshalFile(file string, target interface{}) error {
	if file == "" {
		return fmt.Errorf("file path cannot be empty")
	}
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	if reflect.ValueOf(target).Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", file, err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal file %q into target: %w", file, err)
	}
	return nil
}

// MarshalIndent serializes the target interface into a pretty-printed JSON byte slice.
// Returns an error if marshaling fails.
func MarshalIndent(target interface{}) ([]byte, error) {
	if target == nil {
		return nil, fmt.Errorf("target cannot be nil")
	}
	data, err := json.MarshalIndent(target, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal target: %w", err)
	}
	return data, nil
}

// MarshalIndentToString serializes the target interface into a pretty-printed JSON string.
// If print is true, the JSON string is printed to stdout.
// Returns the JSON string and an error if marshaling fails.
func MarshalIndentToString(target interface{}, print bool) (string, error) {
	if target == nil {
		return "", fmt.Errorf("target cannot be nil")
	}
	data, err := json.MarshalIndent(target, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON string: %w", err)
	}
	result := string(data)
	if print {
		fmt.Println(result)
	}
	return result, nil
}

// MarshalIndentAndSaveFile serializes the target interface into a pretty-printed JSON string
// and saves it to the specified file path.
// Returns an error if marshaling or file writing fails.
func MarshalIndentAndSaveFile(filepath string, target interface{}) error {
	if filepath == "" {
		return fmt.Errorf("file path cannot be empty")
	}
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}

	data, err := json.MarshalIndent(target, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal target: %w", err)
	}

	if err := os.WriteFile(filepath, data, autils.PATH_CHMOD_FILE); err != nil {
		return fmt.Errorf("failed to write to file %q: %w", filepath, err)
	}
	return nil
}

// UnmarshalJSONToMap converts a JSON byte slice into a map[string]interface{}.
// Returns an error if the input is nil, empty, or invalid JSON.
func UnmarshalJSONToMap(input []byte) (map[string]interface{}, error) {
	if input == nil {
		return nil, fmt.Errorf("input cannot be nil")
	}
	if len(input) == 0 {
		return nil, fmt.Errorf("input is empty")
	}
	var output map[string]interface{}
	if err := json.Unmarshal(input, &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to map: %w", err)
	}
	return output, nil
}

// ConvertToMap converts the target interface into a map[string]interface{} by marshaling and unmarshaling.
// Returns an empty map if the target is nil, or an error if conversion fails.
func ConvertToMap(target interface{}) (map[string]interface{}, error) {
	if target == nil {
		return make(map[string]interface{}), nil
	}
	b, err := json.Marshal(target)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal target: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to map: %w", err)
	}
	return result, nil
}

// ConvertFromMap converts a map[string]interface{} into the target struct or map.
// The target must be a non-nil pointer.
// Returns an error if the map or target is nil, or if conversion fails.
func ConvertFromMap(data map[string]interface{}, target interface{}) error {
	if data == nil {
		return fmt.Errorf("data cannot be nil")
	}
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	if reflect.ValueOf(target).Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal map: %w", err)
	}
	if err := json.Unmarshal(b, target); err != nil {
		return fmt.Errorf("failed to unmarshal to target: %w", err)
	}
	return nil
}

// UnmarshalRawMessage unmarshals a json.RawMessage into the target struct or map.
// The target must be a non-nil pointer.
// Returns an error if unmarshaling fails.
func UnmarshalRawMessage(raw json.RawMessage, target interface{}) error {
	if raw == nil {
		return fmt.Errorf("raw message cannot be nil")
	}
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	if reflect.ValueOf(target).Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return fmt.Errorf("failed to unmarshal raw message: %w", err)
	}
	return nil
}
