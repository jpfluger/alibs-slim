package autils

import (
	"testing"
)

// TestIsEmpty checks the IsEmpty method for AppPathKey.
func TestIsEmpty(t *testing.T) {
	tests := []struct {
		key      AppPathKey
		expected bool
	}{
		{"", true},
		{" ", true},
		{"not-empty", false},
	}

	for _, test := range tests {
		if test.key.IsEmpty() != test.expected {
			t.Errorf("IsEmpty() = %v, want %v for key %v", test.key.IsEmpty(), test.expected, test.key)
		}
	}
}

// TestTrimSpace checks the TrimSpace method for AppPathKey.
func TestTrimSpace(t *testing.T) {
	tests := []struct {
		key      AppPathKey
		expected AppPathKey
	}{
		{"  padded  ", "padded"},
		{"\tnotrim\n", "notrim"},
		{" no-change ", "no-change"},
	}

	for _, test := range tests {
		if trimmed := test.key.TrimSpace(); trimmed != test.expected {
			t.Errorf("TrimSpace() = %v, want %v for key %v", trimmed, test.expected, test.key)
		}
	}
}

// TestToStringTrimLower checks the ToStringTrimLower method for AppPathKey.
func TestToStringTrimLower(t *testing.T) {
	key := AppPathKey("  MIXEDcase  ")
	expected := "mixedcase"
	if result := key.ToStringTrimLower(); result != expected {
		t.Errorf("ToStringTrimLower() = %v, want %v for key %v", result, expected, key)
	}
}

// TestType checks the Type method for AppPathKey.
func TestType(t *testing.T) {
	tests := []struct {
		key      AppPathKey
		expected string
	}{
		{"DIR_ROOT", "dir"},
		{"FILE_CONFIG", "file"},
		{"UNKNOWN_TYPE", ""},
	}

	for _, test := range tests {
		if typ := test.key.Type(); typ != test.expected {
			t.Errorf("Type() = %v, want %v for key %v", typ, test.expected, test.key)
		}
	}
}

// TestIsDir checks the IsDir method for AppPathKey.
func TestIsDir(t *testing.T) {
	if !DIR_ROOT.IsDir() {
		t.Errorf("IsDir() = false, want true for key %v", DIR_ROOT)
	}
	FILE_TEST := AppPathKey("FILE_TEST")
	if FILE_TEST.IsDir() {
		t.Errorf("IsDir() = true, want false for key %v", FILE_TEST)
	}
}

// TestIsFile checks the IsFile method for AppPathKey.
func TestIsFile(t *testing.T) {
	fileKey := AppPathKey("FILE_CONFIG")
	if !fileKey.IsFile() {
		t.Errorf("IsFile() = false, want true for key %v", fileKey)
	}
	if DIR_ETC.IsFile() {
		t.Errorf("IsFile() = true, want false for key %v", DIR_ETC)
	}
}

// Add more tests as needed to cover edge cases and other scenarios.
