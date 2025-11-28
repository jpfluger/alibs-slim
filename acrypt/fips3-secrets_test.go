package acrypt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"runtime"
	"testing"
)

// mockReader is a custom reader that always returns an error.
type mockReader struct{}

func (m *mockReader) Read(b []byte) (n int, err error) {
	return 0, fmt.Errorf("mock error")
}

// TestGenerateSecretKeySuccess tests the successful generation of a secret key.
func TestGenerateSecretKeySuccess(t *testing.T) {
	key, err := GenerateSecretKey()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(key) != 32 {
		t.Fatalf("expected key length 32, got %d", len(key))
	}
}

// TestGenerateSecretKeyFailure tests the failure scenario of GenerateSecretKey.
func TestGenerateSecretKeyFailure(t *testing.T) {
	// Get Go version
	goVersion := runtime.Version()
	t.Logf("Go version: %s", goVersion) // Debug log

	// Parse major and minor version for robust check
	var major, minor int
	n, err := fmt.Sscanf(goVersion, "go%d.%d", &major, &minor)
	isGo124OrLater := err == nil && n == 2 && (major > 1 || (major == 1 && minor >= 24))
	t.Logf("isGo124OrLater: %t", isGo124OrLater) // Debug log

	if isGo124OrLater {
		t.Skip("Skipping failure test in Go 1.24+: crypto/rand.Read failures cause irrecoverable crashes, cannot be mocked or recovered.")
	}

	// Replace the global rand.Reader with our mockReader
	originalReader := rand.Reader
	rand.Reader = &mockReader{}
	defer func() { rand.Reader = originalReader }()

	// In Go 1.23 and earlier, expect returned error
	_, err = GenerateSecretKey()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	expected := "failed to generate secret key: mock error"
	if err.Error() != expected {
		t.Fatalf("expected error message '%s', got '%v'", expected, err)
	}
}

// TestRandomStrongOneWaySecret tests the default wrapper function.
func TestRandomStrongOneWaySecret(t *testing.T) {
	pass, err := RandomStrongOneWaySecret()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	length := len(pass)
	if length < 25 || length > 37 {
		t.Errorf("Expected length between 25 and 37, got %d", length)
	}
	if !isBase64Like(pass) {
		t.Errorf("Generated pass does not look like base64: %s", pass)
	}
}

// TestRandomStrongOneWayByVariableLength_CustomRange tests a custom range.
func TestRandomStrongOneWayByVariableLength_CustomRange(t *testing.T) {
	low, high := 10, 20
	pass, err := RandomStrongOneWayByVariableLength(low, high)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	length := len(pass)
	if length < low || length > high {
		t.Errorf("Expected length between %d and %d, got %d", low, high, length)
	}
	if !isBase64Like(pass) {
		t.Errorf("Generated pass does not look like base64: %s", pass)
	}
}

// TestRandomStrongOneWayByVariableLength_InvalidRange tests fallback to defaults.
func TestRandomStrongOneWayByVariableLength_InvalidRange(t *testing.T) {
	pass, err := RandomStrongOneWayByVariableLength(40, 30) // Invalid (low > high)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	length := len(pass)
	if length < 25 || length > 37 {
		t.Errorf("Expected default length between 25 and 37, got %d", length)
	}
	if !isBase64Like(pass) {
		t.Errorf("Generated pass does not look like base64: %s", pass)
	}
}

// TestRandomStrongOneWayByVariableLength_Failure tests rand failure (pre-Go 1.24).
func TestRandomStrongOneWayByVariableLength_Failure(t *testing.T) {
	// Skip in Go 1.24+ where rand.Read panics on failure.
	goVersion := runtime.Version()
	var major, minor int
	n, _ := fmt.Sscanf(goVersion, "go%d.%d", &major, &minor)
	isGo124OrLater := n == 2 && (major > 1 || (major == 1 && minor >= 24))
	if isGo124OrLater {
		t.Skip("Skipping failure test in Go 1.24+: crypto/rand.Read panics on failure")
	}

	// Mock rand.Reader.
	originalReader := rand.Reader
	rand.Reader = &mockReader{}
	defer func() { rand.Reader = originalReader }()

	_, err := RandomStrongOneWayByVariableLength(25, 37)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	expected := "failed to generate secret key: mock error"
	if err.Error() != expected {
		t.Errorf("Expected error '%s', got '%v'", expected, err)
	}
}

// isBase64Like is a helper to validate if string resembles base64 (basic check).
// isBase64Like is a helper to validate if string resembles base64 (basic check).
func isBase64Like(s string) bool {
	l := len(s)
	mod := l % 4
	padded := s
	switch mod {
	case 0:
		// No padding needed.
	case 1:
		return false // Invalid for Base64.
	case 2:
		padded += "=="
	case 3:
		padded += "="
	}
	_, err := base64.StdEncoding.DecodeString(padded)
	return err == nil
}
