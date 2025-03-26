package asessions

import (
	"github.com/jpfluger/alibs-slim/azb"
	"strings"
)

// ActionKey is a type that represents a key for an action.
type ActionKey string

// IsEmpty checks if the ActionKey is empty after trimming whitespace.
func (ak ActionKey) IsEmpty() bool {
	return ak.TrimSpace() == ""
}

// TrimSpace trims whitespace from the ActionKey.
func (ak ActionKey) TrimSpace() ActionKey {
	rt := strings.TrimSpace(string(ak))
	return ActionKey(rt)
}

// String converts the ActionKey to a string.
func (ak ActionKey) String() string {
	return string(ak)
}

// ToZBType converts the ActionKey to a ZBType defined in the azb package.
func (ak ActionKey) ToZBType() azb.ZBType {
	return azb.ZBType(ak)
}

// ActionKeys is a slice of ActionKey.
type ActionKeys []ActionKey

// Find searches for a given key within the ActionKeys slice.
func (aks ActionKeys) Find(key ActionKey) ActionKey {
	if key.IsEmpty() {
		return ""
	}
	for _, ak := range aks {
		if ak == key {
			return ak
		}
	}
	return ""
}

// Has checks if a given key exists within the ActionKeys slice.
func (aks ActionKeys) Has(key ActionKey) bool {
	return aks.HasKey(key)
}

// HasKey is the long version of Has, checking for the existence of a key.
func (aks ActionKeys) HasKey(key ActionKey) bool {
	return !aks.Find(key).IsEmpty()
}

// Add appends a new key to the ActionKeys slice if it doesn't already exist.
func (aks *ActionKeys) Add(key ActionKey) {
	if key.IsEmpty() || aks.HasKey(key) {
		return
	}
	*aks = append(*aks, key)
}

// Remove deletes keys from the ActionKeys slice.
func (aks ActionKeys) Remove(keys ...ActionKey) ActionKeys {
	newActions := ActionKeys{}
	for _, ak := range aks {
		if !ak.IsEmpty() && !containsActionKey(keys, ak) {
			newActions = append(newActions, ak)
		}
	}
	return newActions
}

// containsActionKey checks if a given key is in the slice of ActionKey.
func containsActionKey(slice []ActionKey, key ActionKey) bool {
	for _, item := range slice {
		if item == key {
			return true
		}
	}
	return false
}
