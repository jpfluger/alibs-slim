package atypeconvert

import (
	"testing"
)

// TestConvertArrayFromTo tests the ConvertArrayFromTo function for various array type conversions.
func TestConvertArrayFromTo(t *testing.T) {
	// Test conversion from array of strings to array of ints.
	input := []string{"1", "2", "3"}
	toType := "int"
	expected := []int{1, 2, 3}

	result, err := ConvertArrayFromStringTo(input, toType)
	if err != nil {
		t.Fatalf("ConvertArrayFromStringTo() error = %v", err)
	}
	if intResult, ok := result.([]int); !ok || !equalIntSlices(intResult, expected) {
		t.Errorf("ConvertArrayFromStringTo() got = %v, want %v", result, expected)
	}

	// Add more test cases as needed for other array conversions.
}

// equalIntSlices checks if two slices of ints are equal.
func equalIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// TestConvertArrayFromInterfaceTo tests the ConvertArrayFromInterfaceTo function for interface array conversions.
func TestConvertArrayFromInterfaceTo(t *testing.T) {
	// Test conversion from array of interfaces to array of strings.
	input := []interface{}{"1", "2", "3"}
	toType := "string"
	expected := []string{"1", "2", "3"}

	result, err := ConvertArrayFromInterfaceTo(input, toType)
	if err != nil {
		t.Fatalf("ConvertArrayFromInterfaceTo() error = %v", err)
	}
	if stringResult, ok := result.([]string); !ok || !equalStringSlices(stringResult, expected) {
		t.Errorf("ConvertArrayFromInterfaceTo() got = %v, want %v", result, expected)
	}

	// Add more test cases as needed for other interface array conversions.
}

// equalStringSlices checks if two slices of strings are equal.
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// TestConvertArrayFromIntTo tests the ConvertArrayFromIntTo function for integer array conversions.
func TestConvertArrayFromIntTo(t *testing.T) {
	// Test conversion from array of ints to array of strings.
	input := []int{1, 2, 3}
	toType := "string"
	expected := []string{"1", "2", "3"}

	result, err := ConvertArrayFromIntTo(input, toType)
	if err != nil {
		t.Fatalf("ConvertArrayFromIntTo() error = %v", err)
	}
	if stringResult, ok := result.([]string); !ok || !equalStringSlices(stringResult, expected) {
		t.Errorf("ConvertArrayFromIntTo() got = %v, want %v", result, expected)
	}

	// Add more test cases as needed for other integer array conversions.
}

// TestConvertArrayFromFloatTo tests the ConvertArrayFromFloatTo function for float array conversions.
func TestConvertArrayFromFloatTo(t *testing.T) {
	// Test conversion from array of floats to array of strings.
	input := []float64{1.0, 2.0, 3.0}
	toType := "string"
	expected := []string{"1.000000", "2.000000", "3.000000"}

	result, err := ConvertArrayFromFloatTo(input, toType)
	if err != nil {
		t.Fatalf("ConvertArrayFromFloatTo() error = %v", err)
	}
	if stringResult, ok := result.([]string); !ok || !equalStringSlices(stringResult, expected) {
		t.Errorf("ConvertArrayFromFloatTo() got = %v, want %v", result, expected)
	}

	// Add more test cases as needed for other float array conversions.
}

// TestConvertArrayFromBoolTo tests the ConvertArrayFromBoolTo function for boolean array conversions.
func TestConvertArrayFromBoolTo(t *testing.T) {
	// Test conversion from array of booleans to array of strings.
	input := []bool{true, false, true}
	toType := "string"
	expected := []string{"true", "false", "true"}

	result, err := ConvertArrayFromBoolTo(input, toType)
	if err != nil {
		t.Fatalf("ConvertArrayFromBoolTo() error = %v", err)
	}
	if stringResult, ok := result.([]string); !ok || !equalStringSlices(stringResult, expected) {
		t.Errorf("ConvertArrayFromBoolTo() got = %v, want %v", result, expected)
	}

	// Add more test cases as needed for other boolean array conversions.
}
