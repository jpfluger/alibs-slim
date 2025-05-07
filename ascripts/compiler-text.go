package ascripts

import (
	"bytes"
	"fmt"
	"github.com/jpfluger/alibs-slim/atemplates"
	"strings"
	ttemplate "text/template"
)

// CompilerText is a struct that holds a text/template Template.
type CompilerText struct {
	t *ttemplate.Template // The text/template Template instance.
}

// Compile parses the body as a text/template and stores the result in the CompilerText struct.
func (c *CompilerText) Compile(body string) error {
	// Check if the body is empty after trimming whitespace.
	if strings.TrimSpace(body) == "" {
		return fmt.Errorf("body of text is empty")
	}

	// Retrieve the common text template functions.
	fnMap := atemplates.GetTextTemplateFunctions(atemplates.TEMPLATE_FUNCTIONS_COMMON)

	// Parse the body with the template functions to create a new template.
	t, err := ttemplate.New("").Funcs(*fnMap).Parse(body)
	if err != nil {
		return fmt.Errorf("could not parse text template from code body: %v", err)
	}
	c.t = t // Store the parsed template.
	return nil
}

// Run compiles the body and executes the template with the provided parameters.
func (c *CompilerText) Run(body string, params ...interface{}) (interface{}, error) {
	return c.Render(body, params...)
}

// Render compiles the body if not already compiled and renders the template with the provided parameters.
func (c *CompilerText) Render(body string, params ...interface{}) (string, error) {
	// Compile the body if the template is not already compiled.
	if c.t == nil {
		if err := c.Compile(body); err != nil {
			return "", err
		}
	}
	return c.RenderTemplate(params...)
}

// RenderTemplate executes the compiled template with the provided parameters and returns the result as a string.
func (c *CompilerText) RenderTemplate(params ...interface{}) (string, error) {
	// Check if the template is nil before attempting to execute it.
	if c.t == nil {
		return "", fmt.Errorf("template is nil")
	}

	// Use the first parameter if provided, otherwise pass nil to the template execution.
	var param interface{} = map[string]interface{}{}
	if len(params) > 0 && params[0] != nil {
		param = params[0]
	}

	// Execute the template with the parameter and capture the output in a buffer.
	var tmplBuffer bytes.Buffer
	if err := c.t.Execute(&tmplBuffer, param); err != nil {
		return "", fmt.Errorf("could not execute text template: %v", err)
	}
	return tmplBuffer.String(), nil // Return the rendered string.
}
