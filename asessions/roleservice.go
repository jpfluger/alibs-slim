package asessions

import "strings"

// RoleService represents a service associated with a role.
type RoleService string

// IsEmpty checks if the RoleService is empty after trimming whitespace.
func (rs RoleService) IsEmpty() bool {
	return strings.TrimSpace(string(rs)) == ""
}
