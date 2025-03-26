package atemplates

import (
	"testing"
)

// TestFormatIntegerComma tests the FormatIntegerComma function.
func TestFormatIntegerComma(t *testing.T) {
	tests := []struct {
		name string
		val  int
		want string
	}{
		{"Zero", 0, "0"},
		{"Thousand", 1000, "1,000"},
		{"Million", 1000000, "1,000,000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatIntegerComma(tt.val); got != tt.want {
				t.Errorf("FormatIntegerComma() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFormatBytes tests the FormatBytes function.
func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name string
		val  int
		want string
	}{
		{"Zero", 0, "0 B"},
		{"Kilobyte", 1024, "1.0 kB"},
		{"Megabyte", 1048576, "1.0 MB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatBytes(tt.val); got != tt.want {
				t.Errorf("FormatBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestToInt tests the ToInt function.
func TestToInt(t *testing.T) {
	tests := []struct {
		name string
		val  interface{}
		want int
	}{
		{"Int", 42, 42},
		{"String", "42", 42},
		{"InvalidString", "abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt(tt.val); got != tt.want {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Additional tests for ToInt64, ToUInt64, ToFloat64, ToBool, and ToString can be added similarly.

// TestGetTextTemplateFunctions tests the GetTextTemplateFunctions function.
func TestGetTextTemplateFunctions(t *testing.T) {
	// Assuming TEMPLATE_FUNCTIONS_COMMON is a valid constant in your package
	funcMap := GetTextTemplateFunctions(TEMPLATE_FUNCTIONS_COMMON)
	if funcMap == nil || len(*funcMap) == 0 {
		t.Errorf("GetTextTemplateFunctions() returned an empty or nil FuncMap")
	}
}

// Additional tests can be written to test the individual functions within the FuncMap.
