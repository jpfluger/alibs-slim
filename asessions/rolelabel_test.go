package asessions

import (
	"testing"
)

func TestRoleLabel_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		rl       RoleLabel
		expected bool
	}{
		{"Empty RoleLabel", RoleLabel(""), true},
		{"Non-empty RoleLabel", RoleLabel("Role Title"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rl.IsEmpty(); got != tt.expected {
				t.Errorf("RoleLabel.IsEmpty() = %v, want %v", got, tt.expected)
			}
		})
	}
}
