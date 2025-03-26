package ascripts

import (
	"github.com/jpfluger/alibs-slim/ajson" // Importing a custom JSON utility package
)

// ScriptablesMap is a map that associates ajson.JsonKey with a pointer to a Scriptable
type ScriptablesMap map[ajson.JsonKey]*Scriptable

// Find retrieves a *Scriptable associated with the provided JsonKey, if it exists
func (sim ScriptablesMap) Find(key ajson.JsonKey) *Scriptable {
	// Check if the map is nil, empty, or the key is empty
	if sim == nil || len(sim) == 0 || key.IsEmpty() {
		return nil // Return nil if any of the conditions are true
	}
	val, ok := sim[key] // Attempt to retrieve the value
	if !ok {
		return nil // Return nil if the key does not exist in the map
	}
	return val // Return the retrieved *Scriptable
}

// HasItem checks if a *Scriptable associated with the provided JsonKey exists in the map
func (sim ScriptablesMap) HasItem(key ajson.JsonKey) bool {
	return sim.Find(key) != nil // Utilize the Find method to check for existence
}
