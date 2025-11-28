package atemplates

import (
	"fmt"
	htemplate "html/template"
	"io/fs"
	"sync"
	ttemplate "text/template"

	"github.com/jpfluger/alibs-slim/alog"
)

// FSManage is a struct that holds instances of FS-based template managers for different types of templates.
type FSManage struct {
	PagesText    *FSPagesText    // Manager for text page templates.
	PagesHTML    *FSPagesHTML    // Manager for HTML page templates.
	SnippetsHTML *FSSnippetsHTML // Manager for HTML snippet templates.
	SnippetsText *FSSnippetsText // Manager for text snippet templates.
	mu           sync.RWMutex    // Mutex for concurrent read/write access to the template managers.
}

type FSSourceType string

const (
	FSSOURCETYPE_PAGES_HTML_LAYOUT FSSourceType = "pages-html-layout"
	FSSOURCETYPE_PAGES_HTML_VIEWS  FSSourceType = "pages-html-views"
	FSSOURCETYPE_PAGES_TEXT_LAYOUT FSSourceType = "pages-text-layout"
	FSSOURCETYPE_PAGES_TEXT_VIEWS  FSSourceType = "pages-text-views"
	FSSOURCETYPE_SNIPPETS_HTML     FSSourceType = "snippets-html"
	FSSOURCETYPE_SNIPPETS_TEXT     FSSourceType = "snippets-text"
)

// NewFSManage creates and initializes a new FSManage with the provided FS sources and funcMaps.
// sources is a map like {"pages-html-layouts": fs.FS, "pages-html-views": fs.FS, "snippets-html": fs.FS, ...}.
func NewFSManage(sources map[FSSourceType]fs.FS, funcMapHTML *htemplate.FuncMap, funcMapText *ttemplate.FuncMap) (*FSManage, error) {
	fsm := &FSManage{}

	// PagesHTML (HTML pages).
	pagesHTMLLayoutsFS := sources[FSSOURCETYPE_PAGES_HTML_LAYOUT]
	pagesHTMLViewsFS := sources[FSSOURCETYPE_PAGES_HTML_VIEWS]
	if pagesHTMLViewsFS != nil { // Views required.
		var err error
		fsm.PagesHTML, err = NewFSPagesHTML(pagesHTMLLayoutsFS, []fs.FS{pagesHTMLViewsFS}, funcMapHTML)
		if err != nil {
			return nil, fmt.Errorf("failed to init PagesHTML: %w", err)
		}
	}

	// SnippetsHTML.
	snippetsHTMLFS := sources[FSSOURCETYPE_SNIPPETS_HTML]
	if snippetsHTMLFS != nil {
		var err error
		fsm.SnippetsHTML, err = NewFSSnippetsHTML([]fs.FS{snippetsHTMLFS}, funcMapHTML)
		if err != nil {
			return nil, fmt.Errorf("failed to init SnippetsHTML: %w", err)
		}
	}

	// PagesText (text pages).
	pagesTextLayoutsFS := sources[FSSOURCETYPE_PAGES_TEXT_LAYOUT]
	pagesTextViewsFS := sources[FSSOURCETYPE_PAGES_TEXT_VIEWS]
	if pagesTextViewsFS != nil { // Views required for PagesText.
		var err error
		fsm.PagesText, err = NewFSPagesText(pagesTextLayoutsFS, []fs.FS{pagesTextViewsFS}, funcMapText)
		if err != nil {
			return nil, fmt.Errorf("failed to init PagesText: %w", err)
		}
	}

	// SnippetsText.
	snippetsTextFS := sources[FSSOURCETYPE_SNIPPETS_TEXT]
	if snippetsTextFS != nil {
		var err error
		fsm.SnippetsText, err = NewFSSnippetsText([]fs.FS{snippetsTextFS}, funcMapText)
		if err != nil {
			return nil, fmt.Errorf("failed to init SnippetsText: %w", err)
		}
	}

	return fsm, nil
}

var (
	fsTmanage *FSManage    // Singleton instance of FSManage.
	fsMu      sync.RWMutex // Mutex for concurrent access to the singleton instance.
)

// FSTEMPLATES provides safe access to the singleton instance of FSManage.
func FSTEMPLATES() *FSManage {
	fsMu.RLock()
	defer fsMu.RUnlock()
	return fsTmanage
}

// SetFSTemplates safely sets the singleton instance of FSManage.
func SetFSTemplates(tm *FSManage) {
	fsMu.Lock()
	defer fsMu.Unlock()
	fsTmanage = tm
}

// RenderPageText renders a text page template with the provided data.
func (tm *FSManage) RenderPageText(name string, data interface{}) (string, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.PagesText.RenderPage(name, data)
}

// RenderPageHTML renders an HTML page template with the provided data.
func (tm *FSManage) RenderPageHTML(name string, data interface{}) (string, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.PagesHTML.RenderPage(name, data)
}

// RenderSnippetsHTML renders an HTML snippet template with the provided data.
func (tm *FSManage) RenderSnippetsHTML(name string, data interface{}) (string, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.SnippetsHTML.RenderSnippet(name, data)
}

// RenderSnippetsText renders a text snippet template with the provided data.
func (tm *FSManage) RenderSnippetsText(name string, data interface{}) (string, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.SnippetsText.RenderSnippet(name, data)
}

// MustFSSnippetRenderHTML renders an HTML snippet template and returns an empty string on error.
func MustFSSnippetRenderHTML(name string, data interface{}) string {
	s, err := FSTEMPLATES().RenderSnippetsHTML(name, data)
	if err != nil {
		alog.LOGGER(alog.LOGGER_APP).Err(err).Msg("MustFSSnippetRenderHTML")
		return ""
	}
	return s
}

// MustFSSnippetRenderText renders a text snippet template and returns an empty string on error.
func MustFSSnippetRenderText(name string, data interface{}) string {
	s, err := FSTEMPLATES().RenderSnippetsText(name, data)
	if err != nil {
		alog.LOGGER(alog.LOGGER_APP).Err(err).Msg("MustFSSnippetRenderText")
		return ""
	}
	return s
}

// GetFSSnippetHTMLIsLoaded checks if the specified HTML snippet is loaded.
// It uses a read lock for thread-safe access and returns an error if SnippetsHTML is not initialized or the snippet is not found.
func (tm *FSManage) GetFSSnippetHTMLIsLoaded(name string) error {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if tm.SnippetsHTML == nil {
		return fmt.Errorf("SnippetsHTML not initialized")
	}
	return tm.SnippetsHTML.IsLoaded(name)
}

// GetFSSnippetTextIsLoaded checks if the specified text snippet is loaded.
// It uses a read lock for thread-safe access and returns an error if SnippetsText is not initialized or the snippet is not found.
func (tm *FSManage) GetFSSnippetTextIsLoaded(name string) error {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if tm.SnippetsText == nil {
		return fmt.Errorf("SnippetsText not initialized")
	}
	return tm.SnippetsText.IsLoaded(name)
}
