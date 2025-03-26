package acrypt

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateIDBase36Caps(t *testing.T) {
	length := 10
	id := GenerateIDBase36Caps(length)

	// Ensure the ID has the correct length
	if len(id) != length {
		t.Errorf("Expected length %d, got %d", length, len(id))
	}

	// Ensure the ID contains only allowed characters
	const allowedChars = "123456789ABCDEFGHIJKLMNPQRSTUVWXYZ"
	for _, char := range id {
		if !strings.ContainsRune(allowedChars, char) {
			t.Errorf("Generated ID contains invalid character: %c", char)
		}
	}

	// Ensure "0" and "O" are not included
	if strings.ContainsRune(id, '0') || strings.ContainsRune(id, 'O') {
		t.Errorf("Generated ID contains disallowed characters '0' or 'O'")
	}
}

func TestNewIdGenReadableWithOptions(t *testing.T) {
	format := "%s-%s-%s"
	prefix := "INV"
	date := time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC)
	length := 8

	id := NewIdGenReadableWithOptions(format, prefix, date, length)

	// Check format
	expectedPrefix := "INV-20250109-"
	if !strings.HasPrefix(id, expectedPrefix) {
		t.Errorf("Expected ID to start with %q, got %q", expectedPrefix, id)
	}

	// Check random part length
	randomPart := strings.TrimPrefix(id, expectedPrefix)
	if len(randomPart) != length {
		t.Errorf("Expected random part length %d, got %d", length, len(randomPart))
	}
}

func TestNewIdGenReadableShort(t *testing.T) {
	prefix := "PO"
	date := time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC)

	id := NewIdGenReadableShort(prefix, date)

	// Check format
	expectedPrefix := "PO20250109-"
	if !strings.HasPrefix(id, expectedPrefix) {
		t.Errorf("Expected ID to start with %q, got %q", expectedPrefix, id)
	}

	// Check random part length
	randomPart := strings.TrimPrefix(id, expectedPrefix)
	if len(randomPart) != 7 {
		t.Errorf("Expected random part length 7, got %d", len(randomPart))
	}
}

func TestNewIdGenReadableLong(t *testing.T) {
	prefix := "QUOTE"
	date := time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC)

	id := NewIdGenReadableLong(prefix, date)

	// Check format
	expectedPrefix := "QUOTE-20250109-"
	if !strings.HasPrefix(id, expectedPrefix) {
		t.Errorf("Expected ID to start with %q, got %q", expectedPrefix, id)
	}

	// Check random part length
	randomPart := strings.TrimPrefix(id, expectedPrefix)
	if len(randomPart) != 7 {
		t.Errorf("Expected random part length 7, got %d", len(randomPart))
	}
}

//func TestNewIdGenReadableLongUniqueness(t *testing.T) {
//	const iterations = 100000
//	seen := make(map[string]bool)
//	date := time.Now()
//	seenCounter := 0
//	for i := 0; i < iterations; i++ {
//		id := NewIdGenReadableLong("TEST", date)
//		if seen[id] {
//			seenCounter++
//			//t.Errorf("Collision detected for ID: %s", id)
//		}
//		seen[id] = true
//	}
//	assert.Equal(t, seenCounter, 0, "Should have 0 collisions")
//}
