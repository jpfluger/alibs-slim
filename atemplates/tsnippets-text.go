package atemplates

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

// ISON_SNIPPETS_TEXT_LOADBYFILE indicates whether to load snippets from file on each request.
// Deprecated
var ISON_SNIPPETS_TEXT_LOADBYFILE = false

// TSnippetsText holds the templates and related information for rendering snippets.
// Deprecated
type TSnippetsText struct {
	snippets map[string]*template.Template // Map of snippet names to their compiled templates.
	funcMap  *template.FuncMap             // Function map for template rendering, allows for custom functions.
	dir      string                        // Directory where snippet files are located.
	// snippetsFileMap links snippet names to their file paths, used when ISON_SNIPPETS_TEXT_LOADBYFILE is true.
	snippetsFileMap map[string]string
}

// NewTSnippetsText creates a new TSnippetsText instance, loading snippet templates from the specified directory.
// Deprecated
func NewTSnippetsText(dir string, funcMap *template.FuncMap) *TSnippetsText {
	// Use common template functions if no specific function map is provided.
	if funcMap == nil {
		funcMap = GetTextTemplateFunctions(TEMPLATE_FUNCTIONS_COMMON)
	}

	// Ensure the directory path ends with a slash.
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	// Join the directory with the "snippets-text" subdirectory.
	dir = path.Join(dir, "snippets-text")

	// Check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(fmt.Errorf("Failed to create directory %s: %v", dir, err))
		}
	}

	// Collect all .gohtml files from the snippets directory.
	files := []string{}
	fnFiles := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".gohtml") {
			files = append(files, path)
		}
		return nil
	}

	// Walk through the directory to find snippet files.
	if err := filepath.Walk(dir, fnFiles); err != nil {
		panic(err) // Panic if walking the directory fails.
	}

	// Initialize the map for file paths and the map for compiled templates.
	snippetsFileMap := make(map[string]string)
	templates := make(map[string]*template.Template)

	// Compile templates for each snippet file found.
	for _, file := range files {
		name := path.Base(file)
		templates[name] = template.Must(template.New(name).Funcs(*funcMap).ParseFiles(file))

		// If snippets should be loaded by file, map the snippet names to their file paths.
		if ISON_SNIPPETS_TEXT_LOADBYFILE {
			snippetsFileMap[name] = file
		}
	}

	// Return a new TSnippetsText instance with the compiled templates and file map.
	return &TSnippetsText{
		snippets:        templates,
		funcMap:         funcMap,
		dir:             dir,
		snippetsFileMap: snippetsFileMap,
	}
}

// RenderSnippet renders a snippet with the given name and data.
func (t *TSnippetsText) RenderSnippet(name string, data interface{}) (string, error) {
	var buff bytes.Buffer // Buffer to hold the rendered template output.

	// If snippets should be loaded from file, read the file and execute the template.
	if ISON_SNIPPETS_TEXT_LOADBYFILE {
		snippetsPath, ok := t.snippetsFileMap[name]
		if !ok {
			return "", fmt.Errorf("snippets file map not found for '%s'", name)
		}

		// Parse the snippet file and execute the template with the provided data.
		tmpl := template.Must(template.New(name).Funcs(*t.funcMap).ParseFiles(snippetsPath))
		err := tmpl.Execute(&buff, data)
		if err != nil {
			return "", err
		}
		return buff.String(), nil
	}

	// If snippets are not loaded from file, use the precompiled template.
	tmpl, ok := t.snippets[name]
	if !ok {
		return "", fmt.Errorf("snippet '%s' not found", name)
	}

	// Execute the precompiled template with the provided data.
	err := tmpl.Execute(&buff, data)
	if err != nil {
		return "", fmt.Errorf("failed execute of snippet '%s'; %v", name, err)
	}
	return buff.String(), nil
}

// IsLoaded checks if the target snippet has been loaded as a template.
func (t *TSnippetsText) IsLoaded(name string) error {
	if ISON_SNIPPETS_TEXT_LOADBYFILE {
		_, ok := t.snippetsFileMap[name]
		if !ok {
			return fmt.Errorf("snippets text file map not found for '%s'", name)
		}
	} else {
		// If snippets are not loaded from file, use the precompiled template.
		_, ok := t.snippets[name]
		if !ok {
			fmt.Errorf("snippet text '%s' not found", name)
		}
	}
	return nil
}
