package autils

import (
	"testing"
)

// TestLanguageType_LanguageType tests the LanguageType struct methods for parsing and retrieving language and locale parts.
// It verifies that the GetParts, GetLanguage, and GetLocale methods return the correct values for standard and extended language tags.
func TestLanguageType_LanguageType(t *testing.T) {
	lt := LanguageType("en-US")
	lan, loc := lt.GetParts()
	if lan != "en" {
		t.Errorf("lt.GetParts() lan = %v, want %v", lan, "en")
	}
	if loc != "US" {
		t.Errorf("lt.GetParts() loc = %v, want %v", loc, "US")
	}
	if got := lt.GetLanguage(); got != "en" {
		t.Errorf("lt.GetLanguage() = %v, want %v", got, "en")
	}
	if got := lt.GetLocale(); got != "US" {
		t.Errorf("lt.GetLocale() = %v, want %v", got, "US")
	}

	lt = LanguageType("en-US-really-long")
	lan, loc = lt.GetParts()
	if lan != "en" {
		t.Errorf("lt.GetParts() lan = %v, want %v", lan, "en")
	}
	if loc != "US-really-long" {
		t.Errorf("lt.GetParts() loc = %v, want %v", loc, "US-really-long")
	}
	if got := lt.GetLanguage(); got != "en" {
		t.Errorf("lt.GetLanguage() = %v, want %v", got, "en")
	}
	if got := lt.GetLocale(); got != "US-really-long" {
		t.Errorf("lt.GetLocale() = %v, want %v", got, "US-really-long")
	}
}

// TestLanguageType_IsEmpty checks if the IsEmpty method correctly identifies empty LanguageType values.
func TestLanguageType_IsEmpty(t *testing.T) {
	tests := []struct {
		language LanguageType
		want     bool
	}{
		{"", true},
		{"en-US", false},
		{" ", true},
	}

	for _, test := range tests {
		if got := test.language.IsEmpty(); got != test.want {
			t.Errorf("LanguageType.IsEmpty() = %v, want %v", got, test.want)
		}
	}
}

// TestLanguageType_TrimSpace checks if the TrimSpace method correctly trims whitespace from LanguageType values.
func TestLanguageType_TrimSpace(t *testing.T) {
	language := LanguageType(" en-US ")
	want := LanguageType("en-US")
	if got := language.TrimSpace(); got != want {
		t.Errorf("LanguageType.TrimSpace() = %v, want %v", got, want)
	}
}

// TestLanguageType_GetParts checks if the GetParts method correctly splits LanguageType into language and locale components.
func TestLanguageType_GetParts(t *testing.T) {
	language := LanguageType("en-US")
	wantLang := "en"
	wantLocale := "US"
	if gotLang, gotLocale := language.GetParts(); gotLang != wantLang || gotLocale != wantLocale {
		t.Errorf("LanguageType.GetParts() = %v, %v, want %v, %v", gotLang, gotLocale, wantLang, wantLocale)
	}
}

// TestLanguageType_GetLanguage checks if the GetLanguage method correctly extracts the language part from LanguageType.
func TestLanguageType_GetLanguage(t *testing.T) {
	language := LanguageType("en-US")
	want := "en"
	if got := language.GetLanguage(); got != want {
		t.Errorf("LanguageType.GetLanguage() = %v, want %v", got, want)
	}
}

// TestLanguageType_GetLocale checks if the GetLocale method correctly extracts the locale part from LanguageType.
func TestLanguageType_GetLocale(t *testing.T) {
	language := LanguageType("en-US")
	want := "US"
	if got := language.GetLocale(); got != want {
		t.Errorf("LanguageType.GetLocale() = %v, want %v", got, want)
	}
}

// TestLanguageType_GetHTTPAcceptedLanguage checks if GetHTTPAcceptedLanguage correctly parses the Accept-Language header.
func TestLanguageType_GetHTTPAcceptedLanguage(t *testing.T) {
	// Test with an empty Accept-Language header.
	if got := GetHTTPAcceptedLanguage("").String(); got != "en-US" {
		t.Errorf("GetHTTPAcceptedLanguage(\"\") = %v, want %v", got, "en-US")
	}
	// Test with a valid Accept-Language header.
	if got := GetHTTPAcceptedLanguage("en-US,en;q=0.5").String(); got != "en-US" {
		t.Errorf("GetHTTPAcceptedLanguage(\"en-US,en;q=0.5\") = %v, want %v", got, "en-US")
	}
}

// TestLanguageType_GetHTTPAcceptedLanguageWithDefault checks if GetHTTPAcceptedLanguageWithDefault correctly parses the Accept-Language header and uses the default when necessary.
func TestLanguageType_GetHTTPAcceptedLanguageWithDefault(t *testing.T) {
	header := "de-DE,de;q=0.9"
	defaultLang := LanguageType("en-US")
	want := LanguageType("de-DE")
	if got := GetHTTPAcceptedLanguageWithDefault(header, defaultLang); got != want {
		t.Errorf("GetHTTPAcceptedLanguageWithDefault() = %v, want %v", got, want)
	}

	// Test with an empty header, expecting the default language to be returned.
	header = ""
	want = defaultLang
	if got := GetHTTPAcceptedLanguageWithDefault(header, defaultLang); got != want {
		t.Errorf("GetHTTPAcceptedLanguageWithDefault() = %v, want %v", got, want)
	}
}
