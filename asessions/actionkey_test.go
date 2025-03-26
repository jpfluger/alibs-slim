package asessions

import (
	"testing"

	"github.com/jpfluger/alibs-slim/azb"
)

func TestActionKey_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		ak   ActionKey
		want bool
	}{
		{"empty", "", true},
		{"spaces only", "   ", true},
		{"non-empty", "key", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ak.IsEmpty(); got != tt.want {
				t.Errorf("ActionKey.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionKey_TrimSpace(t *testing.T) {
	tests := []struct {
		name string
		ak   ActionKey
		want ActionKey
	}{
		{"leading and trailing spaces", "  key  ", "key"},
		{"no spaces", "key", "key"},
		{"all spaces", "   ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ak.TrimSpace(); got != tt.want {
				t.Errorf("ActionKey.TrimSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionKey_String(t *testing.T) {
	ak := ActionKey("key")
	want := "key"
	if got := ak.String(); got != want {
		t.Errorf("ActionKey.String() = %v, want %v", got, want)
	}
}

func TestActionKey_ToZBType(t *testing.T) {
	ak := ActionKey("key")
	want := azb.ZBType("key")
	if got := ak.ToZBType(); got != want {
		t.Errorf("ActionKey.ToZBType() = %v, want %v", got, want)
	}
}

func TestActionKeys_Find(t *testing.T) {
	aks := ActionKeys{"key1", "key2", "key3"}
	tests := []struct {
		name string
		key  ActionKey
		want ActionKey
	}{
		{"found", "key2", "key2"},
		{"not found", "key4", ""},
		{"empty key", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aks.Find(tt.key); got != tt.want {
				t.Errorf("ActionKeys.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionKeys_Has(t *testing.T) {
	aks := ActionKeys{"key1", "key2", "key3"}
	tests := []struct {
		name string
		key  ActionKey
		want bool
	}{
		{"has key", "key2", true},
		{"does not have key", "key4", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aks.Has(tt.key); got != tt.want {
				t.Errorf("ActionKeys.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionKeys_Add(t *testing.T) {
	aks := &ActionKeys{"key1", "key2"}
	aks.Add("key3")
	want := ActionKeys{"key1", "key2", "key3"}

	if !equalActionKeys(*aks, want) {
		t.Errorf("ActionKeys.Add() = %v, want %v", *aks, want)
	}

	// Test adding existing key
	aks.Add("key2")
	want = ActionKeys{"key1", "key2", "key3"} // should be unchanged

	if !equalActionKeys(*aks, want) {
		t.Errorf("ActionKeys.Add() did not ignore existing key, got %v, want %v", *aks, want)
	}
}

func TestActionKeys_Remove(t *testing.T) {
	aks := ActionKeys{"key1", "key2", "key3"}
	aks = aks.Remove("key2", "key3")
	want := ActionKeys{"key1"}

	if !equalActionKeys(aks, want) {
		t.Errorf("ActionKeys.Remove() = %v, want %v", aks, want)
	}
}

// Helper function to compare two ActionKeys slices
func equalActionKeys(a, b ActionKeys) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
