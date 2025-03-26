package asessions

import "testing"

func TestRoleName_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		arn      RoleName
		expected bool
	}{
		{"Empty RoleName", RoleName(""), true},
		{"Non-empty RoleName", RoleName("admin:Role Title"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arn.IsEmpty(); got != tt.expected {
				t.Errorf("RoleName.IsEmpty() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRoleName_TrimSpace(t *testing.T) {
	tests := []struct {
		name     string
		arn      RoleName
		expected RoleName
	}{
		{"Trim spaces", RoleName(" admin:Role Title "), RoleName("admin:Role Title")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arn.TrimSpace(); got != tt.expected {
				t.Errorf("RoleName.TrimSpace() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRoleName_ToStringTrimLower(t *testing.T) {
	tests := []struct {
		name     string
		arn      RoleName
		expected string
	}{
		{"Trim and lower", RoleName(" Admin:Role Title "), "admin:role title"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arn.ToStringTrimLower(); got != tt.expected {
				t.Errorf("RoleName.ToStringTrimLower() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRoleName_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		arn      RoleName
		expected bool
	}{
		{"Valid RoleName", RoleName("admin:Role Title"), true},
		{"Invalid RoleName", RoleName("admin:"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arn.IsValid(); got != tt.expected {
				t.Errorf("RoleName.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRoleName_GetRoleInfo(t *testing.T) {
	tests := []struct {
		name            string
		arn             RoleName
		expectedService RoleService
		expectedLabel   RoleLabel
	}{
		{"Valid RoleName", RoleName("admin:Role Title"), RoleService("admin"), RoleLabel("Role Title")},
		{"No RoleService", RoleName(":Role Title"), RoleService(""), RoleLabel("Role Title")},
		{"No RoleLabel", RoleName("admin:"), RoleService("admin"), RoleLabel("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotService, gotLabel := tt.arn.GetRoleInfo()
			if gotService != tt.expectedService || gotLabel != tt.expectedLabel {
				t.Errorf("RoleName.GetRoleInfo() = %v, %v, want %v, %v", gotService, gotLabel, tt.expectedService, tt.expectedLabel)
			}
		})
	}
}

func TestRoleNames_FindByRoleService(t *testing.T) {
	roleNames := RoleNames{
		RoleName("admin:Role Title"),
		RoleName("user:User Role"),
		RoleName("admin:Another Role"),
	}

	tests := []struct {
		name        string
		roleService RoleService
		expected    RoleNames
	}{
		{"Find admin roles", RoleService("admin"), RoleNames{RoleName("admin:Role Title"), RoleName("admin:Another Role")}},
		{"Find user roles", RoleService("user"), RoleNames{RoleName("user:User Role")}},
		{"Find non-existent roles", RoleService("guest"), RoleNames{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := roleNames.FindByRoleService(tt.roleService); !equalRoleNames(got, tt.expected) {
				t.Errorf("RoleNames.FindByRoleService() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRoleNames_ToString(t *testing.T) {
	roleNames := RoleNames{
		RoleName("admin:Role Title"),
		RoleName("user:User Role"),
	}

	tests := []struct {
		name     string
		sep      string
		expected string
	}{
		{"Comma separated", ", ", "admin:Role Title, user:User Role"},
		{"Pipe separated", " | ", "admin:Role Title | user:User Role"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := roleNames.ToString(tt.sep); got != tt.expected {
				t.Errorf("RoleNames.ToString() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Helper function to compare two RoleNames slices for equality.
func equalRoleNames(a, b RoleNames) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
