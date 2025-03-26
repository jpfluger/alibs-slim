package autils

import (
	"testing"
)

func TestGetSystemLocale(t *testing.T) {
	lang, locale := GetSystemLanguage()

	// Check that the language part is not empty.
	if lang == "" {
		t.Errorf("GetSystemLanguage() returned an empty language")
	}

	// Check that the locale part is not empty.
	if locale == "" {
		t.Errorf("GetSystemLanguage() returned an empty locale")
	}
}

func TestGetSystemLanguageType(t *testing.T) {
	lt := GetSystemLanguageType()
	lan, loc := lt.GetParts()

	// Check that the language part is not empty.
	if lan == "" {
		t.Errorf("GetSystemLanguageType() returned an empty language")
	}

	// Check that the locale part is not empty.
	if loc == "" {
		t.Errorf("GetSystemLanguageType() returned an empty locale")
	}
}

// TestGetSystemLanguage tests the GetSystemLanguage function.
func TestGetSystemLanguage(t *testing.T) {
	// Since GetSystemLanguage relies on the underlying OS, we need to mock the OS-specific behavior.
	// For the purpose of this test, we'll assume a default return value.
	expectedLang := "en"
	expectedLoc := "US"

	lang, loc := GetSystemLanguage()
	if lang != expectedLang || loc != expectedLoc {
		t.Errorf("GetSystemLanguage() = %v, %v, want %v, %v", lang, loc, expectedLang, expectedLoc)
	}

	// Additional tests could be added here to mock and test OS-specific commands.
	// However, this would require a more complex setup with interfaces and dependency injection.
}
