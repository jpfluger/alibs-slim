package autils

import (
	"golang.org/x/text/language"
	"strings"
)

// LanguageType represents a language code following the IETF BCP 47 standard.
type LanguageType string

// IsEmpty checks if the LanguageType is empty after trimming whitespace.
func (lant LanguageType) IsEmpty() bool {
	return strings.TrimSpace(string(lant)) == ""
}

// TrimSpace returns a new LanguageType with leading and trailing whitespace removed.
func (lant LanguageType) TrimSpace() LanguageType {
	return LanguageType(strings.TrimSpace(string(lant)))
}

// String returns the string representation of the LanguageType.
func (lant LanguageType) String() string {
	return string(lant)
}

// GetParts splits the LanguageType into its language and locale components.
func (lant LanguageType) GetParts() (lang string, locale string) {
	parts := strings.SplitN(lant.String(), "-", 2)
	switch len(parts) {
	case 0:
		return "", ""
	case 1:
		return parts[0], ""
	default:
		return parts[0], parts[1]
	}
}

// GetLanguage extracts the language part from the LanguageType.
func (lant LanguageType) GetLanguage() string {
	lang, _ := lant.GetParts()
	return lang
}

// GetLocale extracts the locale part from the LanguageType.
func (lant LanguageType) GetLocale() string {
	_, locale := lant.GetParts()
	return locale
}

// GetHTTPAcceptedLanguage parses the Accept-Language header to extract the language preference.
func GetHTTPAcceptedLanguage(target string) LanguageType {
	return GetHTTPAcceptedLanguageWithDefault(target, "")
}

// GetHTTPAcceptedLanguageWithDefault parses the Accept-Language header and returns the language preference,
// falling back to a default if necessary.
func GetHTTPAcceptedLanguageWithDefault(target string, altDefault LanguageType) LanguageType {
	if target != "" {
		// Parse the Accept-Language header using the language package.
		tags, _, err := language.ParseAcceptLanguage(target)
		if err == nil && len(tags) > 0 {
			// If parsing is successful, use the most preferred language tag.
			target = tags[0].String()
		}
	}
	if target == "" {
		// If no language is determined, use the provided default or "en-US" if no default is provided.
		if altDefault.IsEmpty() {
			return "en-US"
		}
		return altDefault
	}
	return LanguageType(target)
}
