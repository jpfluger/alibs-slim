package atemplates

import (
	"html/template"
	"io/fs"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFSSnippetsHTML_RenderSnippet tests the RenderSnippet method for various scenarios.
func TestFSSnippetsHTML_RenderSnippet(t *testing.T) {
	// Mock function map.
	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
	}

	// Test cases.
	tests := []struct {
		name     string
		snippet  string
		data     interface{}
		wantErr  bool
		reload   bool   // For ISON_SNIPPETS_HTML_LOADBYFILE
		useEmbed bool   // File vs. Embed
		wantStr  string // Expected substring in output
	}{
		{
			name:     "ValidSnippet_File_Preload",
			snippet:  "hello.gohtml",
			data:     struct{ Name string }{"World"},
			wantErr:  false,
			reload:   false,
			useEmbed: false,
			wantStr:  "Hello, World!",
		},
		{
			name:     "InvalidSnippet_File_Preload",
			snippet:  "nonexistent.gohtml",
			data:     nil,
			wantErr:  true,
			reload:   false,
			useEmbed: false,
			wantStr:  "",
		},
		{
			name:     "ValidSnippet_File_Reload",
			snippet:  "hello.gohtml",
			data:     struct{ Name string }{"File"},
			wantErr:  false,
			reload:   true,
			useEmbed: false,
			wantStr:  "Hello, File!",
		},
		{
			name:     "ValidSnippet_Embed_Preload",
			snippet:  "hello.gohtml",
			data:     struct{ Name string }{"World"},
			wantErr:  false,
			reload:   false,
			useEmbed: true,
			wantStr:  "Hello, World!",
		},
		{
			name:     "InvalidSnippet_Embed_Preload",
			snippet:  "nonexistent.gohtml",
			data:     nil,
			wantErr:  true,
			reload:   false,
			useEmbed: true,
			wantStr:  "",
		},
		{
			name:     "ValidSnippet_Embed_Reload",
			snippet:  "hello.gohtml",
			data:     struct{ Name string }{"Embed"},
			wantErr:  false,
			reload:   true,
			useEmbed: true,
			wantStr:  "Hello, Embed!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set reload flag.
			ISON_SNIPPETS_HTML_RELOAD_ON_RENDER = tt.reload

			var snippetsFS fs.FS
			if tt.useEmbed {
				// Embedded FS - sub to "snippets-html" (since //go:embed test_data/* puts "snippets-html/" at root).
				var err error
				snippetsFS, err = fs.Sub(testEmbedFS, "test_data/snippets-html")
				if err != nil {
					t.Fatalf("Failed to sub snippets FS: %v", err)
				}
			} else {
				// File-based FS.
				baseDir := "test_data"
				snippetsFS = os.DirFS(path.Join(baseDir, "snippets-html"))
			}

			snippets, err := NewFSSnippetsHTML([]fs.FS{snippetsFS}, &funcMap)
			assert.NoError(t, err)

			// Call RenderSnippet.
			got, err := snippets.RenderSnippet(tt.snippet, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("RenderSnippet() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !strings.Contains(got, tt.wantStr) {
				t.Errorf("RenderSnippet() got = %v, want containing %v", got, tt.wantStr)
			}
		})
	}
}
