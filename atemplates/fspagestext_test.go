package atemplates

import (
	"bytes"
	"io/fs"
	"os"
	"path"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

// TestFSPagesText_Render tests the Render method of FSPagesText.
func TestFSPagesText_Render(t *testing.T) {
	// Mock function map.
	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
	}

	// Test cases.
	tests := []struct {
		name     string
		template string
		data     interface{}
		wantErr  bool
		reload   bool   // For ISON_PAGES_TEXT_LOADBYFILE
		useEmbed bool   // File vs. Embed
		wantStr  string // Expected substring in output
	}{
		{
			name:     "ValidTemplate_File_Preload",
			template: "home.gohtml",
			data:     struct{ Title string }{"Home Page"},
			wantErr:  false,
			reload:   false,
			useEmbed: false,
			wantStr:  "Home Page",
		},
		{
			name:     "InvalidTemplate_File_Preload",
			template: "missing.gohtml",
			data:     nil,
			wantErr:  true,
			reload:   false,
			useEmbed: false,
			wantStr:  "",
		},
		{
			name:     "ValidTemplate_File_Reload",
			template: "home.gohtml",
			data:     struct{ Title string }{"Home Page"},
			wantErr:  false,
			reload:   true,
			useEmbed: false,
			wantStr:  "Home Page",
		},
		{
			name:     "ValidTemplate_Embed_Preload",
			template: "home.gohtml",
			data:     struct{ Title string }{"Home Page"},
			wantErr:  false,
			reload:   false,
			useEmbed: true,
			wantStr:  "Home Page",
		},
		{
			name:     "InvalidTemplate_Embed_Preload",
			template: "missing.gohtml",
			data:     nil,
			wantErr:  true,
			reload:   false,
			useEmbed: true,
			wantStr:  "",
		},
		{
			name:     "ValidTemplate_Embed_Reload",
			template: "home.gohtml",
			data:     struct{ Title string }{"Home Page"},
			wantErr:  false,
			reload:   true,
			useEmbed: true,
			wantStr:  "Home Page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set reload flag.
			ISON_PAGES_TEXT_RELOAD_ON_RENDER = tt.reload

			var layoutsFS fs.FS
			var viewsFS fs.FS
			if tt.useEmbed {
				// Embedded FS.
				var err error
				layoutsFS, err = fs.Sub(testEmbedFS, "test_data/pages-text/layouts")
				if err != nil {
					t.Fatalf("Failed to sub layouts FS: %v", err)
				}
				viewsFS, err = fs.Sub(testEmbedFS, "test_data/pages-text/views")
				if err != nil {
					t.Fatalf("Failed to sub views FS: %v", err)
				}
			} else {
				// File-based FS.
				baseDir := "test_data"
				layoutsFS = os.DirFS(path.Join(baseDir, "pages-text/layouts"))
				viewsFS = os.DirFS(path.Join(baseDir, "pages-text/views"))
			}

			pages, err := NewFSPagesText(layoutsFS, []fs.FS{viewsFS}, &funcMap)
			assert.NoError(t, err)

			// Render to buffer.
			var buf bytes.Buffer
			err = pages.Render(&buf, tt.template, tt.data, mockEchoContext())

			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Check for expected substring in output.
			if err == nil && !strings.Contains(buf.String(), tt.wantStr) {
				t.Errorf("Render() got = %v, want containing %v", buf.String(), tt.wantStr)
			}
		})
	}
}

// TestFSPagesText_RenderPage tests the RenderPage method of FSPagesText.
func TestFSPagesText_RenderPage(t *testing.T) {
	// Mock function map.
	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
	}

	// Test cases (same as Render for consistency).
	tests := []struct {
		name     string
		template string
		data     interface{}
		wantErr  bool
		reload   bool   // For ISON_PAGES_TEXT_LOADBYFILE
		useEmbed bool   // File vs. Embed
		wantStr  string // Expected substring in output
	}{
		{
			name:     "ValidTemplate_File_Preload",
			template: "home.gohtml",
			data:     struct{ Title string }{"Home Page"},
			wantErr:  false,
			reload:   false,
			useEmbed: false,
			wantStr:  "Home Page",
		},
		{
			name:     "InvalidTemplate_File_Preload",
			template: "missing.gohtml",
			data:     nil,
			wantErr:  true,
			reload:   false,
			useEmbed: false,
			wantStr:  "",
		},
		{
			name:     "ValidTemplate_File_Reload",
			template: "home.gohtml",
			data:     struct{ Title string }{"Home Page"},
			wantErr:  false,
			reload:   true,
			useEmbed: false,
			wantStr:  "Home Page",
		},
		{
			name:     "ValidTemplate_Embed_Preload",
			template: "home.gohtml",
			data:     struct{ Title string }{"Home Page"},
			wantErr:  false,
			reload:   false,
			useEmbed: true,
			wantStr:  "Home Page",
		},
		{
			name:     "InvalidTemplate_Embed_Preload",
			template: "missing.gohtml",
			data:     nil,
			wantErr:  true,
			reload:   false,
			useEmbed: true,
			wantStr:  "",
		},
		{
			name:     "ValidTemplate_Embed_Reload",
			template: "home.gohtml",
			data:     struct{ Title string }{"Home Page"},
			wantErr:  false,
			reload:   true,
			useEmbed: true,
			wantStr:  "Home Page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set reload flag.
			ISON_PAGES_TEXT_RELOAD_ON_RENDER = tt.reload

			var layoutsFS fs.FS
			var viewsFS fs.FS
			if tt.useEmbed {
				// Embedded FS.
				var err error
				layoutsFS, err = fs.Sub(testEmbedFS, "test_data/pages-text/layouts")
				if err != nil {
					t.Fatalf("Failed to sub layouts FS: %v", err)
				}
				viewsFS, err = fs.Sub(testEmbedFS, "test_data/pages-text/views")
				if err != nil {
					t.Fatalf("Failed to sub views FS: %v", err)
				}
			} else {
				// File-based FS.
				baseDir := "test_data"
				layoutsFS = os.DirFS(path.Join(baseDir, "pages-text/layouts"))
				viewsFS = os.DirFS(path.Join(baseDir, "pages-text/views"))
			}

			pages, err := NewFSPagesText(layoutsFS, []fs.FS{viewsFS}, &funcMap)
			assert.NoError(t, err)

			// Call RenderPage.
			output, err := pages.RenderPage(tt.template, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("RenderPage() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Check for expected substring in output.
			if err == nil && !strings.Contains(output, tt.wantStr) {
				t.Errorf("RenderPage() got = %v, want containing %v", output, tt.wantStr)
			}
		})
	}
}
