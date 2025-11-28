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

// ISON_SNIPPETS_TEXT_RELOAD_ON_RENDER if true, re-parses snippets from FS on each render (dev/hot-reload); else preloads.
var ISON_SNIPPETS_TEXT_RELOAD_ON_RENDER = false

// FSSnippetsText holds the templates and related information for rendering snippets from FS.
type FSSnippetsText struct {
	snippets   map[string]*template.Template // Preloaded templates (if not reload mode).
	funcMap    *template.FuncMap             // Function map for template rendering.
	snippetsFS []fs.FS                       // FS slice for snippets (e.g., from multiple embed packages).
}

// NewFSSnippetsText creates a new FSSnippetsText instance, loading snippet templates from the provided FS sources.
func NewFSSnippetsText(snippetsFS []fs.FS, funcMap *template.FuncMap) (*FSSnippetsText, error) {
	// Use common template functions if no specific function map is provided.
	if funcMap == nil {
		funcMap = GetTextTemplateFunctions(TEMPLATE_FUNCTIONS_COMMON)
	}

	fst := &FSSnippetsText{
		funcMap:    funcMap,
		snippetsFS: snippetsFS,
	}

	// Always preload for verification.
	snippets, err := fst.preloadSnippets()
	if err != nil {
		return nil, fmt.Errorf("failed to preload snippets: %w", err)
	}
	fst.snippets = snippets

	// If reload enabled, clear preloaded map to force re-parse on each render.
	if ISON_SNIPPETS_TEXT_LOADBYFILE || ISON_SNIPPETS_TEXT_RELOAD_ON_RENDER {
		fst.snippets = make(map[string]*template.Template)
	}

	return fst, nil
}

// preloadSnippets loads and compiles all snippets in-memory from multiple FS sources.
func (t *FSSnippetsText) preloadSnippets() (map[string]*template.Template, error) {
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

	if len(snippets) == 0 && len(t.snippetsFS) > 0 {
		return nil, fmt.Errorf("no snippets found in any FS")
	}

	return snippets, nil
}

// loadFilesFromFS walks the FS and loads file contents into a map (name: content).
func (t *FSSnippetsText) loadFilesFromFS(fsys fs.FS, ext string) (map[string]string, error) {
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
func (t *FSSnippetsText) RenderSnippet(name string, data interface{}) (string, error) {
	var buff bytes.Buffer

	if ISON_SNIPPETS_TEXT_LOADBYFILE || ISON_SNIPPETS_TEXT_RELOAD_ON_RENDER {
		// Reload from multiple FS sources.
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
				break // First found wins; adjust if needed
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
func (t *FSSnippetsText) IsLoaded(name string) error {
	if ISON_SNIPPETS_TEXT_LOADBYFILE || ISON_SNIPPETS_TEXT_RELOAD_ON_RENDER {
		// In reload mode, walk each FS to check if a file with base name 'name' exists.
		for _, sFS := range t.snippetsFS {
			if sFS == nil {
				continue
			}
			found := false
			err := fs.WalkDir(sFS, ".", func(filePath string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				if strings.HasSuffix(filePath, ".gohtml") && path.Base(filePath) == name {
					found = true
					return fs.SkipAll // Stop walking once found
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to check snippet '%s' in FS: %w", name, err)
			}
			if found {
				return nil // Found in this FS
			}
		}
		return fmt.Errorf("snippet text '%s' not found in any FS", name)
	}
	// Preloaded mode.
	_, ok := t.snippets[name]
	if !ok {
		return fmt.Errorf("snippet text '%s' not found", name)
	}
	return nil
}
