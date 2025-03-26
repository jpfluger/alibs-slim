package atemplates

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

// mockEchoContext creates a mock Echo context for testing purposes.
func mockEchoContext() echo.Context {
	e := echo.New()
	req := e.NewContext(nil, nil).Request()
	return e.NewContext(req, nil)
}

// TestTPagesText_Render tests the Render method of TPagesText.
func TestTPagesText_Render(t *testing.T) {
	// Mock function map and directory for testing.
	funcMap := template.FuncMap{
		// Define any necessary template functions here.
		// For example, a function to uppercase strings:
		"upper": strings.ToUpper,
	}
	dir := "test_files/"

	// This must be set to true even if "loadFile" is false, so
	// that a reference to the file is created inside TSnippetsText.
	ISON_PAGES_TEXT_LOADBYFILE = true

	// Initialize TPagesText with the mock data.
	pages := NewTPagesText(dir, &funcMap)

	// Define test cases.
	tests := []struct {
		name     string
		template string
		data     interface{}
		wantErr  bool
		loadFile bool
	}{
		{
			name:     "ValidTemplate",
			template: "home.gohtml",
			data:     struct{ Title string }{"Home Page"},
			wantErr:  false,
			loadFile: false,
		},
		{
			name:     "InvalidTemplate",
			template: "missing.gohtml",
			data:     nil,
			wantErr:  true,
			loadFile: false,
		},
		{
			name:     "LoadFromFile",
			template: "home.gohtml",
			data:     struct{ Title string }{"Home Page"},
			wantErr:  false,
			loadFile: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the global variable for loading from file based on the test case.
			ISON_PAGES_TEXT_LOADBYFILE = tt.loadFile

			// Create a buffer to write the output to.
			var buf bytes.Buffer

			// Call the Render method and capture any errors.
			err := pages.Render(&buf, tt.template, tt.data, mockEchoContext())

			// Check for unexpected errors or lack of expected errors.
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestTPagesText_RenderPage tests the RenderPage method of TPagesText.
func TestTPagesText_RenderPage(t *testing.T) {
	// Mock function map and directory for testing.
	funcMap := template.FuncMap{
		// Define any necessary template functions here.
		// For example, a function to uppercase strings:
		"upper": strings.ToUpper,
	}
	dir := "test_files/"

	// This must be set to true even if "loadFile" is false, so
	// that a reference to the file is created inside TSnippetsText.
	ISON_PAGES_TEXT_LOADBYFILE = true

	// Initialize TPagesText with the mock data.
	pages := NewTPagesText(dir, &funcMap)

	// Define test cases.
	tests := []struct {
		name     string
		template string
		data     interface{}
		wantErr  bool
		loadFile bool
	}{
		{
			name:     "ValidTemplate",
			template: "home.gohtml",
			data:     struct{ Title string }{"Home Page"},
			wantErr:  false,
			loadFile: false,
		},
		{
			name:     "InvalidTemplate",
			template: "missing.gohtml",
			data:     nil,
			wantErr:  true,
			loadFile: false,
		},
		{
			name:     "LoadFromFile",
			template: "home.gohtml",
			data:     struct{ Title string }{"Home Page"},
			wantErr:  false,
			loadFile: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the global variable for loading from file based on the test case.
			ISON_PAGES_TEXT_LOADBYFILE = tt.loadFile

			// Call the RenderPage method and capture any errors.
			_, err := pages.RenderPage(tt.template, tt.data)

			// Check for unexpected errors or lack of expected errors.
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderPage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
