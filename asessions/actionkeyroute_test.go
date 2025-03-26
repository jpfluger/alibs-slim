package asessions

import (
	"reflect"
	"testing"
)

// TestActionKeyUrls_Find tests the Find method of ActionKeyUrls.
func TestActionKeyUrls_Find(t *testing.T) {
	ak1 := &ActionKeyUrl{Key: "key1", Url: "http://example.com/1", CheckPrefix: false}
	ak2 := &ActionKeyUrl{Key: "key2", Url: "http://example.com/2", CheckPrefix: true}
	aks := ActionKeyUrls{ak1, ak2}

	tests := []struct {
		name string
		key  ActionKey
		want *ActionKeyUrl
	}{
		{"find existing key", "key1", ak1},
		{"find non-existing key", "key3", nil},
		{"find empty key", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aks.Find(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ActionKeyUrls.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestActionKeyUrls_Has tests the Has method of ActionKeyUrls.
func TestActionKeyUrls_Has(t *testing.T) {
	ak1 := &ActionKeyUrl{Key: "key1"}
	aks := ActionKeyUrls{ak1}

	tests := []struct {
		name string
		key  ActionKey
		want bool
	}{
		{"has existing key", "key1", true},
		{"has non-existing key", "key2", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aks.Has(tt.key); got != tt.want {
				t.Errorf("ActionKeyUrls.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestActionKeyUrls_Add tests the Add method of ActionKeyUrls.
func TestActionKeyUrls_Add(t *testing.T) {
	ak1 := &ActionKeyUrl{Key: "key1"}
	aks := &ActionKeyUrls{ak1}

	akToAdd := &ActionKeyUrl{Key: "key2"}
	aks.Add(akToAdd)
	want := &ActionKeyUrls{ak1, akToAdd}

	if !reflect.DeepEqual(aks, want) {
		t.Errorf("ActionKeyUrls.Add() = %v, want %v", aks, want)
	}

	// Test adding an existing key
	aks.Add(ak1)
	if !reflect.DeepEqual(aks, want) { // should be unchanged
		t.Errorf("ActionKeyUrls.Add() did not ignore existing key, got %v, want %v", aks, want)
	}
}

// TestActionKeyUrls_Remove tests the Remove function of ActionKeyUrls.
func TestActionKeyUrls_Remove(t *testing.T) {
	ak1 := &ActionKeyUrl{Key: "key1"}
	ak2 := &ActionKeyUrl{Key: "key2"}
	aks := ActionKeyUrls{ak1, ak2}

	aks = aks.Remove("key1")
	want := ActionKeyUrls{ak2}

	if !reflect.DeepEqual(aks, want) {
		t.Errorf("ActionKeyUrls.Remove() = %v, want %v", aks, want)
	}
}
