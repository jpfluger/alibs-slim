package asessions

import "strings"

// RoleLabel represents a label associated with a role.
type RoleLabel string

// IsEmpty checks if the RoleLabel is empty after trimming whitespace.
func (rl RoleLabel) IsEmpty() bool {
	return strings.TrimSpace(string(rl)) == ""
}
