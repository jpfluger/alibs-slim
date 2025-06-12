package ajson

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

// Helper to unmarshal JSON file into map
func loadJSONMapFromFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func TestMergeJSONFiles(t *testing.T) {
	tests := []struct {
		name         string
		options      MergeOptions
		expectedFile string
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
				Files: []string{
					"test_data/base_deep.json",
					"test_data/override_deep.json",
					"test_data/override_deep2.json",
				},
				UseHJSON:      false,
				StripComments: false,
			},
			expectedFile: "test_data/expected_deep2_merged.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			merged, err := MergeConfigs(tc.options)
			if err != nil {
				t.Fatalf("merge failed: %v", err)
			}

			expected, err := loadJSONMapFromFile(tc.expectedFile)
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

// Typed config struct for testing MergeConfigsInto
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

func TestMergeConfigsInto_TypedStruct(t *testing.T) {
	opts := MergeOptions{
		Files: []string{
			"test_data/base_deep.json",
			"test_data/override_deep.json",
			"test_data/override_deep2.json",
		},
		UseHJSON:      false,
		StripComments: false,
	}

	var cfg WebConfig
	err := MergeConfigsInto(&cfg, opts)
	if err != nil {
		t.Fatalf("MergeConfigsInto failed: %v", err)
	}

	expected := WebConfig{
		Services: struct {
			Web struct {
				Port     int             `json:"port"`
				Timeout  int             `json:"timeout"`
				Tags     []string        `json:"tags"`
				Features map[string]bool `json:"features"`
			} `json:"web"`
		}{
			Web: struct {
				Port     int             `json:"port"`
				Timeout  int             `json:"timeout"`
				Tags     []string        `json:"tags"`
				Features map[string]bool `json:"features"`
			}{
				Port:    3000,
				Timeout: 120,
				Tags:    []string{"v2", "stable"},
				Features: map[string]bool{
					"caching": false,
					"gzip":    true,
				},
			},
		},
		Global: struct {
			Logging struct {
				Level   string `json:"level"`
				Enabled bool   `json:"enabled"`
			} `json:"logging"`
			Theme string `json:"theme"`
		}{
			Logging: struct {
				Level   string `json:"level"`
				Enabled bool   `json:"enabled"`
			}{
				Level:   "debug",
				Enabled: true,
			},
			Theme: "dark",
		},
		Metadata: struct {
			Version string `json:"version"`
			Build   int    `json:"build"`
		}{
			Version: "1.0",
			Build:   1002,
		},
	}

	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("Merged config did not match expected.\nGot: %+v\nExpected: %+v", cfg, expected)
	}
}

func TestMergeConfigsSaveAs(t *testing.T) {
	savePath := "test_data/generated_save_output.json"

	opts := MergeOptions{
		Files: []string{
			"test_data/base_deep.json",
			"test_data/override_deep.json",
			"test_data/override_deep2.json",
		},
		UseHJSON:      false,
		StripComments: false,
	}

	// Run the save function
	err := MergeConfigsSaveAs(savePath, opts)
	if err != nil {
		t.Fatalf("MergeConfigsSaveAs failed: %v", err)
	}
	defer os.Remove(savePath) // clean up

	// Load expected and actual output
	actualData, err := os.ReadFile(savePath)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	expectedData, err := os.ReadFile("test_data/expected_deep2_merged.json")
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
		t.Errorf("saved output mismatch\nGot:      %+v\nExpected: %+v", actual, expected)
	}
}
