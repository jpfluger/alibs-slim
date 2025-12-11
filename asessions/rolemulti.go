package asessions

import (
	"fmt"
)

// RoleMulti represents a set of role names and associated permissions.
type RoleMulti struct {
	// Names could be "admin", "domain-admin", "user-admin" or a combination.
	// The idea behind this is to add together all the roles then
	// apply PermsPlus or PermsMinus to the entire roles set.
	Names RoleNames `json:"names,omitempty"`

	// PermsPlus is applied first.
	PermsPlus PermSet `json:"permsPlus,omitempty"`

	// PermsMinus is applied after PermsPlus.
	PermsMinus PermSet `json:"permsMinus,omitempty"`
}

// Validate ensures that essential fields in the Role struct have values.
// If any fields are nil (Names, PermsPlus, or PermsMinus), it initializes them with empty values.
// It returns an error if the Role struct is nil or if required fields like Names are empty after initialization.
func (r *RoleMulti) Validate() error {
	if r == nil {
		return fmt.Errorf("role is nil")
	}
	if r.Names == nil {
		r.Names = RoleNames{}
	}
	if r.PermsPlus == nil {
		r.PermsPlus = PermSet{}
	}
	if r.PermsMinus == nil {
		r.PermsMinus = PermSet{}
	}
	if len(r.Names) == 0 {
		return fmt.Errorf("role names cannot be empty")
	}
	// Check for empty perms (optional but stricter)
	//if len(r.PermsPlus) == 0 && len(r.PermsMinus) == 0 {
	//	return fmt.Errorf("role perms cannot be completely empty")
	//}
	// Optional: Validate PermSets deeper
	// Validate deeper and propagate errors
	if err := r.PermsPlus.Validate(); err != nil {
		return fmt.Errorf("invalid PermsPlus for role %s: %w", r.Names[0], err)
	}
	if err := r.PermsMinus.Validate(); err != nil {
		return fmt.Errorf("invalid PermsMinus for role %s: %w", r.Names[0], err)
	}
	return nil
}

// BuildLimitedPermSet constructs a PermSet for the Role using the specified RoleFactory and ensures
// that permissions do not exceed those in the provided permSetLimit.
// It first validates the Role and verifies that RoleFactory and permSetLimit are non-nil.
// If validation passes, it returns a PermSet that merges the RoleFactory permissions with the Role's PermsPlus and PermsMinus,
// constrained by permSetLimit. Returns an error if validation fails or input values are invalid.
func (r *RoleMulti) BuildLimitedPermSet(factory RoleFactory, permSetLimit PermSet) (PermSet, error) {
	if factory == nil {
		return nil, fmt.Errorf("role factory is nil")
	}
	if permSetLimit == nil {
		return nil, fmt.Errorf("permSetLimit is nil")
	}
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if len(r.Names) == 0 && len(r.PermsPlus) == 0 {
		return nil, fmt.Errorf("role names or permsPlus is empty")
	}
	return factory.BuildPermSetWithLimitMulti(r, permSetLimit), nil
}

func (r *RoleMulti) Clone() *RoleMulti {
	if r == nil {
		return nil
	}
	return &RoleMulti{
		Names:      append(RoleNames{}, r.Names...), // Create a copy of Names slice
		PermsPlus:  r.PermsPlus.Clone(),             // Clone the PermsPlus PermSet
		PermsMinus: r.PermsMinus.Clone(),            // Clone the PermsMinus PermSet
	}
}

// RoleMultis represents a slice of Role pointers.
type RoleMultis []*RoleMulti

func (rs RoleMultis) Validate() error {
	if rs == nil || len(rs) == 0 {
		return nil
	}
	for ii, r := range rs {
		if err := r.Validate(); err != nil {
			return fmt.Errorf("role %d is invalid: %w", ii, err)
		}
	}
	return nil
}

// GetFirst returns the first role otherwise nil, if it does not exist.
func (rs RoleMultis) GetFirst() *RoleMulti {
	if rs == nil || len(rs) == 0 {
		return nil
	}
	return rs[0]
}

// SetRole adds a new role to the Roles slice or updates an existing one by name.
func (rs *RoleMultis) SetRole(role *RoleMulti) {
	if role == nil || len(role.Names) == 0 {
		return // Ignore invalid roles
	}

	if role.PermsPlus == nil {
		role.PermsPlus = PermSet{}
	}
	if role.PermsMinus == nil {
		role.PermsMinus = PermSet{}
	}

	for _, newRoleName := range role.Names {
		for i, existingRole := range *rs {
			if existingRole.Names.Contains(newRoleName) {
				// Update the existing role
				(*rs)[i] = role
				return
			}
		}
	}

	// Role does not exist, so add it
	*rs = append(*rs, role)
}

// RemoveRole removes a role from the Roles slice by its index.
func (rs *RoleMultis) RemoveRole(index int) error {
	if index < 0 || index >= len(*rs) {
		return fmt.Errorf("index out of range")
	}
	*rs = append((*rs)[:index], (*rs)[index+1:]...)
	return nil
}

// FindRoleByName finds a role by its name in the Roles slice.
func (rs RoleMultis) FindRoleByName(name RoleName) *RoleMulti {
	for _, role := range rs {
		for _, roleName := range role.Names {
			if roleName == name {
				return role
			}
		}
	}
	return nil
}

type RoleMultiMap map[RoleName]*RoleMulti

// Set adds or updates a RoleMulti in the RoleMultiMap.
func (rm RoleMultiMap) Set(role *RoleMulti) {
	if role == nil || len(role.Names) == 0 {
		return // Ignore invalid roles
	}
	for _, name := range role.Names {
		rm[name] = role.Clone() // Add or update with a cloned role
	}
}

// Get retrieves a RoleMulti from the RoleMultiMap by RoleName.
func (rm RoleMultiMap) Get(name RoleName) *RoleMulti {
	if role, exists := rm[name]; exists {
		return role.Clone() // Return a cloned role to prevent external modification
	}
	return nil
}

// Remove deletes a RoleMulti from the RoleMultiMap by RoleName.
func (rm RoleMultiMap) Remove(name RoleName) {
	delete(rm, name)
}

// Exists checks if a RoleMulti exists in the RoleMultiMap by RoleName.
func (rm RoleMultiMap) Exists(name RoleName) bool {
	_, exists := rm[name]
	return exists
}

// Clone creates a deep copy of the RoleMultiMap.
func (rm RoleMultiMap) Clone() RoleMultiMap {
	clonedMap := RoleMultiMap{}
	for name, role := range rm {
		if role != nil {
			clonedMap[name] = role.Clone() // Deep clone each RoleMulti
		}
	}
	return clonedMap
}
