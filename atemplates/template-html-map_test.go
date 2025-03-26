package atemplates

import (
	"testing"

	"github.com/gofrs/uuid/v5"
)

// TestToUpperFirst tests the ToUpperFirst function.
func TestToUpperFirst(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"Empty", "", ""},
		{"AlreadyCapital", "Test", "Test"},
		{"Lowercase", "test", "Test"},
		{"NonAlpha", "123abc", "123abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToUpperFirst(tt.s); got != tt.want {
				t.Errorf("ToUpperFirst() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIfUUIDNilElse tests the IfUUIDNilElse function.
func TestIfUUIDNilElse(t *testing.T) {
	tests := []struct {
		name      string
		target    uuid.UUID
		elseValue string
		want      string
	}{
		{"NilUUID", uuid.Nil, "empty", "empty"},
		{"NonNilUUID", uuid.Must(uuid.NewV4()), "empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IfUUIDNilElse(tt.target, tt.elseValue)
			if tt.target == uuid.Nil && got != tt.want {
				t.Errorf("IfUUIDNilElse() with nil UUID = %v, want %v", got, tt.want)
			}
			if tt.target != uuid.Nil && got == tt.want {
				t.Errorf("IfUUIDNilElse() with non-nil UUID = %v, want not %v", got, tt.want)
			}
		})
	}
}

// Additional tests for other functions can be added similarly.
