package autils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalJSON_Array tests unmarshaling from a JSON array of strings.
func TestUnmarshalJSON_Array(t *testing.T) {
	input := `["one", "two", "three"]`
	expected := StringsArray{"one", "two", "three"}

	var sa StringsArray
	err := json.Unmarshal([]byte(input), &sa)

	assert.NoError(t, err, "Unmarshaling should not produce an error")
	assert.Equal(t, expected, sa, "Unmarshaled array should match expected value")
}

// TestUnmarshalJSON_CommaDelimited tests unmarshaling from a comma-delimited string.
func TestUnmarshalJSON_CommaDelimited(t *testing.T) {
	input := `"one,two,three"`
	expected := StringsArray{"one", "two", "three"}

	var sa StringsArray
	err := json.Unmarshal([]byte(input), &sa)

	assert.NoError(t, err, "Unmarshaling should not produce an error")
	assert.Equal(t, expected, sa, "Unmarshaled string should match expected value")
}

// TestUnmarshalJSON_EmptyString tests unmarshaling from an empty JSON string.
func TestUnmarshalJSON_EmptyString(t *testing.T) {
	input := `""`
	expected := StringsArray{}

	var sa StringsArray
	err := json.Unmarshal([]byte(input), &sa)

	assert.NoError(t, err, "Unmarshaling should not produce an error")
	assert.Equal(t, expected, sa, "Unmarshaled empty string should produce an empty StringsArray")
}

// TestUnmarshalJSON_Null tests unmarshaling from a JSON null.
func TestUnmarshalJSON_Null(t *testing.T) {
	input := `null`
	expected := StringsArray{} // An empty StringsArray, not nil.

	var sa StringsArray
	err := json.Unmarshal([]byte(input), &sa)

	assert.NoError(t, err, "Unmarshaling should not produce an error")
	assert.Equal(t, expected, sa, "Unmarshaled null should produce an empty StringsArray")
}

// TestUnmarshalJSON_InvalidJSON tests unmarshaling from an invalid JSON.
func TestUnmarshalJSON_InvalidJSON(t *testing.T) {
	input := `not a valid json`

	var sa StringsArray
	err := json.Unmarshal([]byte(input), &sa)

	assert.Error(t, err, "Unmarshaling should produce an error for invalid JSON")
}

// TestUnmarshalJSON_MixedWhitespace tests unmarshaling from a string with mixed whitespace.
func TestUnmarshalJSON_MixedWhitespace(t *testing.T) {
	input := `" one , two , three "`
	expected := StringsArray{"one", "two", "three"}

	var sa StringsArray
	err := json.Unmarshal([]byte(input), &sa)

	assert.NoError(t, err, "Unmarshaling should not produce an error")
	assert.Equal(t, expected, sa, "Unmarshaled string should match expected value with trimmed whitespace")
}

// TestUnmarshalJSON_EmptyElements tests unmarshaling from a string with empty elements.
func TestUnmarshalJSON_EmptyElements(t *testing.T) {
	input := `"one,, ,two, ,three,"`
	expected := StringsArray{"one", "two", "three"}

	var sa StringsArray
	err := json.Unmarshal([]byte(input), &sa)

	assert.NoError(t, err, "Unmarshaling should not produce an error")
	assert.Equal(t, expected, sa, "Unmarshaled string should match expected value without empty elements")
}
