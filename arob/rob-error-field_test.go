package arob

import (
	"testing"
)

// TestROBErrorField_IsEmpty tests the IsEmpty method to ensure it correctly identifies empty fields.
func TestROBErrorField_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		fld  ROBErrorField
		want bool
	}{
		{"Empty", "", true},
		{"WhitespaceOnly", "   ", true},
		{"NonEmpty", "field", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fld.IsEmpty(); got != tt.want {
				t.Errorf("ROBErrorField.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestROBErrorField_TrimSpace tests the TrimSpace method to ensure it correctly trims whitespace.
func TestROBErrorField_TrimSpace(t *testing.T) {
	tests := []struct {
		name string
		fld  ROBErrorField
		want ROBErrorField
	}{
		{"NoTrim", "field", "field"},
		{"TrimLeading", "  field", "field"},
		{"TrimTrailing", "field  ", "field"},
		{"TrimBoth", "  field  ", "field"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fld.TrimSpace(); got != tt.want {
				t.Errorf("ROBErrorField.TrimSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestROBErrorField_GetType tests the GetType method to ensure it correctly extracts the type part of the field.
func TestROBErrorField_GetType(t *testing.T) {
	tests := []struct {
		name string
		fld  ROBErrorField
		want string
	}{
		{"NoColon", "field", "field"},
		{"WithColon", "field:type", "field"},
		{"MultipleColons", "field:type:subtype", "field"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fld.GetType(); got != tt.want {
				t.Errorf("ROBErrorField.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestROBErrorField_GetSubType tests the GetSubType method to ensure it correctly extracts the subtype part of the field.
func TestROBErrorField_GetSubType(t *testing.T) {
	tests := []struct {
		name string
		fld  ROBErrorField
		want string
	}{
		{"NoColon", "field", ""},
		{"WithColon", "field:type", "type"},
		{"MultipleColons", "field:type:subtype", "type:subtype"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fld.GetSubType(); got != tt.want {
				t.Errorf("ROBErrorField.GetSubType() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Additional tests for other methods can be added similarly.
