package ajson

import (
	"testing"
)

// TestIsEmpty checks if the IsEmpty method correctly identifies empty JsonKeys.
func TestIsEmpty(t *testing.T) {
	tests := []struct {
		key      JsonKey
		expected bool
	}{
		{"", true},
		{" ", true},
		{"key", false},
	}

	for _, test := range tests {
		if test.key.IsEmpty() != test.expected {
			t.Errorf("IsEmpty() for key '%s' - expected %v, got %v", test.key, test.expected, !test.expected)
		}
	}
}

// TestTrimSpace checks if the TrimSpace method correctly trims whitespace from JsonKeys.
func TestTrimSpace(t *testing.T) {
	tests := []struct {
		key      JsonKey
		expected JsonKey
	}{
		{" key ", "key"},
		{"  key", "key"},
		{"key  ", "key"},
	}

	for _, test := range tests {
		if test.key.TrimSpace() != test.expected {
			t.Errorf("TrimSpace() for key '%s' - expected '%s', got '%s'", test.key, test.expected, test.key.TrimSpace())
		}
	}
}

// TestIsRoot checks if the IsRoot method correctly identifies root JsonKeys.
func TestIsRoot(t *testing.T) {
	tests := []struct {
		key      JsonKey
		expected bool
	}{
		{"key", true},
		{"key.subkey", false},
	}

	for _, test := range tests {
		if test.key.IsRoot() != test.expected {
			t.Errorf("IsRoot() for key '%s' - expected %v, got %v", test.key, test.expected, !test.expected)
		}
	}
}

// TestGetRoot checks if the GetRoot method correctly extracts the root from JsonKeys.
func TestGetRoot(t *testing.T) {
	tests := []struct {
		key      JsonKey
		expected JsonKey
	}{
		{"key.subkey", "key"},
		{"key", "key"},
	}

	for _, test := range tests {
		if test.key.GetRoot() != test.expected {
			t.Errorf("GetRoot() for key '%s' - expected '%s', got '%s'", test.key, test.expected, test.key.GetRoot())
		}
	}
}

// TestGetPathLeaf checks if the GetPathLeaf method correctly extracts the leaf from JsonKeys.
func TestGetPathLeaf(t *testing.T) {
	tests := []struct {
		key      JsonKey
		expected JsonKey
	}{
		{"key.subkey.leaf", "leaf"},
		{"key", "key"},
	}

	for _, test := range tests {
		if test.key.GetPathLeaf() != test.expected {
			t.Errorf("GetPathLeaf() for key '%s' - expected '%s', got '%s'", test.key, test.expected, test.key.GetPathLeaf())
		}
	}
}

// TestGetPathParts checks if the GetPathParts method correctly splits the JsonKeys into parts.
func TestGetPathParts(t *testing.T) {
	key := JsonKey("key.subkey.leaf")
	expected := []string{"key", "subkey", "leaf"}

	parts := key.GetPathParts()
	for i, part := range parts {
		if part != expected[i] {
			t.Errorf("GetPathParts() - expected '%s', got '%s'", expected[i], part)
		}
	}
}

// TestGetPathParent checks if the GetPathParent method correctly extracts the parent path from JsonKeys.
func TestGetPathParent(t *testing.T) {
	tests := []struct {
		key      JsonKey
		expected JsonKey
	}{
		{"key.subkey.leaf", "key.subkey"},
		{"key", ""},
	}

	for _, test := range tests {
		if test.key.GetPathParent() != test.expected {
			t.Errorf("GetPathParent() for key '%s' - expected '%s', got '%s'", test.key, test.expected, test.key.GetPathParent())
		}
	}
}

// TestAdd checks if the Add method correctly appends a target JsonKey to the current JsonKey path.
func TestAdd(t *testing.T) {
	key := JsonKey("key")
	target := JsonKey("subkey")
	expected := JsonKey("key.subkey")

	result := key.Add(target)
	if result != expected {
		t.Errorf("Add() - expected '%s', got '%s'", expected, result)
	}
}

