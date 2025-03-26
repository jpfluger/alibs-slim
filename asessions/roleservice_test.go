package asessions

import "testing"

func TestRoleService_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		rs       RoleService
		expected bool
	}{
		{"Empty RoleService", RoleService(""), true},
		{"Non-empty RoleService", RoleService("admin"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rs.IsEmpty(); got != tt.expected {
				t.Errorf("RoleService.IsEmpty() = %v, want %v", got, tt.expected)
			}
		})
	}
}
