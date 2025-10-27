package atemplates

//import (
//	"bytes"
//	"html/template"
//	"strings"
//	"testing"
//)
//
//// TestTPagesHTML_Render tests the Render method of TPagesHTML.
//func TestTPagesHTML_Render(t *testing.T) {
//	// Mock function map and directory for testing.
//	funcMap := template.FuncMap{
//		// Define any necessary template functions here.
//		// For example, a function to uppercase strings:
//		"upper": strings.ToUpper,
//	}
//	dir := "test_data/"
//
//	// This must be set to true even if "loadFile" is false, so
//	// that a reference to the file is created inside TSnippetsText.
//	ISON_PAGES_HTML_LOADBYFILE = true
//
//	// Initialize TPagesHTML with the mock data.
//	pages := NewTPagesHTML(dir, &funcMap)
//
//	// Define test cases.
//	tests := []struct {
//		name     string
//		template string
//		data     interface{}
//		wantErr  bool
//		loadFile bool
//	}{
//		{
//			name:     "ValidTemplate",
//			template: "home.gohtml",
//			data:     struct{ Title string }{"Home Page"},
//			wantErr:  false,
//			loadFile: false,
//		},
//		{
//			name:     "InvalidTemplate",
//			template: "missing.gohtml",
//			data:     nil,
//			wantErr:  true,
//			loadFile: false,
//		},
//		{
//			name:     "LoadFromFile",
//			template: "home.gohtml",
//			data:     struct{ Title string }{"Home Page"},
//			wantErr:  false,
//			loadFile: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			// Set the global variable for loading from file based on the test case.
//			ISON_PAGES_HTML_LOADBYFILE = tt.loadFile
//
//			// Create a buffer to write the output to.
//			var buf bytes.Buffer
//
//			// Call the Render method and capture any errors.
//			err := pages.Render(&buf, tt.template, tt.data, mockEchoContext())
//
//			// Check for unexpected errors or lack of expected errors.
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//// TestTPagesHTML_RenderPage tests the RenderPage method of TPagesHTML.
//func TestTPagesHTML_RenderPage(t *testing.T) {
//	// Mock function map and directory for testing.
//	funcMap := template.FuncMap{
//		// Define any necessary template functions here.
//		// For example, a function to uppercase strings:
//		"upper": strings.ToUpper,
//	}
//	dir := "test_data/"
//
//	// This must be set to true even if "loadFile" is false, so
//	// that a reference to the file is created inside TSnippetsText.
//	ISON_PAGES_HTML_LOADBYFILE = true
//
//	// Initialize TPagesHTML with the mock data.
//	pages := NewTPagesHTML(dir, &funcMap)
//
//	// Define test cases.
//	tests := []struct {
//		name     string
//		template string
//		data     interface{}
//		wantErr  bool
//		loadFile bool
//	}{
//		{
//			name:     "ValidTemplate",
//			template: "home.gohtml",
//			data:     struct{ Title string }{"Home Page"},
//			wantErr:  false,
//			loadFile: false,
//		},
//		{
//			name:     "InvalidTemplate",
//			template: "missing.gohtml",
//			data:     nil,
//			wantErr:  true,
//			loadFile: false,
//		},
//		{
//			name:     "LoadFromFile",
//			template: "home.gohtml",
//			data:     struct{ Title string }{"Home Page"},
//			wantErr:  false,
//			loadFile: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			// Set the global variable for loading from file based on the test case.
//			ISON_PAGES_HTML_LOADBYFILE = tt.loadFile
//
//			// Call the RenderPage method and capture any errors.
//			_, err := pages.RenderPage(tt.template, tt.data)
//
//			// Check for unexpected errors or lack of expected errors.
//			if (err != nil) != tt.wantErr {
//				t.Errorf("RenderPage() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
