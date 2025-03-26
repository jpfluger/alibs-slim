package autils

import (
	"testing"
)

// TestStrings_ToStringTrimLower checks if the ToStringTrimLower function correctly trims and converts a string to lowercase.
func TestStrings_ToStringTrimLower(t *testing.T) {
	input := " HeLLo WoRLD "
	expected := "hello world"
	if got := ToStringTrimLower(input); got != expected {
		t.Errorf("ToStringTrimLower() = %v, want %v", got, expected)
	}
}

// TestStrings_ToStringTrimUpper checks if the ToStringTrimUpper function correctly trims and converts a string to uppercase.
func TestStrings_ToStringTrimUpper(t *testing.T) {
	input := " HeLLo WoRLD "
	expected := "HELLO WORLD"
	if got := ToStringTrimUpper(input); got != expected {
		t.Errorf("ToStringTrimUpper() = %v, want %v", got, expected)
	}
}

// TestStrings_HasPrefixPath checks if the HasPrefixPath function correctly identifies if a string has any of the provided prefixes.
func TestStrings_HasPrefixPath(t *testing.T) {
	prefixes := []string{"/home/", "/var/"}
	input := "/home/user"
	if !HasPrefixPath(input, prefixes) {
		t.Errorf("HasPrefixPath() = false, want true")
	}
}

// TestStrings_ExtractPrefixBrackets checks if the ExtractPrefixBrackets function correctly extracts content within brackets.
func TestStrings_ExtractPrefixBrackets(t *testing.T) {
	input := "[inner]outer"
	expectedInner := "inner"
	expectedOuter := "outer"
	inner, outer := ExtractPrefixBrackets(input)
	if inner != expectedInner || outer != expectedOuter {
		t.Errorf("ExtractPrefixBrackets() = %v, %v, want %v, %v", inner, outer, expectedInner, expectedOuter)
	}
}

// TestStrings_ExtractPrefixParenthesis checks if the ExtractPrefixParenthesis function correctly extracts content within parentheses.
func TestStrings_ExtractPrefixParenthesis(t *testing.T) {
	input := "(inner)outer"
	expectedInner := "inner"
	expectedOuter := "outer"
	inner, outer := ExtractPrefixParenthesis(input)
	if inner != expectedInner || outer != expectedOuter {
		t.Errorf("ExtractPrefixParenthesis() = %v, %v, want %v, %v", inner, outer, expectedInner, expectedOuter)
	}
}

// TestStrings_ExtractPrefixBraces checks if the ExtractPrefixBraces function correctly extracts content within braces.
func TestStrings_ExtractPrefixBraces(t *testing.T) {
	input := "{inner}outer"
	expectedInner := "inner"
	expectedOuter := "outer"
	inner, outer := ExtractPrefixBraces(input)
	if inner != expectedInner || outer != expectedOuter {
		t.Errorf("ExtractPrefixBraces() = %v, %v, want %v, %v", inner, outer, expectedInner, expectedOuter)
	}
}

// TestStrings_ExtractPlaceholderLeftRight checks if the ExtractPlaceholderLeftRight function correctly extracts content between delimiters.
func TestStrings_ExtractPlaceholderLeftRight(t *testing.T) {
	input := "{inner}outer"
	left := "{"
	right := "}"
	expectedInner := "inner"
	expectedOuter := "outer"
	inner, outer := ExtractPlaceholderLeftRight(input, left, right)
	if inner != expectedInner || outer != expectedOuter {
		t.Errorf("ExtractPlaceholderLeftRight() = %v, %v, want %v, %v", inner, outer, expectedInner, expectedOuter)
	}
}
