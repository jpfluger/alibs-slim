package ascripts

import (
	"bytes"
	"fmt"
	"github.com/jpfluger/alibs-slim/atemplates"
	"strings"
	ttemplate "text/template"
)

// CompilerCSS is a struct that holds a text/template Template for CSS.
type CompilerCSS struct {
	t *ttemplate.Template // The text/template Template instance for CSS.
}

// Run compiles the CSS body and executes the template with the provided parameters.
func (c *CompilerCSS) Run(body string, params ...interface{}) (interface{}, error) {
	return c.Render(body, params...)
}

// Render compiles the CSS body if not already compiled and renders the template with the provided parameters.
func (c *CompilerCSS) Render(body string, params ...interface{}) (string, error) {
	// Check if the body is empty after trimming whitespace.
	if strings.TrimSpace(body) == "" {
		return "", fmt.Errorf("body of CSS is empty")
	}

	// Compile the body as a text/template if it has not been compiled yet.
	if c.t == nil {
		// Retrieve the common text template functions.
		fnMap := atemplates.GetTextTemplateFunctions(atemplates.TEMPLATE_FUNCTIONS_COMMON)

		// Parse the body with the template functions to create a new template.
		t, err := ttemplate.New("").Funcs(*fnMap).Parse(body)
		if err != nil {
			return "", fmt.Errorf("could not parse CSS template from code body: %v", err)
		}
		c.t = t // Store the parsed template.
	}

	// Use the first parameter if provided, otherwise pass nil to the template execution.
	var param interface{}
	if len(params) > 0 {
		param = params[0]
	}

	// Execute the template with the parameter and capture the output in a buffer.
	var tmplBuffer bytes.Buffer
	if err := c.t.Execute(&tmplBuffer, param); err != nil {
		return "", fmt.Errorf("could not execute CSS template: %v", err)
	}
	return tmplBuffer.String(), nil // Return the rendered CSS.
}
