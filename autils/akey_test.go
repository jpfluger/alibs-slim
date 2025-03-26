package autils

import (
	"fmt"
	"testing"
)

// TestAKeyIsEmpty checks the IsEmpty method.
func TestAKeyIsEmpty(t *testing.T) {
	tests := []struct {
		key      AKey
		expected bool
	}{
		{"", true},
		{" ", true},
		{"not-empty", false},
	}

	for _, test := range tests {
		if test.key.IsEmpty() != test.expected {
			t.Errorf("Expected '%v' for key '%v', got '%v'", test.expected, test.key, !test.expected)
		}
	}
}

// TestAKeyTrimSpace checks the TrimSpace method.
func TestAKeyTrimSpace(t *testing.T) {
	tests := []struct {
		key      AKey
		expected AKey
	}{
		{"  padded  ", "padded"},
		{"\tnotrim\n", "notrim"},
		{" no-change ", "no-change"},
	}

	for _, test := range tests {
		if trimmed := test.key.TrimSpace(); trimmed != test.expected {
			t.Errorf("Expected '%v' for key '%v', got '%v'", test.expected, test.key, trimmed)
		}
	}
}

// TestAKeyHasMatch checks the HasMatch method.
func TestAKeyHasMatch(t *testing.T) {
	key := AKey("match")
	if !key.HasMatch("match") {
		t.Errorf("Expected true for matching keys, got false")
	}
}

// TestAKeyMatchesOne checks the MatchesOne method.
func TestAKeyMatchesOne(t *testing.T) {
	key := AKey("one")
	if !key.MatchesOne("zero", "one", "two") {
		t.Errorf("Expected true for one matching key, got false")
	}
}

// TestAKeyHasPrefix checks the HasPrefix method.
func TestAKeyHasPrefix(t *testing.T) {
	key := AKey("prefix-key")
	if !key.HasPrefix("prefix") {
		t.Errorf("Expected true for key with prefix, got false")
	}
}

// TestAKeyHasSuffix checks the HasSuffix method.
func TestAKeyHasSuffix(t *testing.T) {
	key := AKey("key-suffix")
	if !key.HasSuffix("suffix") {
		t.Errorf("Expected true for key with suffix, got false")
	}
}

// TestAKeyValidate checks the Validate method.
func TestAKeyValidate(t *testing.T) {
	tests := []struct {
		key      AKey
		expected error
	}{
		{"valid-key_1.2", nil},
		{"-invalid", fmt.Errorf("invalid char '-' at position 0; only alphas and numbers can begin or end the AKey")},
		{"invalid.", fmt.Errorf("invalid char '.' at position 7; only alphas and numbers can begin or end the AKey")},
		{"in*valid", fmt.Errorf("invalid char '*'; AKey allows alpha A-Z, a-z, numbers 0-9, special '-_.'")},
	}

	for _, test := range tests {
		if err := test.key.Validate(); err != nil && err.Error() != test.expected.Error() {
			t.Errorf("Expected '%v' for key '%v', got '%v'", test.expected, test.key, err)
		}
	}
}

// TestAKeysHasValues checks the HasValues method of AKeys.
func TestAKeysHasValues(t *testing.T) {
	keys := AKeys{"one", "two"}
	if !keys.HasValues() {
		t.Errorf("Expected true for AKeys with values, got false")
	}
}

// TestAKeysHasMatch checks the HasMatch method of AKeys.
func TestAKeysHasMatch(t *testing.T) {
	keys := AKeys{"one", "two"}
	if !keys.HasMatch("two") {
		t.Errorf("Expected true for AKeys with matching key, got false")
	}
}

// TestAKeysHasPrefix checks the HasPrefix method of AKeys.
func TestAKeysHasPrefix(t *testing.T) {
	keys := AKeys{"start-middle-end", "start-another"}
	if !keys.HasPrefix("start") {
		t.Errorf("Expected true for AKeys with prefix 'start', got false")
	}
}

// TestAKeysClone checks the Clone method of AKeys.
func TestAKeysClone(t *testing.T) {
	original := AKeys{"one", "two", "three"}
	cloned := original.Clone()
	if len(cloned) != len(original) {
		t.Errorf("Cloned AKeys length mismatch, expected %d, got %d", len(original), len(cloned))
	}
	for i, key := range cloned {
		if key != original[i] {
			t.Errorf("Cloned AKeys content mismatch, expected %v, got %v", original[i], key)
		}
	}
}

// TestAKeysToArrStrings checks the ToArrStrings method of AKeys.
func TestAKeysToArrStrings(t *testing.T) {
	keys := AKeys{"one", "two", "three"}
	expected := []string{"one", "two", "three"}
	result := keys.ToArrStrings()
	if len(result) != len(expected) {
		t.Errorf("ToArrStrings length mismatch, expected %d, got %d", len(expected), len(result))
	}
	for i, str := range result {
		if str != expected[i] {
			t.Errorf("ToArrStrings content mismatch, expected %s, got %s", expected[i], str)
		}
	}
}

// TestAKeysIncludeIfInTargets checks the IncludeIfInTargets method of AKeys.
func TestAKeysIncludeIfInTargets(t *testing.T) {
	keys := AKeys{"one", "two", "three"}
	targets := AKeys{"two", "four"}
	expected := AKeys{"two"}
	result := keys.IncludeIfInTargets(targets)
	if len(result) != len(expected) {
		t.Errorf("IncludeIfInTargets length mismatch, expected %d, got %d", len(expected), len(result))
	}
	for i, key := range result {
		if key != expected[i] {
			t.Errorf("IncludeIfInTargets content mismatch, expected %v, got %v", expected[i], key)
		}
	}
}

// TestAKeysClean checks the Clean method of AKeys.
func TestAKeysClean(t *testing.T) {
	keys := AKeys{"", "one", "", "two", "three", ""}
	expected := AKeys{"one", "two", "three"}
	result := keys.Clean()
	if len(result) != len(expected) {
		t.Errorf("Clean length mismatch, expected %d, got %d", len(expected), len(result))
	}
	for i, key := range result {
		if key != expected[i] {
			t.Errorf("Clean content mismatch, expected %v, got %v", expected[i], key)
		}
	}
}
