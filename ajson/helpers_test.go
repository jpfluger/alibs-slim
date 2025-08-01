package ajson

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// TestWebConfig is a simplified version of WebConfig for testing.
type TestWebConfig struct {
	App      string `json:"app"`
	Version  string `json:"version"`
	Settings struct {
		Port  int  `json:"port"`
		Debug bool `json:"debug"`
	} `json:"settings"`
}

func TestUnmarshalFile(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		target   interface{}
		expected TestWebConfig
		wantErr  bool
		errMsg   string
	}{
		{
			name:   "valid_json",
			file:   "test_data/base.json",
			target: &TestWebConfig{},
			expected: TestWebConfig{
				App:     "myapp",
				Version: "1.0",
				Settings: struct {
					Port  int  `json:"port"`
					Debug bool `json:"debug"`
				}{Port: 8080, Debug: false},
			},
		},
		{
			name:    "empty_file_path",
			file:    "",
			target:  &TestWebConfig{},
			wantErr: true,
			errMsg:  "file path cannot be empty",
		},
		{
			name:    "nil_target",
			file:    "test_data/base.json",
			target:  nil,
			wantErr: true,
			errMsg:  "target cannot be nil",
		},
		{
			name:    "non_pointer_target",
			file:    "test_data/base.json",
			target:  TestWebConfig{},
			wantErr: true,
			errMsg:  "target must be a pointer",
		},
		{
			name:    "non_existent_file",
			file:    "test_data/nonexistent.json",
			target:  &TestWebConfig{},
			wantErr: true,
			errMsg:  "failed to read file",
		},
		{
			name:    "malformed_json",
			file:    "test_data/malformed.json",
			target:  &TestWebConfig{},
			wantErr: true,
			errMsg:  "failed to unmarshal file",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := UnmarshalFile(tc.file, tc.target)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got none", tc.errMsg)
				}
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Fatalf("expected error containing %q, got %q", tc.errMsg, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalFile failed: %v", err)
			}
			if !reflect.DeepEqual(tc.target, &tc.expected) {
				t.Errorf("unmarshaled result does not match expected\nGot: %+v\nExpected: %+v", tc.target, tc.expected)
			}
		})
	}
}

func TestMarshalIndent(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid_struct",
			input: TestWebConfig{
				App:     "myapp",
				Version: "1.0",
				Settings: struct {
					Port  int  `json:"port"`
					Debug bool `json:"debug"`
				}{Port: 8080, Debug: false},
			},
		},
		{
			name:    "nil_input",
			input:   nil,
			wantErr: true,
			errMsg:  "target cannot be nil",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := MarshalIndent(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got none", tc.errMsg)
				}
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Fatalf("expected error containing %q, got %q", tc.errMsg, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("MarshalIndent failed: %v", err)
			}
			var result TestWebConfig
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("failed to unmarshal result: %v", err)
			}
			if !reflect.DeepEqual(result, tc.input) {
				t.Errorf("marshaled result does not match input\nGot: %+v\nExpected: %+v", result, tc.input)
			}
		})
	}
}

func TestMarshalIndentToString(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		print   bool
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid_struct",
			input: TestWebConfig{
				App:     "myapp",
				Version: "1.0",
				Settings: struct {
					Port  int  `json:"port"`
					Debug bool `json:"debug"`
				}{Port: 8080, Debug: false},
			},
			print: false,
		},
		{
			name:    "nil_input",
			input:   nil,
			print:   false,
			wantErr: true,
			errMsg:  "target cannot be nil",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MarshalIndentToString(tc.input, tc.print)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got none", tc.errMsg)
				}
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Fatalf("expected error containing %q, got %q", tc.errMsg, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("MarshalIndentToString failed: %v", err)
			}
			var parsed TestWebConfig
			if err := json.Unmarshal([]byte(result), &parsed); err != nil {
				t.Fatalf("failed to unmarshal result: %v", err)
			}
			if !reflect.DeepEqual(parsed, tc.input) {
				t.Errorf("marshaled result does not match input\nGot: %+v\nExpected: %+v", parsed, tc.input)
			}
		})
	}
}

func TestMarshalIndentAndSaveFile(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		input    interface{}
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid_save",
			filepath: "test_output.json",
			input: TestWebConfig{
				App:     "myapp",
				Version: "1.0",
				Settings: struct {
					Port  int  `json:"port"`
					Debug bool `json:"debug"`
				}{Port: 8080, Debug: false},
			},
		},
		{
			name:     "empty_filepath",
			filepath: "",
			input:    TestWebConfig{},
			wantErr:  true,
			errMsg:   "file path cannot be empty",
		},
		{
			name:     "nil_input",
			filepath: "test_output.json",
			input:    nil,
			wantErr:  true,
			errMsg:   "target cannot be nil",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath := tc.filepath
			if filePath != "" {
				filePath = filepath.Join(tempDir, filePath)
			}

			err := MarshalIndentAndSaveFile(filePath, tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got none", tc.errMsg)
				}
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Fatalf("expected error containing %q, got %q", tc.errMsg, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("MarshalIndentAndSaveFile failed: %v", err)
			}

			data, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("failed to read saved file: %v", err)
			}
			var result TestWebConfig
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("failed to unmarshal saved file: %v", err)
			}
			if !reflect.DeepEqual(result, tc.input) {
				t.Errorf("saved result does not match input\nGot: %+v\nExpected: %+v", result, tc.input)
			}
		})
	}
}

func TestUnmarshalJSONToMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected map[string]interface{}
		wantErr  bool
		errMsg   string
	}{
		{
			name:  "valid_json",
			input: []byte(`{"app":"myapp","version":"1.0","settings":{"port":8080,"debug":false}}`),
			expected: map[string]interface{}{
				"app":     "myapp",
				"version": "1.0",
				"settings": map[string]interface{}{
					"port":  float64(8080),
					"debug": false,
				},
			},
		},
		{
			name:    "nil_input",
			input:   nil,
			wantErr: true,
			errMsg:  "input cannot be nil",
		},
		{
			name:    "empty_input",
			input:   []byte{},
			wantErr: true,
			errMsg:  "input is empty",
		},
		{
			name:    "invalid_json",
			input:   []byte(`{"app":"myapp`),
			wantErr: true,
			errMsg:  "failed to unmarshal JSON",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := UnmarshalJSONToMap(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got none", tc.errMsg)
				}
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Fatalf("expected error containing %q, got %q", tc.errMsg, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalJSONToMap failed: %v", err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("unmarshaled result does not match expected\nGot: %+v\nExpected: %+v", result, tc.expected)
			}
		})
	}
}

func TestConvertToMap(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected map[string]interface{}
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid_struct",
			input: TestWebConfig{
				App:     "myapp",
				Version: "1.0",
				Settings: struct {
					Port  int  `json:"port"`
					Debug bool `json:"debug"`
				}{Port: 8080, Debug: false},
			},
			expected: map[string]interface{}{
				"app":     "myapp",
				"version": "1.0",
				"settings": map[string]interface{}{
					"port":  float64(8080),
					"debug": false,
				},
			},
		},
		{
			name:     "nil_input",
			input:    nil,
			expected: make(map[string]interface{}),
		},
		{
			name:    "invalid_struct",
			input:   make(chan int),
			wantErr: true,
			errMsg:  "failed to marshal target",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ConvertToMap(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got none", tc.errMsg)
				}
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Fatalf("expected error containing %q, got %q", tc.errMsg, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ConvertToMap failed: %v", err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("converted result does not match expected\nGot: %+v\nExpected: %+v", result, tc.expected)
			}
		})
	}
}

func TestConvertFromMap(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		target   interface{}
		expected TestWebConfig
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid_map",
			input: map[string]interface{}{
				"app":     "myapp",
				"version": "1.0",
				"settings": map[string]interface{}{
					"port":  float64(8080),
					"debug": false,
				},
			},
			target: &TestWebConfig{},
			expected: TestWebConfig{
				App:     "myapp",
				Version: "1.0",
				Settings: struct {
					Port  int  `json:"port"`
					Debug bool `json:"debug"`
				}{Port: 8080, Debug: false},
			},
		},
		{
			name:    "nil_map",
			input:   nil,
			target:  &TestWebConfig{},
			wantErr: true,
			errMsg:  "data cannot be nil",
		},
		{
			name:    "nil_target",
			input:   map[string]interface{}{"app": "myapp"},
			target:  nil,
			wantErr: true,
			errMsg:  "target cannot be nil",
		},
		{
			name:    "non_pointer_target",
			input:   map[string]interface{}{"app": "myapp"},
			target:  TestWebConfig{},
			wantErr: true,
			errMsg:  "target must be a pointer",
		},
		{
			name:    "invalid_map",
			input:   map[string]interface{}{"app": make(chan int)},
			target:  &TestWebConfig{},
			wantErr: true,
			errMsg:  "failed to marshal map",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ConvertFromMap(tc.input, tc.target)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got none", tc.errMsg)
				}
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Fatalf("expected error containing %q, got %q", tc.errMsg, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ConvertFromMap failed: %v", err)
			}
			if !reflect.DeepEqual(tc.target, &tc.expected) {
				t.Errorf("converted result does not match expected\nGot: %+v\nExpected: %+v", tc.target, tc.expected)
			}
		})
	}
}

func TestUnmarshalRawMessage(t *testing.T) {
	tests := []struct {
		name     string
		raw      json.RawMessage
		target   interface{}
		expected TestWebConfig
		wantErr  bool
		errMsg   string
	}{
		{
			name:   "valid_raw_message",
			raw:    json.RawMessage(`{"app":"myapp","version":"1.0","settings":{"port":8080,"debug":false}}`),
			target: &TestWebConfig{},
			expected: TestWebConfig{
				App:     "myapp",
				Version: "1.0",
				Settings: struct {
					Port  int  `json:"port"`
					Debug bool `json:"debug"`
				}{Port: 8080, Debug: false},
			},
		},
		{
			name:    "nil_raw_message",
			raw:     nil,
			target:  &TestWebConfig{},
			wantErr: true,
			errMsg:  "raw message cannot be nil",
		},
		{
			name:    "nil_target",
			raw:     json.RawMessage(`{"app":"myapp"}`),
			target:  nil,
			wantErr: true,
			errMsg:  "target cannot be nil",
		},
		{
			name:    "non_pointer_target",
			raw:     json.RawMessage(`{"app":"myapp"}`),
			target:  TestWebConfig{},
			wantErr: true,
			errMsg:  "target must be a pointer",
		},
		{
			name:    "invalid_raw_message",
			raw:     json.RawMessage(`{"app":"myapp`),
			target:  &TestWebConfig{},
			wantErr: true,
			errMsg:  "failed to unmarshal raw message",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := UnmarshalRawMessage(tc.raw, tc.target)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got none", tc.errMsg)
				}
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Fatalf("expected error containing %q, got %q", tc.errMsg, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalRawMessage failed: %v", err)
			}
			if !reflect.DeepEqual(tc.target, &tc.expected) {
				t.Errorf("unmarshaled result does not match expected\nGot: %+v\nExpected: %+v", tc.target, tc.expected)
			}
		})
	}
}
