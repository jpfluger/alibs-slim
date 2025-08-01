package acrypt

import (
	"crypto/rand"
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
