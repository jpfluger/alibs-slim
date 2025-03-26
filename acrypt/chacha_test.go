package acrypt

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/chacha20"
	"testing"
)

// mockReader is a custom reader that always returns an error.
type mockReader struct{}

func (m *mockReader) Read(b []byte) (n int, err error) {
	return 0, fmt.Errorf("mock error")
}

// TestGenerateChaCha20KeySuccess tests the successful generation of a ChaCha20 key.
func TestGenerateChaCha20KeySuccess(t *testing.T) {
	key, err := GenerateChaCha20Key()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(key) != chacha20.KeySize {
		t.Fatalf("expected key length %d, got %d", chacha20.KeySize, len(key))
	}
}

// TestGenerateChaCha20KeyFailure tests the failure scenario of GenerateChaCha20Key function.
func TestGenerateChaCha20KeyFailure(t *testing.T) {
	// Replace the global rand.Reader with our mockReader
	originalReader := rand.Reader
	rand.Reader = &mockReader{}
	defer func() { rand.Reader = originalReader }()

	_, err := GenerateChaCha20Key()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "failed to generate ChaCha20 key: mock error" {
		t.Fatalf("expected error message 'failed to generate ChaCha20 key: mock error', got %v", err)
	}
}
