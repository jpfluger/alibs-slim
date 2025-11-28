package ascripts

import (
	"reflect"

	"github.com/jpfluger/alibs-slim/areflect"
)

// TYPEMANAGER_SCRIPT_COMPILER is a constant key used for registering script compilers.
const (
	TYPEMANAGER_SCRIPT_COMPILER = "script-compiler"
)

// init registers the script compiler types with the type manager upon package initialization.
func init() {
	// Ignore the error returned by Register as it's not being handled.
	_ = areflect.TypeManager().Register(TYPEMANAGER_SCRIPT_COMPILER, "autils/scripts", returnScriptCompiler)
}

// returnScriptCompiler returns the reflect.Type corresponding to the typeName of the script compiler.
func returnScriptCompiler(typeName string) (reflect.Type, error) {
	var rtype reflect.Type
	switch typeName {
	case SCRIPTTYPE_GO.String():
		rtype = reflect.TypeOf(CompilerGolangYaegi{}) // Example compiler for Go.
	case SCRIPTTYPE_HTML.String():
		rtype = reflect.TypeOf(CompilerHTML{}) // Example compiler for HTML.
	case SCRIPTTYPE_MARKDOWN.String():
		rtype = reflect.TypeOf(CompilerMarkdown{}) // Example compiler for Markdown.
	case SCRIPTTYPE_CSS.String():
		rtype = reflect.TypeOf(CompilerCSS{}) // Example compiler for CSS.
	case SCRIPTTYPE_TEXT.String():
		rtype = reflect.TypeOf(CompilerText{}) // Example compiler for plain text.
	case SCRIPTTYPE_MARKDOWN_HTML.String():
		rtype = reflect.TypeOf(CompilerMarkdownHtml{}) // Example compiler for Markdown Html (potentially unsafe).
	}
	return rtype, nil
}

// ICompiler is an interface that defines methods for compiling and rendering scripts.
type ICompiler interface {
	// Run compiles the script body with the given parameters and returns the result or an error.
	Run(body string, params ...interface{}) (returnedItem interface{}, err error)

	// Render processes the script body with the given parameters and returns the result as a string or an error.
	Render(body string, params ...interface{}) (string, error)
}

// IRenderer is an interface used by certain Go compilers to obtain the rendered string result.
type IRenderer interface {
	// GetRendered retrieves the rendered string result.
	GetRendered() string
}
