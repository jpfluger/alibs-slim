package ajson

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/autils"
	"os"
	"strings"
)

// UnmarshalFromFile reads a JSON file and unmarshals its content into the target interface.
func UnmarshalFromFile(targetFile string, targetInterface interface{}) error {
	if targetInterface == nil {
		return fmt.Errorf("target interface is nil")
	}
	b, err := os.ReadFile(targetFile)
	if err != nil {
		return fmt.Errorf("failed to read target file: %v", err)
	}
	return UnmarshalBytes(b, targetInterface)
}

// UnmarshalBytes unmarshals a byte slice into the target interface.
func UnmarshalBytes(b []byte, targetInterface interface{}) error {
	if b == nil {
		return fmt.Errorf("bytes is nil")
	}
	if targetInterface == nil {
		return fmt.Errorf("target interface is nil")
	}
	if err := json.Unmarshal(b, targetInterface); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}
	return nil
}

// UnmarshalFile reads a JSON file and unmarshals its content into the target interface.
func UnmarshalFile(fpath string, target interface{}) error {
	if target == nil {
		return fmt.Errorf("target interface is nil")
	}

	fpath = strings.TrimSpace(fpath)
	if fpath == "" {
		return fmt.Errorf("file path is empty")
	}

	b, err := os.ReadFile(fpath)
	if err != nil {
		return fmt.Errorf("failed to read file path: %v", err)
	}
	if err = json.Unmarshal(b, target); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}

	return nil
}

// Marshal serializes the target interface into a JSON byte slice.
func Marshal(targetInterface interface{}) ([]byte, error) {
	if targetInterface == nil {
		return nil, fmt.Errorf("target interface is nil")
	}
	return json.Marshal(targetInterface)
}

// MarshalIndent serializes the target interface into a pretty-printed JSON byte slice.
func MarshalIndent(targetInterface interface{}) ([]byte, error) {
	if targetInterface == nil {
		return nil, fmt.Errorf("target interface is nil")
	}
	return json.MarshalIndent(targetInterface, "", "  ")
}

// MarshalIndentToString serializes the target interface into a pretty-printed JSON string.
func MarshalIndentToString(targetInterface interface{}) (string, error) {
	b, err := MarshalIndent(targetInterface)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// MarshalIndentWithPrint serializes the target interface into a pretty-printed JSON string and prints it.
func MarshalIndentWithPrint(targetInterface interface{}) (string, error) {
	s, err := MarshalIndentToString(targetInterface)
	if err != nil {
		return "", err
	}
	fmt.Println(s)
	return s, nil
}

// MarshalIndentAndSaveFile serializes the target interface into a pretty-printed JSON string
// and saves it to the specified file path.
func MarshalIndentAndSaveFile(filepath string, targetInterface interface{}) error {
	// Serialize the interface to a JSON string with indentation
	data, err := json.MarshalIndent(targetInterface, "", "  ")
	if err != nil {
		return err
	}

	// Save the JSON string to the specified file path
	if err = os.WriteFile(filepath, data, autils.PATH_CHMOD_FILE); err != nil {
		return err
	}

	return nil
}

// ConvertToMapOfInterfaces converts the target interface into a map of interfaces by marshaling and unmarshaling.
func ConvertToMapOfInterfaces(target interface{}) (map[string]interface{}, error) {
	var data map[string]interface{}
	if target == nil {
		return data, nil
	}

	b, err := json.Marshal(target)
	if err != nil {
		return data, fmt.Errorf("cannot marshal json: %v", err)
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return data, fmt.Errorf("cannot unmarshal json: %v", err)
	}

	return data, nil
}

// ConvertFromMapOfInterfaces converts a map of interfaces into the target interface by marshaling and unmarshaling.
func ConvertFromMapOfInterfaces(data map[string]interface{}, target interface{}) error {
	if data == nil {
		return fmt.Errorf("data is nil")
	}
	if target == nil {
		return fmt.Errorf("target is nil")
	}

	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %v", err)
	}

	if err := json.Unmarshal(b, target); err != nil {
		return fmt.Errorf("failed to unmarshal json: %v", err)
	}

	return nil
}

// ConvertMapToJSON converts a map into a JSON byte slice.
func ConvertMapToJSON(input map[string]interface{}) ([]byte, error) {
	return json.Marshal(input)
}

// ConvertJSONToMap converts a JSON byte slice into a map.
func ConvertJSONToMap(input []byte) (map[string]interface{}, error) {
	var output map[string]interface{}
	if len(input) == 0 {
		return output, fmt.Errorf("input is empty")
	}
	if err := json.Unmarshal(input, &output); err != nil {
		return output, fmt.Errorf("failed to unmarshal JSON to map: %v", err)
	}
	return output, nil
}
