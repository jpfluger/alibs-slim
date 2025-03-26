package atypeconvert

import (
	"testing"
)

// TestConvertToStringFrom tests the ConvertToStringFrom function for various types.
func TestConvertToStringFrom(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{"test", "test"},
		{123, "123"},
		{123.456, "123.456000"},
		{true, "true"},
		{false, "false"},
	}

	for _, test := range tests {
		result, err := ConvertToStringFrom(test.input)
		if err != nil {
			t.Errorf("ConvertToStringFrom(%v) unexpected error: %v", test.input, err)
		}
		if result != test.expected {
			t.Errorf("ConvertToStringFrom(%v) = %v, want %v", test.input, result, test.expected)
		}
	}
}

// TestConvertToIntFrom tests the ConvertToIntFrom function for various types.
func TestConvertToIntFrom(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected int
	}{
		{"123", 123},
		{123, 123},
		{123.456, 123},
		{true, 1},
		{false, 0},
	}

	for _, test := range tests {
		result, err := ConvertToIntFrom(test.input)
		if err != nil {
			t.Errorf("ConvertToIntFrom(%v) unexpected error: %v", test.input, err)
		}
		if result != test.expected {
			t.Errorf("ConvertToIntFrom(%v) = %v, want %v", test.input, result, test.expected)
		}
	}
}

// TestConvertToFloatFrom tests the ConvertToFloatFrom function for various types.
func TestConvertToFloatFrom(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected float64
	}{
		{"123.456", 123.456},
		{123, 123.0},
		{123.456, 123.456},
		{true, 1.0},
		{false, 0.0},
	}

	for _, test := range tests {
		result, err := ConvertToFloatFrom(test.input)
		if err != nil {
			t.Errorf("ConvertToFloatFrom(%v) unexpected error: %v", test.input, err)
		}
		if result != test.expected {
			t.Errorf("ConvertToFloatFrom(%v) = %v, want %v", test.input, result, test.expected)
		}
	}
}

// TestConvertToBoolFrom tests the ConvertToBoolFrom function for various types.
func TestConvertToBoolFrom(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected bool
	}{
		{"true", true},
		{"false", false},
		{1, true},
		{0, false},
		{123.456, true},
		{0.0, false},
		{true, true},
		{false, false},
	}

	for _, test := range tests {
		result, err := ConvertToBoolFrom(test.input)
		if err != nil {
			t.Errorf("ConvertToBoolFrom(%v) unexpected error: %v", test.input, err)
		}
		if result != test.expected {
			t.Errorf("ConvertToBoolFrom(%v) = %v, want %v", test.input, result, test.expected)
		}
	}
}
