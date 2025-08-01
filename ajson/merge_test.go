package ajson

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// WebConfig defines a test struct for typed configuration merging.
type WebConfig struct {
	Services struct {
		Web struct {
			Port     int             `json:"port"`
			Timeout  int             `json:"timeout"`
			Tags     []string        `json:"tags"`
			Features map[string]bool `json:"features"`
		} `json:"web"`
	} `json:"services"`
	Global struct {
		Logging struct {
			Level   string `json:"level"`
			Enabled bool   `json:"enabled"`
		} `json:"logging"`
		Theme string `json:"theme"`
	} `json:"global"`
	Metadata struct {
		Version string `json:"version"`
		Build   int    `json:"build"`
	} `json:"metadata"`
}

// loadJSONMapFromFile loads a JSON file into a map[string]interface{}.
func loadJSONMapFromFile(t *testing.T, file string) (map[string]interface{}, error) {
	t.Helper()
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// loadWebConfigFromFile loads a JSON file into a WebConfig struct.
func loadWebConfigFromFile(t *testing.T, file string) (*WebConfig, error) {
	t.Helper()
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var cfg WebConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func TestMergeConfigs(t *testing.T) {
	tests := []struct {
		name         string
		options      MergeOptions
		expectedFile string
		wantErr      bool
		errMsg       string
	}{
		{
			name: "basic_merge",
			options: MergeOptions{
				Files:         []string{"test_data/base.json", "test_data/override.json"},
				UseHJSON:      false,
				StripComments: true,
			},
			expectedFile: "test_data/expected_merged.json",
		},
		{
			name: "commented_json_merge",
			options: MergeOptions{
				Files:         []string{"test_data/base_with_comments.json", "test_data/override.json"},
				UseHJSON:      false,
				StripComments: true,
			},
			expectedFile: "test_data/expected_merged.json",
		},
		{
			name: "three_file_merge",
			options: MergeOptions{
				Files:         []string{"test_data/base.json", "test_data/override.json", "test_data/env.json"},
				UseHJSON:      false,
				StripComments: true,
			},
			expectedFile: "test_data/expected_merged_env.json",
		},
		{
			name: "hjson_merge",
			options: MergeOptions{
				Files:         []string{"test_data/base.hjson", "test_data/override.hjson"},
				UseHJSON:      true,
				StripComments: false,
			},
			expectedFile: "test_data/expected_merged_hjson.json",
		},
		{
			name: "deep_merge",
			options: MergeOptions{
				Files:         []string{"test_data/base_deep.json", "test_data/override_deep.json"},
				UseHJSON:      false,
				StripComments: false,
			},
			expectedFile: "test_data/expected_deep_merged.json",
		},
		{
			name: "chained_merge_with_third_override",
			options: MergeOptions{
				Files:         []string{"test_data/base_deep.json", "test_data/override_deep.json", "test_data/override_deep2.json"},
				UseHJSON:      false,
				StripComments: false,
			},
			expectedFile: "test_data/expected_deep2_merged.json",
		},
		{
			name:    "empty_files",
			options: MergeOptions{Files: []string{}, UseHJSON: false, StripComments: true},
			wantErr: true,
			errMsg:  "no config files provided",
		},
		{
			name: "invalid_file",
			options: MergeOptions{
				Files:         []string{"test_data/nonexistent.json"},
				UseHJSON:      false,
				StripComments: true,
			},
			wantErr: true,
			errMsg:  "invalid file path",
		},
		{
			name: "malformed_json",
			options: MergeOptions{
				Files:         []string{"test_data/malformed.json"},
				UseHJSON:      false,
				StripComments: true,
			},
			wantErr: true,
			errMsg:  "failed to parse JSON file",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			merged, err := MergeConfigs(tc.options)
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
				t.Fatalf("merge failed: %v", err)
			}

			expected, err := loadJSONMapFromFile(t, tc.expectedFile)
			if err != nil {
				t.Fatalf("failed to load expected file: %v", err)
			}

			if !reflect.DeepEqual(merged, expected) {
				mergedBytes, _ := json.MarshalIndent(merged, "", "  ")
				expectedBytes, _ := json.MarshalIndent(expected, "", "  ")
				t.Errorf("merged result does not match expected\nGot:\n%s\n\nExpected:\n%s\n",
					string(mergedBytes), string(expectedBytes))
			}
		})
	}
}

