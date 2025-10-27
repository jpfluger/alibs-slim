package atemplates

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

//go:embed test_data/*
var testEmbedFS embed.FS

// mockEchoContext returns a mock Echo context.
func mockEchoContext() echo.Context {
	return echo.New().NewContext(nil, nil)
}

// TestFSPagesHTML_Render tests the Render method of FSPagesHTML.
func TestFSPagesHTML_Render(t *testing.T) {
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
		reload   bool   // For ISON_PAGES_HTML_LOADBYFILE
		useEmbed bool   // File vs. Embed
		wantStr  string // Expected substring in output
	}{
		{
			name:     "ValidTemplate_File_Preload",
			template: "home.gohtml",
			data:     struct{ Title, Header, MainHeading, Content, Footer string }{"Home Page", "Test Header", "Main", "Test Content", "Test Footer"},
			wantErr:  false,
			reload:   false,
			useEmbed: false,
			wantStr:  "Home Page", // From title block
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
			name:     "SFPPrefix_File_Preload",
			template: "sfp-special.gohtml",
			data:     struct{ Title, Content string }{"Special Page", "Special Content"},
			wantErr:  false,
			reload:   false,
			useEmbed: false,
			wantStr:  "Special Standalone Content", // Standalone check
		},
		{
			name:     "ValidTemplate_File_Reload",
			template: "home.gohtml",
			data:     struct{ Title, Header, MainHeading, Content, Footer string }{"Home Page", "Test Header", "Main", "Test Content", "Test Footer"},
			wantErr:  false,
			reload:   true,
			useEmbed: false,
			wantStr:  "Home Page",
		},
		{
			name:     "ValidTemplate_Embed_Preload",
			template: "home.gohtml",
			data:     struct{ Title, Header, MainHeading, Content, Footer string }{"Home Page", "Test Header", "Main", "Test Content", "Test Footer"},
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
			name:     "SFPPrefix_Embed_Preload",
			template: "sfp-special.gohtml",
			data:     struct{ Title, Content string }{"Special Page", "Special Content"},
			wantErr:  false,
			reload:   false,
			useEmbed: true,
			wantStr:  "Special Standalone Content",
		},
		{
			name:     "ValidTemplate_Embed_Reload",
			template: "home.gohtml",
			data:     struct{ Title, Header, MainHeading, Content, Footer string }{"Home Page", "Test Header", "Main", "Test Content", "Test Footer"},
			wantErr:  false,
			reload:   true,
			useEmbed: true,
			wantStr:  "Home Page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set reload flag.
			ISON_PAGES_HTML_RELOAD_ON_RENDER = tt.reload

			var layoutsFS fs.FS
			var viewsFS fs.FS
			if tt.useEmbed {
				// Embedded FS.
				var err error
				layoutsFS, err = fs.Sub(testEmbedFS, "test_data/pages-html/layouts")
				if err != nil {
					t.Fatalf("Failed to sub layouts FS: %v", err)
				}
				viewsFS, err = fs.Sub(testEmbedFS, "test_data/pages-html/views")
				if err != nil {
					t.Fatalf("Failed to sub views FS: %v", err)
				}
			} else {
				// File-based FS.
				baseDir := "test_data"
				layoutsFS = os.DirFS(path.Join(baseDir, "pages-html/layouts"))
				viewsFS = os.DirFS(path.Join(baseDir, "pages-html/views"))
			}

			pages, err := NewFSPagesHTML(layoutsFS, []fs.FS{viewsFS}, &funcMap)
			assert.NoError(t, err)

			// Render to buffer.
			var buf bytes.Buffer
			err = pages.Render(&buf, tt.template, tt.data, mockEchoContext())

			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Check for expected substring in output.
			if err == nil && !strings.Contains(buf.String(), tt.wantStr) {
				t.Errorf("Unexpected output: got %s, want containing %s", buf.String(), tt.wantStr)
			}
		})
	}
}

// TestFSPagesHTML_RenderPage tests the RenderPage method of FSPagesHTML.
func TestFSPagesHTML_RenderPage(t *testing.T) {
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
		reload   bool   // For ISON_PAGES_HTML_LOADBYFILE
		useEmbed bool   // File vs. Embed
		wantStr  string // Expected substring in output
	}{
		// Copy the tests array from TestFSPagesHTML_Render...
		{
			name:     "ValidTemplate_File_Preload",
			template: "home.gohtml",
			data:     struct{ Title, Header, MainHeading, Content, Footer string }{"Home Page", "Test Header", "Main", "Test Content", "Test Footer"},
			wantErr:  false,
			reload:   false,
			useEmbed: false,
			wantStr:  "Home Page",
		},
		// ... (add all other cases similarly)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set reload flag.
			ISON_PAGES_HTML_RELOAD_ON_RENDER = tt.reload

			var layoutsFS fs.FS
			var viewsFS fs.FS
			if tt.useEmbed {
				// Embedded FS.
				var err error
				layoutsFS, err = fs.Sub(testEmbedFS, "test_data/pages-html/layouts")
				if err != nil {
					t.Fatalf("Failed to sub layouts FS: %v", err)
				}
				viewsFS, err = fs.Sub(testEmbedFS, "test_data/pages-html/views")
				if err != nil {
					t.Fatalf("Failed to sub views FS: %v", err)
				}
			} else {
				// File-based FS.
				baseDir := "test_data"
				layoutsFS = os.DirFS(path.Join(baseDir, "pages-html/layouts"))
				viewsFS = os.DirFS(path.Join(baseDir, "pages-html/views"))
			}

			pages, err := NewFSPagesHTML(layoutsFS, []fs.FS{viewsFS}, &funcMap)
			assert.NoError(t, err)

			// Call RenderPage.
			output, err := pages.RenderPage(tt.template, tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("RenderPage() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Check for expected substring in output.
			if err == nil && !strings.Contains(output, tt.wantStr) {
				t.Errorf("Unexpected output: got %s, want containing %s", output, tt.wantStr)
			}
		})
	}
}
