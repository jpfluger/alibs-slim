package autils

import (
	"encoding/json"
	"strings"
)

// StringsArray is a custom type that extends []string.
// It provides a custom unmarshaling from JSON to handle both
// an array of strings and a single comma-delimited string.
type StringsArray []string

// UnmarshalJSON is a custom unmarshaller for StringsArray.
// It supports unmarshaling from both an array of strings in JSON
// and a single comma-separated string.
func (sa *StringsArray) UnmarshalJSON(b []byte) error {
	// If the byte slice is empty, "null", or an empty JSON string, set sa to an empty StringsArray.
	if len(b) == 0 || string(b) == "null" || string(b) == `""` {
		*sa = StringsArray{}
		return nil
	}

	// Attempt to unmarshal as an array of strings.
	var values []string
	err := json.Unmarshal(b, &values)
	if err == nil {
		// If successful, clean up the strings.
		*sa = cleanStrings(values)
		return nil
	}

	// If unmarshaling as an array fails, attempt to unmarshal as a single string.
	var singleStr string
	err = json.Unmarshal(b, &singleStr)
	if err != nil {
		// If this also fails, return the error.
		return err
	}

	// If the string is "null", treat it as an empty slice.
	if singleStr == "null" {
		*sa = StringsArray{}
		return nil
	}

	// Split the string by commas and clean up the resulting strings.
	*sa = cleanStrings(strings.Split(singleStr, ","))
	return nil
}

// cleanStrings takes a slice of strings and trims whitespace from each string.
// It also filters out any empty strings.
func cleanStrings(input []string) StringsArray {
	var cleaned StringsArray
	for _, str := range input {
		trimmed := strings.TrimSpace(str)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}

// StringArrContainsString checks if a string slice contains a given string.
func StringArrContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
