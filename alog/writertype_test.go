package alog

import (
	"testing"
)

// TestIsEmpty checks the IsEmpty method for various WriterType values.
func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		writer   WriterType
		expected bool
	}{
		{"Empty", "", true},
		{"Whitespace", "   ", true},
		{"NotEmpty", "file", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.writer.IsEmpty(); got != test.expected {
				t.Errorf("WriterType.IsEmpty() = %v, want %v", got, test.expected)
			}
		})
	}
}

// TestTrimSpace checks the TrimSpace method for various WriterType values.
func TestTrimSpace(t *testing.T) {
	tests := []struct {
		name     string
		writer   WriterType
		expected WriterType
	}{
		{"LeadingSpace", "  file", "file"},
		{"TrailingSpace", "file  ", "file"},
		{"Both", "  file  ", "file"},
		{"NoSpace", "file", "file"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.writer.TrimSpace(); got != test.expected {
				t.Errorf("WriterType.TrimSpace() = %v, want %v", got, test.expected)
			}
		})
	}
}

// TestString checks the String method for various WriterType values.
func TestString(t *testing.T) {
	tests := []struct {
		name     string
		writer   WriterType
		expected string
	}{
		{"Simple", "file", "file"},
		{"Complex", "console-stdout", "console-stdout"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.writer.String(); got != test.expected {
				t.Errorf("WriterType.String() = %v, want %v", got, test.expected)
			}
		})
	}
}

// TestHasMatch checks the HasMatch method for various WriterType values.
func TestHasMatch(t *testing.T) {
	tests := []struct {
		name     string
		writer   WriterType
		arg      WriterType
		expected bool
	}{
		{"Match", "file", "file", true},
		{"NoMatch", "file", "console", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.writer.HasMatch(test.arg); got != test.expected {
				t.Errorf("WriterType.HasMatch() = %v, want %v", got, test.expected)
			}
		})
	}
}

// TestMatchesOne checks the MatchesOne method for various WriterType values.
func TestMatchesOne(t *testing.T) {
	tests := []struct {
		name     string
		writer   WriterType
		args     []WriterType
		expected bool
	}{
		{"OneMatch", "file", []WriterType{"console", "file"}, true},
		{"NoMatches", "file", []WriterType{"console", "stdout"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.writer.MatchesOne(test.args...); got != test.expected {
				t.Errorf("WriterType.MatchesOne() = %v, want %v", got, test.expected)
			}
		})
	}
}
