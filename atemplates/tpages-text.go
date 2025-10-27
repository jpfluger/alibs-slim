package atemplates

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jpfluger/alibs-slim/autils"

	"github.com/jpfluger/alibs-slim/alog"
	"github.com/labstack/echo/v4"
)

// ISON_PAGES_TEXT_LOADBYFILE determines if pages should be loaded from files on each request.
// Deprecated
var ISON_PAGES_TEXT_LOADBYFILE = false

// TPagesText is a custom HTML/template renderer for the Echo framework.
// Deprecated
type TPagesText struct {
	templates   map[string]*template.Template // Map of template names to their compiled templates.
	funcMap     *template.FuncMap             // Function map for template rendering, allows for custom functions.
	viewFileMap map[string]string             // Map of view names to their file paths.
	dirLayouts  string                        // Directory where layout files are located.
}

// NewTPagesText creates a new TPagesText instance, loading templates from the specified directory.
// Deprecated
func NewTPagesText(dir string, funcMap *template.FuncMap) *TPagesText {
	// Use common text template functions if no specific function map is provided.
	if funcMap == nil {
		funcMap = GetTextTemplateFunctions(TEMPLATE_FUNCTIONS_COMMON)
	}

	// List of directories to ensure existence
	if err := autils.CleanDirsWithMkdirOption([]string{
		"pages-text/views/",
		"pages-text/layouts/",
	}, dir, true); err != nil {
		panic(err)
	}

	// Ensure the directory path ends with a slash.
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	// Collect all .gohtml files from the views directory.
	views := []string{}
	fnFiles := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".gohtml") {
			views = append(views, path)
		}
		return nil
	}

	// Walk through the views directory to find view files.
	if err := filepath.Walk(path.Join(dir, "pages-text/views/"), fnFiles); err != nil {
		panic(err) // Panic if walking the directory fails.
	}

	// Collect all layout files from the layouts directory.
	dirLayouts := path.Join(dir, "pages-text/layouts")
	layouts, err := filepath.Glob(dirLayouts + "/*")
	if err != nil {
		panic(err) // Panic if globbing the directory fails.
	} else if layouts == nil {
		layouts = []string{}
	}

	// Initialize the map for view file paths and the map for compiled templates.
	viewFileMap := make(map[string]string)
	templates := make(map[string]*template.Template)

	// Compile templates for each view file found, combining them with layout files.
	for _, view := range views {
		combined := append(layouts, view)
		templates[filepath.Base(view)] = template.Must(template.New(filepath.Base(view)).Funcs(*funcMap).ParseFiles(combined...))

		// If pages should be loaded by file, map the view names to their file paths.
		if ISON_PAGES_TEXT_LOADBYFILE {
			viewFileMap[filepath.Base(view)] = view
		}
	}

	// Return a new TPagesText instance with the compiled templates and file map.
	return &TPagesText{
		templates:   templates,
		funcMap:     funcMap,
		viewFileMap: viewFileMap,
		dirLayouts:  dirLayouts,
	}
}

// Render renders a template document using Echo's rendering interface.
func (t *TPagesText) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if ISON_PAGES_TEXT_LOADBYFILE {
		// Load layouts and view from files if enabled.
		layouts, err := filepath.Glob(t.dirLayouts + "/*")
		if err != nil {
			return fmt.Errorf("could not load layouts")
		}
		viewPath := t.viewFileMap[name]
		if viewPath == "" {
			return fmt.Errorf("view file map not found for '%s'", name)
		}
		files := append(layouts, viewPath)
		myTemplate := template.Must(template.New(name).Funcs(*t.funcMap).ParseFiles(files...))
		err = myTemplate.ExecuteTemplate(w, "base", data)
		if ISON_PAGES_TEXT_LOG_ERROR && err != nil {
			alog.LOGGER(alog.LOGGER_APP).Error().Err(err).Msg("exe text template loadbyfile")
		}
		return err
	}

	// Use precompiled template if loading from files is not enabled.
	tmpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("template '%s' not found", name)
	}
	err := tmpl.ExecuteTemplate(w, "base", data)
	if ISON_PAGES_TEXT_LOG_ERROR && err != nil {
		alog.LOGGER(alog.LOGGER_APP).Error().Err(err).Msg("exe text template")
	}
	return err
}

// RenderPage renders a page to a string.
func (t *TPagesText) RenderPage(name string, data interface{}) (string, error) {
	var buff bytes.Buffer
	if ISON_PAGES_TEXT_LOADBYFILE {
		// Load layouts and view from files if enabled.
		layouts, err := filepath.Glob(t.dirLayouts + "/*")
		if err != nil {
			return "", fmt.Errorf("could not load layouts")
		}
		viewPath := t.viewFileMap[name]
		if viewPath == "" {
			return "", fmt.Errorf("view file map not found for '%s'", name)
		}
		files := append(layouts, viewPath)
		myTemplate := template.Must(template.New(name).Funcs(*t.funcMap).ParseFiles(files...))
		err = myTemplate.Execute(&buff, data)
		if err != nil {
			if ISON_PAGES_TEXT_LOG_ERROR {
				alog.LOGGER(alog.LOGGER_APP).Error().Err(err).Msg("exe render text template loadbyfile")
			}
			return "", err
		}
		return buff.String(), nil
	}

	// Use precompiled template if loading from files is not enabled.
	tmpl, ok := t.templates[name]
	if !ok {
		return "", fmt.Errorf("template '%s' not found", name)
	}
	err := tmpl.Execute(&buff, data)
	if err != nil {
		if ISON_PAGES_TEXT_LOG_ERROR {
			alog.LOGGER(alog.LOGGER_APP).Error().Err(err).Msg("exe render text template")
		}
		return "", fmt.Errorf("failed execute of TPageText '%s'; %v", name, err)
	}
	return buff.String(), nil
}
