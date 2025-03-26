package atemplates

import (
	"bytes"
	"fmt"
	"github.com/jpfluger/alibs-slim/autils"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/jpfluger/alibs-slim/alog"
)

// ISON_PAGES_HTML_LOADBYFILE indicates whether to load HTML pages from files on each request.
var ISON_PAGES_HTML_LOADBYFILE = false

// ISON_PAGES_HTML_LOG_ERROR enables logging of errors during HTML template execution.
var ISON_PAGES_HTML_LOG_ERROR = false

// TPagesHTML is a custom HTML/template renderer for the Echo framework.
type TPagesHTML struct {
	templates   map[string]*template.Template // Map of template names to their compiled templates.
	funcMap     *template.FuncMap             // Function map for template rendering, allows for custom functions.
	viewFileMap map[string]string             // Map of view names to their file paths.
	dirLayouts  string                        // Directory where layout files are located.
}

// NewTPagesHTML creates a new TPagesHTML instance, loading templates from the specified directory.
func NewTPagesHTML(dir string, funcMap *template.FuncMap) *TPagesHTML {
	// Use common HTML template functions if no specific function map is provided.
	if funcMap == nil {
		funcMap = GetHTMLTemplateFunctions(TEMPLATE_FUNCTIONS_COMMON)
	}

	// List of directories to ensure existence
	if err := autils.CleanDirsWithMkdirOption([]string{
		"pages-html/views/",
		"pages-html/layouts/",
	}, dir, true); err != nil {
		panic(err)
	}

	// Ensure the directory path ends with a slash.
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	// Initialize the map for view file paths and the map for compiled templates.
	viewFileMap := make(map[string]string)
	templates := make(map[string]*template.Template)

	// Walk through the views directory to find view files.
	views := []string{}
	fnFiles := func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".gohtml") {
			views = append(views, filePath)
		}
		return nil
	}
	if err := filepath.Walk(path.Join(dir, "pages-html/views/"), fnFiles); err != nil {
		panic(err) // Panic if walking the directory fails.
	}

	// Collect all layout files from the layouts directory.
	dirLayouts := path.Join(dir, "pages-html/layouts")
	layouts, err := filepath.Glob(dirLayouts + "/*")
	if err != nil {
		panic(err) // Panic if globbing the directory fails.
	} else if layouts == nil {
		layouts = []string{}
	}

	// Compile templates for each view file found, combining them with layout files.
	for _, viewPath := range views {
		name := filepath.Base(viewPath)
		var files []string
		if strings.HasPrefix(name, "sfp-") {
			files = []string{viewPath}
		} else {
			files = append(layouts, viewPath)
		}
		templates[name] = template.Must(template.New(name).Funcs(*funcMap).ParseFiles(files...))

		// If pages should be loaded by file, map the view names to their file paths.
		if ISON_PAGES_HTML_LOADBYFILE {
			viewFileMap[name] = viewPath
		}
	}

	// Return a new TPagesHTML instance with the compiled templates and file map.
	return &TPagesHTML{
		templates:   templates,
		funcMap:     funcMap,
		viewFileMap: viewFileMap,
		dirLayouts:  dirLayouts,
	}
}

// Render renders a template document using Echo's rendering interface.
func (t *TPagesHTML) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if ISON_PAGES_HTML_LOADBYFILE {
		// Load layouts and view from files if enabled.
		viewPath, ok := t.viewFileMap[name]
		if !ok {
			return fmt.Errorf("view file map not found for '%s'", name)
		}

		var files []string
		if strings.HasPrefix(filepath.Base(viewPath), "sfp-") {
			files = []string{viewPath}
		} else {
			layouts, err := filepath.Glob(t.dirLayouts + "/*")
			if err != nil {
				return fmt.Errorf("could not load layouts")
			}
			files = append(layouts, viewPath)
		}

		myTemplate := template.Must(template.New(name).Funcs(*t.funcMap).ParseFiles(files...))
		err := myTemplate.ExecuteTemplate(w, "base", data)
		if ISON_PAGES_HTML_LOG_ERROR && err != nil {
			alog.LOGGER(alog.LOGGER_APP).Error().Err(err).Msg("exe html template loadbyfile")
		}
		return err
	}

	// Use precompiled template if loading from files is not enabled.
	tmpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("template '%s' not found", name)
	}
	err := tmpl.ExecuteTemplate(w, "base", data)
	if ISON_PAGES_HTML_LOG_ERROR && err != nil {
		alog.LOGGER(alog.LOGGER_APP).Error().Err(err).Msg("exe html template")
	}
	return err
}

// RenderPage renders a page to a string.
func (t *TPagesHTML) RenderPage(name string, data interface{}) (string, error) {
	var buff bytes.Buffer
	if ISON_PAGES_HTML_LOADBYFILE {
		// Load layouts and view from files if enabled.
		viewPath, ok := t.viewFileMap[name]
		if !ok {
			return "", fmt.Errorf("view file map not found for '%s'", name)
		}

		var files []string
		if strings.HasPrefix(filepath.Base(viewPath), "sfp-") {
			files = []string{viewPath}
		} else {
			layouts, err := filepath.Glob(t.dirLayouts + "/*")
			if err != nil {
				return "", fmt.Errorf("could not load layouts")
			}
			files = append(layouts, viewPath)
		}

		myTemplate := template.Must(template.New(name).Funcs(*t.funcMap).ParseFiles(files...))
		err := myTemplate.Execute(&buff, data)
		if err != nil {
			if ISON_PAGES_HTML_LOG_ERROR {
				alog.LOGGER(alog.LOGGER_APP).Error().Err(err).Msg("exe render html template loadbyfile")
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
		if ISON_PAGES_HTML_LOG_ERROR {
			alog.LOGGER(alog.LOGGER_APP).Error().Err(err).Msg("exe render html template")
		}
		return "", fmt.Errorf("failed execute of TPageText '%s'; %v", name, err)
	}
	return buff.String(), nil
}
