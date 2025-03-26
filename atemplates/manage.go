package atemplates

import (
	"github.com/jpfluger/alibs-slim/alog"
	"sync"
)

// TemplateFunctions is an enumeration of available template function sets.
type TemplateFunctions uint32

// Enumeration values for TemplateFunctions.
const (
	TEMPLATE_FUNCTIONS_COMMON TemplateFunctions = iota
)

// Manage is a struct that holds instances of template managers for different types of templates.
type Manage struct {
	PagesText    *TPagesText    // Manager for text page templates.
	PagesHTML    *TPagesHTML    // Manager for HTML page templates.
	SnippetsHTML *TSnippetsHTML // Manager for HTML snippet templates.
	SnippetsText *TSnippetsText // Manager for text snippet templates.
	mu           sync.RWMutex   // Mutex for concurrent read/write access to the template managers.
}

var (
	tmanage *Manage      // Singleton instance of Manage.
	mu      sync.RWMutex // Mutex for concurrent access to the singleton instance.
)

// TEMPLATES provides safe access to the singleton instance of Manage.
func TEMPLATES() *Manage {
	mu.RLock()
	defer mu.RUnlock()
	return tmanage
}

// SetTemplates safely sets the singleton instance of Manage.
func SetTemplates(tm *Manage) {
	mu.Lock()
	defer mu.Unlock()
	tmanage = tm
}

// RenderPageText renders a text page template with the provided data.
func (tm *Manage) RenderPageText(name string, data interface{}) (string, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.PagesText.RenderPage(name, data)
}

// RenderPageHTML renders an HTML page template with the provided data.
func (tm *Manage) RenderPageHTML(name string, data interface{}) (string, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.PagesHTML.RenderPage(name, data)
}

// RenderSnippetsHTML renders an HTML snippet template with the provided data.
func (tm *Manage) RenderSnippetsHTML(name string, data interface{}) (string, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.SnippetsHTML.RenderSnippet(name, data)
}

// RenderSnippetsText renders a text snippet template with the provided data.
func (tm *Manage) RenderSnippetsText(name string, data interface{}) (string, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.SnippetsText.RenderSnippet(name, data)
}

// MustSnippetRenderHTML renders an HTML snippet template and returns an empty string on error.
func MustSnippetRenderHTML(name string, data interface{}) string {
	s, err := TEMPLATES().RenderSnippetsHTML(name, data)
	if err != nil {
		alog.LOGGER(alog.LOGGER_APP).Err(err).Msg("MustSnippetRenderHTML")
		return ""
	}
	return s
}

// MustSnippetRenderText renders a text snippet template and returns an empty string on error.
func MustSnippetRenderText(name string, data interface{}) string {
	s, err := TEMPLATES().RenderSnippetsText(name, data)
	if err != nil {
		alog.LOGGER(alog.LOGGER_APP).Err(err).Msg("MustSnippetRenderText")
		return ""
	}
	return s
}
