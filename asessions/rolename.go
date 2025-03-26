package asessions

import (
	"github.com/jpfluger/alibs-slim/autils"
	"strings"
)

// RoleName combines RoleService and RoleLabel separated by a colon (":").
// Examples:
// * admin:Role Title
// * domain:MyDomain
// * org:Admin
// * org:Operator
type RoleName string

// IsEmpty checks if the RoleName is empty after trimming whitespace.
func (arn RoleName) IsEmpty() bool {
	return strings.TrimSpace(string(arn)) == ""
}

// TrimSpace trims leading and trailing whitespace from the RoleName.
func (arn RoleName) TrimSpace() RoleName {
	return RoleName(strings.TrimSpace(string(arn)))
}

// String returns the RoleName as a trimmed string.
func (arn RoleName) String() string {
	return strings.TrimSpace(string(arn))
}

// ToStringTrimLower returns the RoleName as a trimmed and lowercase string.
func (arn RoleName) ToStringTrimLower() string {
	return autils.ToStringTrimLower(arn.String())
}

// IsValid checks if the RoleName is valid by ensuring both RoleService and RoleLabel are not empty.
func (arn RoleName) IsValid() bool {
	tierService, tierLabel := arn.GetRoleInfo()
	return !tierService.IsEmpty() && !tierLabel.IsEmpty()
}

// GetRoleInfo splits the RoleName into RoleService and RoleLabel.
func (arn RoleName) GetRoleInfo() (RoleService, RoleLabel) {
	ss := strings.Split(string(arn), ":")
	if len(ss) == 0 {
		return "", ""
	} else if len(ss) == 1 {
		return "", RoleLabel(ss[0])
	}
	return RoleService(ss[0]), RoleLabel(ss[1])
}

// RoleNames represents a slice of RoleName.
type RoleNames []RoleName

// FindByRoleService finds all RoleNames with the specified RoleService.
func (nms RoleNames) FindByRoleService(roleService RoleService) RoleNames {
	var arr RoleNames
	for _, nm := range nms {
		rs, _ := nm.GetRoleInfo()
		if rs == roleService {
			arr = append(arr, nm)
		}
	}
	return arr
}

// ToString converts RoleNames to a single string separated by the specified separator.
func (nms RoleNames) ToString(sep string) string {
	var sb strings.Builder
	for ii, nm := range nms {
		if ii > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(string(nm))
	}
	return sb.String()
}

// Contains checks if the RoleNames slice contains the specified RoleName.
func (nms RoleNames) Contains(name RoleName) bool {
	if name.IsEmpty() {
		return false
	}
	for _, nm := range nms {
		if nm == name {
			return true
		}
	}
	return false
}
