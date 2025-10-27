package ascripts

import (
	"bytes"
	"fmt"
	htemplate "html/template"
	"strings"

	"github.com/jpfluger/alibs-slim/atemplates"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// CompilerMarkdown is a struct that holds an HTML template.
type CompilerMarkdown struct {
	t *htemplate.Template // The HTML template instance.
}

// Run compiles the markdown body and executes the template with the provided parameters.
func (c *CompilerMarkdown) Run(body string, params ...interface{}) (returnedItem interface{}, err error) {
	return c.Render(body, params...)
}

// Render converts markdown to HTML, parses it as an HTML template, and renders it with the provided parameters.
func (c *CompilerMarkdown) Render(body string, params ...interface{}) (string, error) {
	// Check if the body is empty after trimming whitespace.
	if strings.TrimSpace(body) == "" {
		return "", fmt.Errorf("body of markdown is empty")
	}

	// Convert markdown to HTML and parse it as an HTML template if not already done.
	if c.t == nil {
		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM), // Use GitHub Flavored Markdown extensions.
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(), // Automatically generate heading IDs.
			),
			goldmark.WithRendererOptions(
				html.WithHardWraps(), // Use hard wraps.
				html.WithUnsafe(),    // Added for unsafe HTML: Enables raw HTML rendering.
				// html.WithXHTML(), // Uncomment to enable XHTML output.
			),
		)

		var mdBuffer bytes.Buffer
		if err := md.Convert([]byte(body), &mdBuffer); err != nil {
			return "", fmt.Errorf("failed markdown conversion: %v", err)
		}

		content := mdBuffer.String()

		// Retrieve the common HTML template functions.
		fnMap := atemplates.GetHTMLTemplateFunctions(atemplates.TEMPLATE_FUNCTIONS_COMMON)
		t, err := htemplate.New("").Funcs(*fnMap).Parse(content)
		if err != nil {
			return "", fmt.Errorf("could not parse HTML template from markdown content: %v", err)
		}
		c.t = t // Store the parsed template.
	}

	// Use the first parameter if provided, otherwise pass nil to the template execution.
	var param interface{} = map[string]interface{}{}
	if len(params) > 0 && params[0] != nil {
		param = params[0]
	}

	// Execute the template with the parameter and capture the output in a buffer.
	var tmplBuffer bytes.Buffer
	if err := c.t.Execute(&tmplBuffer, param); err != nil {
		return "", fmt.Errorf("could not execute HTML template: %v", err)
	}
	return tmplBuffer.String(), nil // Return the rendered string.
}

// MarkdownToHTML converts a markdown string to HTML.
func MarkdownToHTML(source string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(source), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil // Return the converted HTML string.
}
