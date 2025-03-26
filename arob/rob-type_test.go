package arob

import (
	"testing"
)

// TestROBType_IsEmpty tests the IsEmpty method to ensure it correctly identifies empty ROBType values.
func TestROBType_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		rt   ROBType
		want bool
	}{
		{"Empty", "", true},
		{"WhitespaceOnly", "   ", true},
		{"NonEmpty", "error", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rt.IsEmpty(); got != tt.want {
				t.Errorf("ROBType.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestROBType_TrimSpace tests the TrimSpace method to ensure it correctly trims whitespace from ROBType values.
func TestROBType_TrimSpace(t *testing.T) {
	tests := []struct {
		name string
		rt   ROBType
		want ROBType
	}{
		{"NoTrim", "error", "error"},
		{"TrimLeading", "  error", "error"},
		{"TrimTrailing", "error  ", "error"},
		{"TrimBoth", "  error  ", "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rt.TrimSpace(); got != tt.want {
				t.Errorf("ROBType.TrimSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestROBType_Validate tests the Validate method to ensure it correctly validates ROBType values.
func TestROBType_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rt      ROBType
		wantErr bool
	}{
		{"Valid", "error", false},
		{"InvalidChar", "error!", true},
		{"InvalidStart", "-error", true},
		{"InvalidEnd", "error-", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rt.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ROBType.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestROBTypes_HasValues tests the HasValues method to ensure it correctly identifies if ROBTypes has any values.
func TestROBTypes_HasValues(t *testing.T) {
	tests := []struct {
		name string
		rts  ROBTypes
		want bool
	}{
		{"Empty", nil, false},
		{"HasValues", []ROBType{"error"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rts.HasValues(); got != tt.want {
				t.Errorf("ROBTypes.HasValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Additional tests for other methods can be added similarly.
