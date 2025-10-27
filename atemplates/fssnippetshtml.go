package atemplates

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"path"
	"strings"

	"github.com/jpfluger/alibs-slim/alog"
)

// ISON_SNIPPETS_HTML_RELOAD_ON_RENDER if true, re-parses snippets from FS on each render (dev/hot-reload); else preloads.
var ISON_SNIPPETS_HTML_RELOAD_ON_RENDER = false

// FSSnippetsHTML holds HTML templates and related data for rendering snippets from multiple FS sources.
type FSSnippetsHTML struct {
	snippets   map[string]*template.Template // Preloaded templates (if not reload mode).
	funcMap    *template.FuncMap             // Function map for template rendering.
	snippetsFS []fs.FS                       // FS slice for snippets (e.g., from multiple embed packages).
}

// NewFSSnippetsHTML creates a new FSSnippetsHTML instance, loading snippet templates from the provided FS sources.
func NewFSSnippetsHTML(snippetsFS []fs.FS, funcMap *template.FuncMap) (*FSSnippetsHTML, error) {
	// Use common HTML template functions if no specific function map is provided.
	if funcMap == nil {
		funcMap = GetHTMLTemplateFunctions(TEMPLATE_FUNCTIONS_COMMON)
	}

	fsh := &FSSnippetsHTML{
		funcMap:    funcMap,
		snippetsFS: snippetsFS,
	}

	// Always preload for verification.
	snippets, err := fsh.preloadSnippets()
	if err != nil {
		return nil, fmt.Errorf("failed to preload snippets: %w", err)
	}
	fsh.snippets = snippets

	// If reload enabled, clear preloaded map to force re-parse on each render.
	if ISON_SNIPPETS_HTML_RELOAD_ON_RENDER {
		fsh.snippets = make(map[string]*template.Template)
	}

	return fsh, nil
}

// preloadSnippets loads and compiles all snippets in-memory from multiple FS sources.
func (t *FSSnippetsHTML) preloadSnippets() (map[string]*template.Template, error) {
	snippets := make(map[string]*template.Template)

	// Load snippets from all FS sources.
	for _, sFS := range t.snippetsFS {
		if sFS == nil {
			continue // Skip nil FS
		}
		filesContent, err := t.loadFilesFromFS(sFS, ".gohtml")
		if err != nil {
			return nil, fmt.Errorf("failed to load snippets from FS: %w", err)
		}
		for name, content := range filesContent {
			if _, exists := snippets[name]; exists {
				alog.LOGGER(alog.LOGGER_APP).Warn().Msgf("Duplicate snippet '%s' found; overwriting with latest", name)
			}
			tmpl, err := template.New(name).Funcs(*t.funcMap).Parse(content)
			if err != nil {
				return nil, fmt.Errorf("failed to parse snippet '%s': %w", name, err)
			}
			snippets[name] = tmpl
		}
	}

	if len(snippets) == 0 {
		return nil, fmt.Errorf("no snippets found in any FS")
	}

	return snippets, nil
}

// loadFilesFromFS walks the FS and loads file contents into a map (name: content).
func (t *FSSnippetsHTML) loadFilesFromFS(fsys fs.FS, ext string) (map[string]string, error) {
	if fsys == nil {
		return nil, fmt.Errorf("FS is nil")
	}

	files := make(map[string]string)
	err := fs.WalkDir(fsys, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ext) {
			content, err := fs.ReadFile(fsys, filePath)
			if err != nil {
				return err
			}
			files[path.Base(filePath)] = string(content)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load files from FS: %w", err)
	}
	return files, nil
}

// RenderSnippet renders a snippet with the given name and data.
func (t *FSSnippetsHTML) RenderSnippet(name string, data interface{}) (string, error) {
	var buff bytes.Buffer

	if ISON_SNIPPETS_HTML_RELOAD_ON_RENDER {
		// Reload from FS.
		snippetContent := ""
		found := false
		for _, sFS := range t.snippetsFS {
			if sFS == nil {
				continue
			}
			filesContent, err := t.loadFilesFromFS(sFS, ".gohtml")
			if err != nil {
				return "", err
			}
			if content, ok := filesContent[name]; ok {
				snippetContent = content
				found = true
				break
			}
		}
		if !found {
			return "", fmt.Errorf("snippet '%s' not found in any FS", name)
		}
		tmpl, err := template.New(name).Funcs(*t.funcMap).Parse(snippetContent)
		if err != nil {
			return "", fmt.Errorf("failed to parse snippet '%s': %w", name, err)
		}
		err = tmpl.Execute(&buff, data)
		if err != nil {
			return "", fmt.Errorf("failed execute of snippet '%s': %w", name, err)
		}
		return buff.String(), nil
	}

	// Use preloaded.
	tmpl, ok := t.snippets[name]
	if !ok {
		return "", fmt.Errorf("snippet '%s' not found", name)
	}
	err := tmpl.Execute(&buff, data)
	if err != nil {
		return "", fmt.Errorf("failed execute of snippet '%s': %w", name, err)
	}
	return buff.String(), nil
}

// IsLoaded checks if the target snippet has been loaded as a template.
func (t *FSSnippetsHTML) IsLoaded(name string) error {
	if ISON_SNIPPETS_HTML_RELOAD_ON_RENDER {
		// In reload mode, check if file exists in any FS.
		for _, sFS := range t.snippetsFS {
			if sFS == nil {
				continue
			}
			_, err := fs.Stat(sFS, name)
			if err == nil {
				return nil // Found in this FS
			}
		}
		return fmt.Errorf("snippet html '%s' not found in any FS", name)
	}
	// Preloaded mode.
	_, ok := t.snippets[name]
	if !ok {
		return fmt.Errorf("snippet html '%s' not found", name)
	}
	return nil
}
