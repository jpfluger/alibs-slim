package ascripts

import (
	"github.com/jpfluger/alibs-slim/ajson"
	"testing"
)

// TestFind checks if the Find method correctly retrieves a *Scriptable
func TestFind(t *testing.T) {
	// Create a mock Scriptable and JsonKey
	mockKey := ajson.JsonKey("testKey")
	mockValue := &Scriptable{Key: mockKey}

	// Create a ScriptablesMap with the mock Scriptable
	sMap := ScriptablesMap{mockKey: mockValue}

	// Test finding an existing key
	if sMap.Find(mockKey) != mockValue {
		t.Errorf("Find method failed to retrieve the existing key")
	}

	// Test finding a non-existing key
	if sMap.Find(ajson.JsonKey("nonExistingKey")) != nil {
		t.Errorf("Find method incorrectly retrieved a non-existing key")
	}
}

// TestHasItem checks if the HasItem method correctly identifies existence of a *Scriptable
func TestHasItem(t *testing.T) {
	// Create a mock Scriptable and JsonKey
	mockKey := ajson.JsonKey("testKey")
	mockValue := &Scriptable{Key: mockKey}

	// Create a ScriptablesMap with the mock Scriptable
	sMap := ScriptablesMap{mockKey: mockValue}

	// Test for an existing key
	if !sMap.HasItem(mockKey) {
		t.Errorf("HasItem method failed to identify the existing key")
	}

	// Test for a non-existing key
	if sMap.HasItem(ajson.JsonKey("nonExistingKey")) {
		t.Errorf("HasItem method incorrectly identified a non-existing key")
	}
}
