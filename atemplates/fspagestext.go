package atemplates

import (
	"bytes"
	"fmt"
	"github.com/jpfluger/alibs-slim/alog"
	"github.com/labstack/echo/v4"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"
)

// ISON_PAGES_TEXT_RELOAD_ON_RENDER if true, re-parses templates from FS on each render (dev/hot-reload); else preloads.
var ISON_PAGES_TEXT_RELOAD_ON_RENDER = false

// ISON_PAGES_TEXT_LOG_ERROR enables logging of errors during template execution.
var ISON_PAGES_TEXT_LOG_ERROR = false

// FSPagesText is an FS-based text/template renderer for the Echo framework.
type FSPagesText struct {
	templates map[string]*template.Template // Preloaded templates (if not reload mode).
	funcMap   *template.FuncMap             // Function map for template rendering.
	layoutsFS fs.FS                         // FS for layouts (optional: e.g., fs.Sub(embedweb.Files, "pages-text/layouts")).
	viewsFS   []fs.FS                       // FS slice for views (multiple sources: e.g., fs.Sub(embedweb.Files, "pages-text/views")).
}

// NewFSPagesText creates a new FSPagesText instance, loading templates from the provided FS sources.
func NewFSPagesText(layoutsFS fs.FS, viewsFS []fs.FS, funcMap *template.FuncMap) (*FSPagesText, error) {
	// Use common text template functions if no specific function map is provided.
	if funcMap == nil {
		funcMap = GetTextTemplateFunctions(TEMPLATE_FUNCTIONS_COMMON)
	}

	fpt := &FSPagesText{
		funcMap:   funcMap,
		layoutsFS: layoutsFS,
		viewsFS:   viewsFS,
	}

	// Always preload for verification.
	templates, err := fpt.preloadTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to preload templates: %w", err)
	}
	fpt.templates = templates

	// If reload enabled, clear preloaded map to force re-parse on each render.
	if ISON_PAGES_TEXT_LOADBYFILE || ISON_PAGES_TEXT_RELOAD_ON_RENDER {
		fpt.templates = make(map[string]*template.Template)
	}

	return fpt, nil
}

// preloadTemplates loads and compiles all templates in-memory.
func (t *FSPagesText) preloadTemplates() (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)

	// Load layouts (optional).
	layouts, err := t.loadFilesFromFS(t.layoutsFS, ".gohtml", true) // optional=true
	if err != nil {
		return nil, fmt.Errorf("failed to load layouts: %w", err)
	}

	// Load all views from multiple sources (required).
	views := make(map[string]string)
	for _, vFS := range t.viewsFS {
		subViews, err := t.loadFilesFromFS(vFS, ".gohtml", false)
		if err != nil {
			return nil, fmt.Errorf("failed to load views from FS: %w", err)
		}
		for name, content := range subViews {
			if _, exists := views[name]; exists {
				alog.LOGGER(alog.LOGGER_APP).Warn().Msgf("Duplicate view '%s' found; overwriting with latest", name)
			}
			views[name] = content
		}
	}
	if len(views) == 0 {
		return nil, fmt.Errorf("no views found in any FS")
	}

	// Compile each view with layouts (if any).
	for name, viewContent := range views {
		tmpl := template.New(name).Funcs(*t.funcMap)

		// Skip layouts for "sfp-" prefixed views.
		if !strings.HasPrefix(name, "sfp-") {
			for _, layoutContent := range layouts {
				_, err := tmpl.Parse(layoutContent)
				if err != nil {
					return nil, fmt.Errorf("failed to parse layout for view '%s': %w", name, err)
				}
			}
		}

		// Parse the view.
		_, err := tmpl.Parse(viewContent)
		if err != nil {
			return nil, fmt.Errorf("failed to parse view '%s': %w", name, err)
		}

		templates[name] = tmpl
	}

	return templates, nil
}

// loadFilesFromFS walks the FS and loads file contents into a map (name: content).
func (t *FSPagesText) loadFilesFromFS(fsys fs.FS, ext string, optional bool) (map[string]string, error) {
	if fsys == nil && optional {
		return make(map[string]string), nil // Optional and nil: empty
	} else if fsys == nil {
		return nil, fmt.Errorf("required FS is nil")
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
		if optional && (os.IsNotExist(err) || strings.Contains(err.Error(), "file does not exist")) {
			return make(map[string]string), nil // Optional and missing: empty
		}
		return nil, fmt.Errorf("failed to load files from FS: %w", err)
	}
	return files, nil
}

// Render renders a template document using Echo's rendering interface.
func (t *FSPagesText) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if ISON_PAGES_TEXT_LOADBYFILE || ISON_PAGES_TEXT_RELOAD_ON_RENDER {
		return t.renderReload(name, data, w)
	}
	tmpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("template '%s' not found", name)
	}
	var err error
	if strings.HasPrefix(name, "sfp-") {
		err = tmpl.Execute(w, data) // Standalone: execute root
	} else {
		err = tmpl.ExecuteTemplate(w, "base", data) // With layout: execute "base"
	}
	if ISON_PAGES_TEXT_LOG_ERROR && err != nil {
		alog.LOGGER(alog.LOGGER_APP).Error().Err(err).Msg("exe text template")
	}
	return err
}

// renderReload re-parses and renders on-the-fly (for dev/hot-reload).
func (t *FSPagesText) renderReload(name string, data interface{}, w io.Writer) error {
	// Reload layouts (optional).
	layouts, err := t.loadFilesFromFS(t.layoutsFS, ".gohtml", true)
	if err != nil {
		return err
	}

	// Reload views from multiple sources to find the specific view.
	viewContent := ""
	found := false
	for _, vFS := range t.viewsFS {
		views, err := t.loadFilesFromFS(vFS, ".gohtml", false)
		if err != nil {
			return err
		}
		if content, ok := views[name]; ok {
			viewContent = content
			found = true
			break // First found wins; adjust if needed
		}
	}
	if !found {
		return fmt.Errorf("view '%s' not found in any FS", name)
	}

	tmpl := template.New(name).Funcs(*t.funcMap)

	// Skip layouts for "sfp-" prefixed views.
	if !strings.HasPrefix(name, "sfp-") {
		for _, layoutContent := range layouts {
			_, err := tmpl.Parse(layoutContent)
			if err != nil {
				return err
			}
		}
	}

	// Parse the view.
	_, err = tmpl.Parse(viewContent)
	if err != nil {
		return err
	}

	var executeErr error
	if strings.HasPrefix(name, "sfp-") {
		executeErr = tmpl.Execute(w, data) // Standalone
	} else {
		executeErr = tmpl.ExecuteTemplate(w, "base", data) // With layout
	}
	if ISON_PAGES_TEXT_LOG_ERROR && executeErr != nil {
		alog.LOGGER(alog.LOGGER_APP).Error().Err(executeErr).Msg("exe text template reload")
	}
	return executeErr
}

// RenderPage renders a page to a string.
func (t *FSPagesText) RenderPage(name string, data interface{}) (string, error) {
	var buff bytes.Buffer
	err := t.Render(&buff, name, data, nil)
	return buff.String(), err
}
