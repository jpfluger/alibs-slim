package ascripts

import (
	"testing"
)

// TestFindScriptTypeFilter checks if the Find method correctly retrieves a ScriptTypeFilter
func TestFindScriptTypeFilter(t *testing.T) {
	// Create a slice of ScriptTypeFilters with some default values
	filters := ScriptTypeFilters{
		&ScriptTypeFilter{Type: SCRIPTTYPE_GO, Title: "Go"},
		&ScriptTypeFilter{Type: SCRIPTTYPE_HTML, Title: "HTML"},
	}

	// Test finding an existing filter by type
	if filter := filters.Find(SCRIPTTYPE_GO); filter == nil || filter.Title != "Go" {
		t.Errorf("Find did not return the correct ScriptTypeFilter for SCRIPTTYPE_GO")
	}

	// Test finding a non-existing filter
	if filter := filters.Find(ScriptType("nonexistent")); filter != nil {
		t.Errorf("Find returned a ScriptTypeFilter for a non-existent type")
	}
}

// TestHasScriptTypeFilter checks if the Has method correctly identifies the existence of a ScriptTypeFilter
func TestHasScriptTypeFilter(t *testing.T) {
	// Create a slice of ScriptTypeFilters with some default values
	filters := ScriptTypeFilters{
		&ScriptTypeFilter{Type: SCRIPTTYPE_GO, Title: "Go"},
		&ScriptTypeFilter{Type: SCRIPTTYPE_HTML, Title: "HTML"},
	}

	// Test for an existing filter by type
	if !filters.Has(SCRIPTTYPE_GO) {
		t.Errorf("Has did not identify the existence of SCRIPTTYPE_GO")
	}

	// Test for a non-existing filter
	if filters.Has(ScriptType("nonexistent")) {
		t.Errorf("Has incorrectly identified the existence of a non-existent type")
	}
}

// TestGetScriptTypeFilterDefaults checks if default filters are created correctly
func TestGetScriptTypeFilterDefaults(t *testing.T) {
	// Create a slice of ScriptTypes
	scriptTypes := ScriptTypes{SCRIPTTYPE_GO, SCRIPTTYPE_HTML}

	// Get default filters for the script types
	defaultFilters := GetScriptTypeFilterDefaults(scriptTypes)

	// Check if the correct number of default filters are created
	if len(defaultFilters) != len(scriptTypes) {
		t.Errorf("Expected %d default filters, got %d", len(scriptTypes), len(defaultFilters))
	}

	// Check if the default filters have the correct types and titles
	for _, filter := range defaultFilters {
		if filter.Type == SCRIPTTYPE_GO && filter.Title != "Go" {
			t.Errorf("Default filter for SCRIPTTYPE_GO does not have the title 'Go'")
		}
		if filter.Type == SCRIPTTYPE_HTML && filter.Title != "HTML" {
			t.Errorf("Default filter for SCRIPTTYPE_HTML does not have the title 'HTML'")
		}
	}
}