func TestMergeConfigsInto_TypedStruct(t *testing.T) {
	tests := []struct {
		name         string
		options      MergeOptions
		expectedFile string
		wantErr      bool
		errMsg       string
	}{
		{
			name: "valid_merge",
			options: MergeOptions{
				Files:         []string{"test_data/base_deep.json", "test_data/override_deep.json", "test_data/override_deep2.json"},
				UseHJSON:      false,
				StripComments: false,
			},
			expectedFile: "test_data/expected_deep2_merged.json",
		},
		{
			name:    "nil_target",
			options: MergeOptions{Files: []string{"test_data/base.json"}},
			wantErr: true,
			errMsg:  "target cannot be nil",
		},
		{
			name:    "non_pointer_target",
			options: MergeOptions{Files: []string{"test_data/base.json"}},
			wantErr: true,
			errMsg:  "target must be a pointer",
		},
		{
			name: "invalid_file",
			options: MergeOptions{
				Files:         []string{"test_data/nonexistent.json"},
				UseHJSON:      false,
				StripComments: true,
			},
			wantErr: true,
			errMsg:  "invalid file path",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var cfg WebConfig
			var target interface{}
			if tc.name == "non_pointer_target" {
				target = WebConfig{}
			} else {
				target = &cfg
			}
			if tc.name == "nil_target" {
				target = nil
			}

			err := MergeConfigsInto(target, tc.options)
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
				t.Fatalf("MergeConfigsInto failed: %v", err)
			}

			expected, err := loadWebConfigFromFile(t, tc.expectedFile)
			if err != nil {
				t.Fatalf("failed to load expected file: %v", err)
			}

			if !reflect.DeepEqual(cfg, *expected) {
				cfgBytes, _ := json.MarshalIndent(cfg, "", "  ")
				expectedBytes, _ := json.MarshalIndent(expected, "", "  ")
				t.Errorf("merged config does not match expected\nGot:\n%s\n\nExpected:\n%s\n",
					string(cfgBytes), string(expectedBytes))
			}
		})
	}
}

func TestMergeConfigsSaveAs(t *testing.T) {
	tests := []struct {
		name         string
		savePath     string
		options      MergeOptions
		expectedFile string
		wantErr      bool
		errMsg       string
	}{
		{
			name:     "valid_save",
			savePath: "generated_save_output.json",
			options: MergeOptions{
				Files:         []string{"test_data/base_deep.json", "test_data/override_deep.json", "test_data/override_deep2.json"},
				UseHJSON:      false,
				StripComments: false,
			},
			expectedFile: "test_data/expected_deep2_merged.json",
		},
		{
			name:     "empty_save_path",
			savePath: "",
			options:  MergeOptions{Files: []string{"test_data/base.json"}},
			wantErr:  true,
			errMsg:   "save file path cannot be empty",
		},
		{
			name:     "invalid_file",
			savePath: "generated_save_output.json",
			options: MergeOptions{
				Files:         []string{"test_data/nonexistent.json"},
				UseHJSON:      false,
				StripComments: true,
			},
			wantErr: true,
			errMsg:  "invalid file path",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var savePath string
			if tc.savePath != "" {
				tempDir := t.TempDir()
				savePath = filepath.Join(tempDir, tc.savePath)
			}

			err := MergeConfigsSaveAs(savePath, tc.options)
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
				t.Fatalf("MergeConfigsSaveAs failed: %v", err)
			}

			actualData, err := os.ReadFile(savePath)
			if err != nil {
				t.Fatalf("failed to read saved file: %v", err)
			}

			expectedData, err := os.ReadFile(tc.expectedFile)
			if err != nil {
				t.Fatalf("failed to read expected file: %v", err)
			}

			var actual, expected map[string]interface{}
			if err := json.Unmarshal(actualData, &actual); err != nil {
				t.Fatalf("failed to unmarshal actual output: %v", err)
			}
			if err := json.Unmarshal(expectedData, &expected); err != nil {
				t.Fatalf("failed to unmarshal expected output: %v", err)
			}

			if !reflect.DeepEqual(actual, expected) {
				actualBytes, _ := json.MarshalIndent(actual, "", "  ")
				expectedBytes, _ := json.MarshalIndent(expected, "", "  ")
				t.Errorf("saved output does not match expected\nGot:\n%s\n\nExpected:\n%s\n",
					string(actualBytes), string(expectedBytes))
			}
		})
	}
}