// TestCopyPlusAdd checks if the CopyPlusAdd method correctly creates a new JsonKey by appending a target JsonKey.
func TestCopyPlusAdd(t *testing.T) {
	key := JsonKey("key")
	target := JsonKey("subkey")
	expected := JsonKey("key.subkey")

	result := key.CopyPlusAdd(target)
	if result != expected {
		t.Errorf("CopyPlusAdd() - expected '%s', got '%s'", expected, result)
	}
}

// TestCopyPlusAddInt checks if the CopyPlusAddInt method correctly creates a new JsonKey by appending an integer.
func TestCopyPlusAddInt(t *testing.T) {
	key := JsonKey("key")
	target := 1
	expected := JsonKey("key.1")

	result := key.CopyPlusAddInt(target)
	if result != expected {
		t.Errorf("CopyPlusAddInt() - expected '%s', got '%s'", expected, result)
	}
}

func TestJsonKey(t *testing.T) {
	key := JsonKey("root")
	if key.IsRoot() != true {
		t.Errorf("Expected IsRoot to be true, got false")
	}
	if key.String() != "root" {
		t.Errorf("Expected String to be 'root', got '%s'", key.String())
	}
	if key.TrimSpace() != JsonKey("root") {
		t.Errorf("Expected TrimSpace to be 'root', got '%s'", key.TrimSpace())
	}
	if key.IsEmpty() != false {
		t.Errorf("Expected IsEmpty to be false, got true")
	}
	if key.GetPathLeaf() != JsonKey("root") {
		t.Errorf("Expected GetPathLeaf to be 'root', got '%s'", key.GetPathLeaf())
	}
	if key.GetPathParent() != JsonKey("") {
		t.Errorf("Expected GetPathParent to be '', got '%s'", key.GetPathParent())
	}
	if len(key.GetPathParts()) != 1 {
		t.Errorf("Expected GetPathParts length to be 1, got %d", len(key.GetPathParts()))
	}

	key = JsonKey("root.obj1")
	if key.IsRoot() != false {
		t.Errorf("Expected IsRoot to be false, got true")
	}
	if key.String() != "root.obj1" {
		t.Errorf("Expected String to be 'root.obj1', got '%s'", key.String())
	}
	if key.GetPathLeaf() != JsonKey("obj1") {
		t.Errorf("Expected GetPathLeaf to be 'obj1', got '%s'", key.GetPathLeaf())
	}
	if key.GetPathParent() != JsonKey("root") {
		t.Errorf("Expected GetPathParent to be 'root', got '%s'", key.GetPathParent())
	}
	if len(key.GetPathParts()) != 2 {
		t.Errorf("Expected GetPathParts length to be 2, got %d", len(key.GetPathParts()))
	}

	key = key.Add("arr1")
	if key.String() != "root.obj1.arr1" {
		t.Errorf("Expected String after Add to be 'root.obj1.arr1', got '%s'", key.String())
	}
	if key.GetPathLeaf() != JsonKey("arr1") {
		t.Errorf("Expected GetPathLeaf to be 'arr1', got '%s'", key.GetPathLeaf())
	}
	if key.GetPathParent() != JsonKey("root.obj1") {
		t.Errorf("Expected GetPathParent to be 'root.obj1', got '%s'", key.GetPathParent())
	}
	if len(key.GetPathParts()) != 3 {
		t.Errorf("Expected GetPathParts length to be 3, got %d", len(key.GetPathParts()))
	}

	jkNew := key.CopyPlusAdd("fld1")
	if jkNew.String() != "root.obj1.arr1.fld1" {
		t.Errorf("Expected String after CopyPlusAdd to be 'root.obj1.arr1.fld1', got '%s'", jkNew.String())
	}
	if jkNew.GetPathLeaf() != JsonKey("fld1") {
		t.Errorf("Expected GetPathLeaf to be 'fld1', got '%s'", jkNew.GetPathLeaf())
	}
	if jkNew.GetPathParent() != JsonKey("root.obj1.arr1") {
		t.Errorf("Expected GetPathParent to be 'root.obj1.arr1', got '%s'", jkNew.GetPathParent())
	}
	if len(jkNew.GetPathParts()) != 4 {
		t.Errorf("Expected GetPathParts length to be 4, got %d", len(jkNew.GetPathParts()))
	}
}
