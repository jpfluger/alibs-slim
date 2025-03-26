package atypeconvert

// DetectType determines the type of the given value and whether it is an array.
// It returns a string representing the type and a boolean indicating if it's an array.
func DetectType(value interface{}) (valueType string, isArray bool) {
	// If the value is nil, no type can be detected.
	if value == nil {
		return "", false
	}

	// Use a type switch to determine the type of the value.
	switch value.(type) {
	case string, *string:
		// Detected a string or pointer to a string.
		return "string", false
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8,
		*int, *int64, *int32, *int16, *int8, *uint, *uint64, *uint32, *uint16, *uint8:
		// Detected an integer or pointer to an integer of various sizes.
		return "int", false
	case float64, float32, *float64, *float32:
		// Detected a float or pointer to a float of various sizes.
		return "float", false
	case bool, *bool:
		// Detected a boolean or pointer to a boolean.
		return "bool", false
	case map[string]interface{}:
		// Detected a map with string keys and interface{} values.
		return "obj", false
	case []interface{}:
		// Detected a slice of interfaces.
		return "arr", true
	case []string:
		// Detected a slice of strings.
		return "string", true
	case []int, []int64, []int32, []int16, []int8, []uint, []uint64, []uint32, []uint16, []uint8:
		// Detected a slice of integers of various sizes.
		return "int", true
	case []float64, []float32:
		// Detected a slice of floats of various sizes.
		return "float", true
	case []bool:
		// Detected a slice of booleans.
		return "bool", true
	default:
		// If none of the above types match, return an empty string and false.
		return "", false
	}
}