func TestMergeFilesIntoMap(t *testing.T) {
	tests := []struct {
		name         string
		files        []string
		writeToFile  string
		expectedFile string
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "valid_merge",
			files:        []string{"test_data/base_deep.json", "test_data/override_deep.json", "test_data/override_deep2.json"},
			writeToFile:  "generated_map_output.json",
			expectedFile: "test_data/expected_deep2_merged.json",
		},
		{
			name:    "empty_files",
			files:   []string{},
			wantErr: true,
			errMsg:  "no config files provided",
		},
		{
			name:    "nil_files",
			files:   nil,
			wantErr: true,
			errMsg:  "no config files provided",
		},
		{
			name:    "invalid_file",
			files:   []string{"test_data/nonexistent.json"},
			wantErr: true,
			errMsg:  "invalid file path",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			writeToFile := ""
			if tc.writeToFile != "" {
				writeToFile = filepath.Join(tempDir, tc.writeToFile)
			}

			merged, err := MergeFilesIntoMap(tc.files, writeToFile)
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
				t.Fatalf("MergeFilesIntoMap failed: %v", err)
			}

			expected, err := loadJSONMapFromFile(t, tc.expectedFile)
			if err != nil {
				t.Fatalf("failed to load expected file: %v", err)
			}

			var expectedBytes []byte
			if !reflect.DeepEqual(merged, expected) {
				mergedBytes, _ := json.MarshalIndent(merged, "", "  ")
				expectedBytes, _ = json.MarshalIndent(expected, "", "  ")
				t.Errorf("merged result does not match expected\nGot:\n%s\n\nExpected:\n%s\n",
					string(mergedBytes), string(expectedBytes))
			}

			if tc.writeToFile != "" {
				actualData, err := os.ReadFile(writeToFile)
				if err != nil {
					t.Fatalf("failed to read saved file: %v", err)
				}
				var actual map[string]interface{}
				if err := json.Unmarshal(actualData, &actual); err != nil {
					t.Fatalf("failed to unmarshal saved output: %v", err)
				}
				if !reflect.DeepEqual(actual, expected) {
					actualBytes, _ := json.MarshalIndent(actual, "", "  ")
					t.Errorf("saved output does not match expected\nGot:\n%s\n\nExpected:\n%s\n",
						string(actualBytes), string(expectedBytes))
				}
			}
		})
	}
}

func TestLoadFileIntoElseMerge(t *testing.T) {
	tests := []struct {
		name         string
		loadFile     string
		files        []string
		expectedFile string
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "load_single_file",
			loadFile:     "test_data/base_deep.json",
			files:        []string{},
			expectedFile: "test_data/base_deep.json",
		},
		{
			name:         "merge_files",
			loadFile:     "test_data/nonexistent.json",
			files:        []string{"test_data/base_deep.json", "test_data/override_deep.json"},
			expectedFile: "test_data/expected_deep_merged.json",
		},
		{
			name:     "nil_target",
			loadFile: "test_data/base.json",
			wantErr:  true,
			errMsg:   "target cannot be nil",
		},
		{
			name:     "non_pointer_target",
			loadFile: "test_data/base.json",
			wantErr:  true,
			errMsg:   "target must be a pointer",
		},
		{
			name:     "no_files",
			loadFile: "",
			files:    []string{},
			wantErr:  true,
			errMsg:   "no files specified",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var cfg WebConfig
			var target interface{}
			if tc.name == "non_pointer_target" {
				target = WebConfig{}
			} else if tc.name != "nil_target" {
				target = &cfg
			}

			err := LoadFileIntoElseMerge(tc.loadFile, tc.files, target)
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
				t.Fatalf("LoadFileIntoElseMerge failed: %v", err)
			}

			expected, err := loadWebConfigFromFile(t, tc.expectedFile)
			if err != nil {
				t.Fatalf("failed to load expected file: %v", err)
			}

			if !reflect.DeepEqual(cfg, *expected) {
				cfgBytes, _ := json.MarshalIndent(cfg, "", "  ")
				expectedBytes, _ := json.MarshalIndent(expected, "", "  ")
				t.Errorf("merged config does not match expected\nGot:\n%s\n\nExpected:\n%s\n",
					string(cfgBytes), string(expectedBytes))
			}
		})
	}
}

func TestStripComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single_line_comments",
			input:    `{"key": "value"} // comment`,
			expected: `{"key": "value"} `,
		},
		{
			name:     "multi_line_comments",
			input:    `{"key": /* comment */ "value"}`,
			expected: `{"key":  "value"}`,
		},
		{
			name:     "mixed_comments",
			input:    `{"key": "value" /* multi */} // single`,
			expected: `{"key": "value" } `,
		},
		{
			name:     "no_comments",
			input:    `{"key": "value"}`,
			expected: `{"key": "value"}`,
		},
		{
			name: "json_with_comments",
			input: `{
  "app": "myapp",
  "version": "1.0" // version comment
}`,
			expected: `{
  "app": "myapp",
  "version": "1.0" 
}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := string(StripComments([]byte(tc.input)))
			if result != tc.expected {
				t.Errorf("StripComments failed\nGot: %q\nExpected: %q", result, tc.expected)
			}
		})
	}
}
