package atemplates

import (
	"testing"
)

// TestTEMPLATES ensures that the TEMPLATES function returns a singleton instance of Manage.
func TestTEMPLATES(t *testing.T) {
	// Set up the initial instance with mock implementations.
	original := &Manage{
		PagesText:    &TPagesText{},
		PagesHTML:    &TPagesHTML{},
		SnippetsHTML: &TSnippetsHTML{},
		SnippetsText: &TSnippetsText{},
	}
	SetTemplates(original)

	// Retrieve the instance and verify it's the same as the original.
	retrieved := TEMPLATES()
	if retrieved != original {
		t.Errorf("TEMPLATES() did not return the original instance")
	}
}
