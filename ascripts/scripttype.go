package ascripts

import (
	"strings" // Importing strings package for string manipulation

	"github.com/jpfluger/alibs-slim/autils" // Importing custom utility package
)

// ScriptType defines a type for different script types as a string
type ScriptType string

// Constants for different script types
const (
	SCRIPTTYPE_GO            ScriptType = "go"
	SCRIPTTYPE_HTML          ScriptType = "html"
	SCRIPTTYPE_MARKDOWN      ScriptType = "markdown"
	SCRIPTTYPE_MARKDOWN_HTML ScriptType = "markdown-html"
	SCRIPTTYPE_CSS           ScriptType = "css"
	SCRIPTTYPE_TEXT          ScriptType = "text"
	SCRIPTTYPE_JS            ScriptType = "js"  // JavaScript type, implementation left to higher-order apps
	SCRIPTTYPE_SQL           ScriptType = "sql" // SQL type, similar to JS, implemented by higher-order apps
)

// IsEmpty checks if the ScriptType is empty after trimming spaces
func (sType ScriptType) IsEmpty() bool {
	return strings.TrimSpace(string(sType)) == ""
}

// TrimSpace trims spaces from ScriptType and returns a new ScriptType
func (sType ScriptType) TrimSpace() ScriptType {
	return ScriptType(strings.TrimSpace(string(sType)))
}

// String returns the ScriptType as a trimmed string
func (sType ScriptType) String() string {
	return strings.TrimSpace(string(sType))
}

// TrimSpaceToLower trims spaces from ScriptType and converts it to lower case
func (sType ScriptType) TrimSpaceToLower() ScriptType {
	return ScriptType(autils.ToStringTrimLower(string(sType)))
}

// GetExt returns the file extension associated with the ScriptType
func (sType ScriptType) GetExt() string {
	switch sType {
	case SCRIPTTYPE_GO:
		return ".go"
	case SCRIPTTYPE_HTML:
		return ".html"
	case SCRIPTTYPE_MARKDOWN:
		return ".md"
	case SCRIPTTYPE_MARKDOWN_HTML:
		return ".md"
	case SCRIPTTYPE_CSS:
		return ".css"
	case SCRIPTTYPE_JS:
		return ".js"
	case SCRIPTTYPE_TEXT:
		return ".txt" // Can also represent templates
	case SCRIPTTYPE_SQL:
		return ".sql"
	default:
		return "." + sType.String() // Default case returns the ScriptType as an extension
	}
}

// FilePathToScriptType converts a file path to a ScriptType based on its extension
func FilePathToScriptType(filePath string) ScriptType {
	return ExtToScriptType(autils.GetFileNamePartExtNoDotPrefixToLower(filePath))
}

// ExtToScriptType converts a file extension to a ScriptType
func ExtToScriptType(ext string) ScriptType {
	ext = autils.StripExtensionPrefix(autils.ToStringTrimLower(ext))
	switch ext {
	case "go":
		return SCRIPTTYPE_GO
	case "html", "htm", "gohtml":
		return SCRIPTTYPE_HTML
	case "md", "markdown":
		return SCRIPTTYPE_MARKDOWN
	case "mdh", "markdownh":
		return SCRIPTTYPE_MARKDOWN_HTML
	case "css":
		return SCRIPTTYPE_CSS
	case "js":
		return SCRIPTTYPE_JS
	case "txt", "text", "tmpl", "template":
		return SCRIPTTYPE_TEXT
	case "sql":
		return SCRIPTTYPE_SQL
	default:
		return ScriptType(ext) // Default case returns the extension as a ScriptType
	}
}

// ScriptTypes defines a slice of ScriptType
type ScriptTypes []ScriptType

// HasMatch checks if the target ScriptType is present in the slice
func (sts ScriptTypes) HasMatch(target ScriptType) bool {
	for _, st := range sts {
		if st == target {
			return true
		}
	}
	return false // Returns false if no match is found
}
