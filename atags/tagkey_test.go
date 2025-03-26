package atags

import (
	"testing"
)

// TestNewTagKey tests the NewTagKey function for various inputs.
func TestNewTagKey(t *testing.T) {
	tests := []struct {
		kType, kId, expected TagKey
	}{
		{"type", "id", "type:id"},
		{"type", "", "type"},
		{"", "id", "id"},
		{"", "", ""},
	}

	for _, test := range tests {
		result := NewTagKey(test.kType.String(), test.kId.String())
		if result != test.expected {
			t.Errorf("NewTagKey(%q, %q) = %q; want %q", test.kType, test.kId, result, test.expected)
		}
	}
}

// TestGetType tests the GetType method.
func TestGetType(t *testing.T) {
	key := NewTagKey("type", "id")
	if key.GetType() != "type" {
		t.Errorf("Expected type to be 'type', got %s", key.GetType())
	}
}

// TestGetUniqueId tests the GetUniqueId method.
func TestGetUniqueId(t *testing.T) {
	key := NewTagKey("type", "id")
	if key.GetUniqueId() != "id" {
		t.Errorf("Expected unique ID to be 'id', got %s", key.GetUniqueId())
	}
}

// TestIsEmpty tests the IsEmpty method.
func TestIsEmpty(t *testing.T) {
	if !NewTagKey("", "").IsEmpty() {
		t.Error("Expected key to be empty")
	}
	if NewTagKey("type", "id").IsEmpty() {
		t.Error("Expected key not to be empty")
	}
}

// TestTrimSpace tests the TrimSpace method.
func TestTrimSpace(t *testing.T) {
	key := TagKey(" type : id ")
	if key.TrimSpace() != "type:id" {
		t.Errorf("Expected trimmed key to be 'type:id', got %q", key.TrimSpace())
	}
}

// TestString tests the String method.
func TestString(t *testing.T) {
	key := NewTagKey("type", "id")
	if key.String() != "type:id" {
		t.Errorf("Expected String() to return 'type:id', got %q", key.String())
	}
}
