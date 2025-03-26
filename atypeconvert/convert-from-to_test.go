package atypeconvert

import (
	"testing"
)

// TestConvertFromTo tests the ConvertFromTo function for various type conversions.
func TestConvertFromTo(t *testing.T) {
	// Test string to int conversion.
	value, err := ConvertFromTo("123", "int")
	if err != nil {
		t.Fatalf("ConvertFromTo() error = %v", err)
	}
	if intValue, ok := value.(int); !ok || intValue != 123 {
		t.Errorf("ConvertFromTo() got = %v, want %v", value, 123)
	}

	// Test int to string conversion.
	value, err = ConvertFromTo(123, "string")
	if err != nil {
		t.Fatalf("ConvertFromTo() error = %v", err)
	}
	if strValue, ok := value.(string); !ok || strValue != "123" {
		t.Errorf("ConvertFromTo() got = %v, want %v", value, "123")
	}

	// Add more test cases as needed for other conversions.
}

// TestConvertFromStringTo tests the ConvertFromStringTo function for string conversions.
func TestConvertFromStringTo(t *testing.T) {
	// Test string to bool conversion.
	value, err := ConvertFromStringTo("true", "bool")
	if err != nil {
		t.Fatalf("ConvertFromStringTo() error = %v", err)
	}
	if boolValue, ok := value.(bool); !ok || !boolValue {
		t.Errorf("ConvertFromStringTo() got = %v, want %v", value, true)
	}

	// Add more test cases as needed for other string conversions.
}

// TestConvertFromIntTo tests the ConvertFromIntTo function for integer conversions.
func TestConvertFromIntTo(t *testing.T) {
	// Test int to float conversion.
	value, err := ConvertFromIntTo(123, "float")
	if err != nil {
		t.Fatalf("ConvertFromIntTo() error = %v", err)
	}
	if floatValue, ok := value.(float64); !ok || floatValue != 123.0 {
		t.Errorf("ConvertFromIntTo() got = %v, want %v", value, 123.0)
	}

	// Add more test cases as needed for other integer conversions.
}

// TestConvertFromFloatTo tests the ConvertFromFloatTo function for float conversions.
func TestConvertFromFloatTo(t *testing.T) {
	// Test float to int conversion.
	value, err := ConvertFromFloatTo(123.45, "int")
	if err != nil {
		t.Fatalf("ConvertFromFloatTo() error = %v", err)
	}
	if intValue, ok := value.(int); !ok || intValue != 123 {
		t.Errorf("ConvertFromFloatTo() got = %v, want %v", value, 123)
	}

	// Add more test cases as needed for other float conversions.
}

// TestConvertFromBoolTo tests the ConvertFromBoolTo function for boolean conversions.
func TestConvertFromBoolTo(t *testing.T) {
	// Test bool to string conversion.
	value, err := ConvertFromBoolTo(true, "string")
	if err != nil {
		t.Fatalf("ConvertFromBoolTo() error = %v", err)
	}
	if strValue, ok := value.(string); !ok || strValue != "true" {
		t.Errorf("ConvertFromBoolTo() got = %v, want %v", value, "true")
	}

	// Add more test cases as needed for other boolean conversions.
}
