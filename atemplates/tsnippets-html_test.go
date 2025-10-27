package atemplates

import (
	"html/template"
	"strings"
	"testing"
)

// TestRenderSnippetHTML tests the RenderSnippet method for various scenarios.
func TestRenderSnippetHTML(t *testing.T) {
	// Mock function map and directory for testing.
	funcMap := template.FuncMap{
		// Define any necessary template functions here.
		// For example, a function to uppercase strings:
		"upper": strings.ToUpper,
	}
	dir := "test_data/"

	// This must be set to true even if "loadFile" is false, so
	// that a reference to the file is created inside TSnippetsHTML.
	ISON_SNIPPETS_HTML_LOADBYFILE = true

	// Initialize TSnippetsHTML with the mock data.
	snippets := NewTSnippetsHTML(dir, &funcMap)

	// Define test cases.
	tests := []struct {
		name     string
		snippet  string
		data     interface{}
		want     string
		wantErr  bool
		loadFile bool
	}{
		{
			name:    "ValidSnippet",
			snippet: "hello.gohtml",
			data:    struct{ Name string }{"World"},
			want:    "Hello, World!",
			wantErr: false,
		},
		{
			name:    "InvalidSnippet",
			snippet: "nonexistent.gohtml",
			data:    nil,
			want:    "",
			wantErr: true,
		},
		{
			name:     "LoadFromFile",
			snippet:  "hello.gohtml",
			data:     struct{ Name string }{"File"},
			want:     "Hello, File!",
			wantErr:  false,
			loadFile: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the global variable for loading from file based on the test case.
			ISON_SNIPPETS_HTML_LOADBYFILE = tt.loadFile

			got, err := snippets.RenderSnippet(tt.snippet, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderSnippet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(got, tt.want) {
				t.Errorf("RenderSnippet() got = %v, want %v", got, tt.want)
			}
		})
	}
}
