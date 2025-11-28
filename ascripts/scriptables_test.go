package ascripts

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/jpfluger/alibs-slim/autils"
	"github.com/stretchr/testify/assert"
)

// TestScriptables_A tests the rendering of Scriptable objects with provided data.
func TestScriptables_A(t *testing.T) {
	// Initialize a slice of Scriptable objects with different types and bodies.
	ss := Scriptables{
		{Key: "text.1.a", Type: SCRIPTTYPE_TEXT, Body: "This is a {{.ReplaceTEXT}}!"},
		{Key: "html.1.a", Type: SCRIPTTYPE_HTML, Body: "<html>something here <div style=\"{{.AttribHTML}}\">MyStuff</div></html>"},
		{Key: "css.1.a", Type: SCRIPTTYPE_CSS, Body: "body .mydiv {background-color: \"{{.ColorCSS}}\"}"},
	}

	// Data to be used for rendering the Scriptable objects.
	data := map[string]interface{}{
		"ReplaceTEXT": "sentence",
		"AttribHTML":  "background-color:#fff",
		"ColorCSS":    "#fff",
	}

	// Iterate over the Scriptable objects and test the rendering.
	for ii, s := range ss {
		content, err := s.Render(data)
		assert.NoError(t, err, "Rendering failed for Scriptable at index %d", ii)
		expectedContents := []string{
			"This is a sentence!",
			"<html>something here <div style=\"background-color:#fff\">MyStuff</div></html>",
			"body .mydiv {background-color: \"#fff\"}",
		}
		assert.Equal(t, expectedContents[ii], content, "Rendered content does not match expected for Scriptable at index %d", ii)
	}

	// Convert the slice of Scriptable objects to a map and test the rendering.
	ssMap, err := ss.ToMap()
	assert.NoError(t, err, "Conversion to map failed")
	assert.Equal(t, 3, len(ssMap), "The length of the map does not match the expected value")

	// Test the rendering of each Scriptable object in the map.
	for key, val := range ssMap {
		content, err := val.Render(data)
		assert.NoError(t, err, "Rendering failed for Scriptable with key %s", key)
		expectedContents := map[string]string{
			"text.1.a": "This is a sentence!",
			"html.1.a": "<html>something here <div style=\"background-color:#fff\">MyStuff</div></html>",
			"css.1.a":  "body .mydiv {background-color: \"#fff\"}",
		}
		assert.Equal(t, expectedContents[key.String()], content, "Rendered content does not match expected for Scriptable with key %s", key)
	}
}

// TestNewScriptableFromPath_A tests the creation of Scriptable objects from file paths.
func TestNewScriptableFromPath_A(t *testing.T) {
	// Array of file paths to be used for creating Scriptable objects.
	arrFiles := []string{"test-data/a.txt", "test-data/b.js", "test-data/c.html", "test-data/d.md"}

	// Get the current working directory.
	wd, err := os.Getwd()
	assert.NoError(t, err, "Getting current working directory failed")

	// Iterate over the file paths and test the creation of Scriptable objects.
	for _, arrFile := range arrFiles {
		filePath := path.Join(wd, arrFile)
		scriptable, err := NewScriptableFromPath(filePath)
		assert.NoError(t, err, "Creating Scriptable from path failed for file %s", arrFile)

		full, name, ext := autils.FileNameParts(arrFile)
		assert.Equal(t, path.Base(full), scriptable.Key.String(), "Key does not match the base of the file path for file %s", arrFile)
		assert.Equal(t, name, scriptable.Body, "Body does not match the name part of the file path for file %s", arrFile)
		assert.Equal(t, ExtToScriptType(ext).String(), scriptable.Type.String(), "Type does not match the extension for file %s", arrFile)
	}

	// Test the creation of default ScriptTypeFilters.
	sfts := GetScriptTypeFilterDefaults(ScriptTypes{SCRIPTTYPE_HTML, SCRIPTTYPE_MARKDOWN})
	assert.Equal(t, 2, len(sfts), "The length of ScriptTypeFilters does not match the expected value")
	assert.Equal(t, SCRIPTTYPE_HTML.String(), sfts[0].Type.String(), "The first ScriptTypeFilter does not match SCRIPTTYPE_HTML")
	assert.Equal(t, SCRIPTTYPE_MARKDOWN.String(), sfts[1].Type.String(), "The second ScriptTypeFilter does not match SCRIPTTYPE_MARKDOWN")
}

// TestMarkdownSample tests the creation of the test-data/samples.md file.
func TestMarkdownSample(t *testing.T) {
	// Get the current working directory.
	wd, err := os.Getwd()
	assert.NoError(t, err, "Getting current working directory failed")

	filePath := path.Join(wd, "test-data/samples.md")
	scriptable, err := NewScriptableFromPath(filePath)
	assert.NoError(t, err, "Creating Scriptable from path failed for file %s", filePath)

	output, err := scriptable.Render(nil)
	assert.NoError(t, err, "Rendering failed for file %s", filePath)

	filePathControl := path.Join(wd, "test-data/sample-test-output.html")
	control, err := os.ReadFile(filePathControl)
	assert.NoError(t, err, "Reading control file failed for file %s", filePathControl)
	assert.Equal(t, strings.TrimSpace(string(control)), strings.TrimSpace(output), "Rendered content does not match expected for file %s", filePathControl)
}

// TestMarkdownHtmlSample tests the creation of the test-data/samples.mdh file.
func TestMarkdownHtmlSample(t *testing.T) {
	// Get the current working directory.
	wd, err := os.Getwd()
	assert.NoError(t, err, "Getting current working directory failed")

	filePath := path.Join(wd, "test-data/samples.mdh")
	scriptable, err := NewScriptableFromPath(filePath)
	assert.NoError(t, err, "Creating Scriptable from path failed for file %s", filePath)

	output, err := scriptable.Render(nil)
	assert.NoError(t, err, "Rendering failed for file %s", filePath)

	filePathControl := path.Join(wd, "test-data/sample-test-html-output.html")
	control, err := os.ReadFile(filePathControl)
	assert.NoError(t, err, "Reading control file failed for file %s", filePathControl)
	assert.Equal(t, strings.TrimSpace(string(control)), strings.TrimSpace(output), "Rendered content does not match expected for file %s", filePathControl)
}
