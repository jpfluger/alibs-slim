package asessions

import (
	"fmt"
)

// Role represents a set of role names and associated permissions.
type Role struct {
	// Name could be "admin", "domain-admin", "user-admin" or a combination.
	// The idea behind this is to add together all the roles then
	// apply PermsPlus or PermsMinus to the entire roles set.
	Name RoleName `json:"name,omitempty"`

	// PermsPlus is applied first.
	PermsPlus PermSet `json:"permsPlus,omitempty"`

	// PermsMinus is applied after PermsPlus.
	PermsMinus PermSet `json:"permsMinus,omitempty"`
}

// Validate ensures that essential fields in the Role struct have values.
// If any fields are nil (Names, PermsPlus, or PermsMinus), it initializes them with empty values.
// It returns an error if the Role struct is nil or if required fields like Names are empty after initialization.
func (r *Role) Validate() error {
	if r == nil {
		return fmt.Errorf("role is nil")
	}
	if r.Name.IsEmpty() {
		return fmt.Errorf("role name is empty")
	}
	if r.PermsPlus == nil {
		r.PermsPlus = PermSet{}
	}
	if r.PermsMinus == nil {
		r.PermsMinus = PermSet{}
	}
	return nil
}

// BuildLimitedPermSet constructs a PermSet for the Role using the specified RoleFactory and ensures
// that permissions do not exceed those in the provided permSetLimit.
// It first validates the Role and verifies that RoleFactory and permSetLimit are non-nil.
// If validation passes, it returns a PermSet that merges the RoleFactory permissions with the Role's PermsPlus and PermsMinus,
// constrained by permSetLimit. Returns an error if validation fails or input values are invalid.
func (r *Role) BuildLimitedPermSet(factory RoleFactory, permSetLimit PermSet) (PermSet, error) {
	if factory == nil {
		return nil, fmt.Errorf("role factory is nil")
	}
	if permSetLimit == nil {
		return nil, fmt.Errorf("permSetLimit is nil")
	}
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if r.Name.IsValid() && len(r.PermsPlus) == 0 {
		return nil, fmt.Errorf("role name or permsPlus is empty")
	}
	return factory.BuildPermSetWithLimit(r, permSetLimit), nil
}

func (r *Role) Clone() *Role {
	if r == nil {
		return nil
	}
	r.PermsPlus.Validate()  // Clean invalid entries
	r.PermsMinus.Validate() // Clean invalid entries
	return &Role{
		Name:       r.Name,               // Copy the role name
		PermsPlus:  r.PermsPlus.Clone(),  // Clone the PermsPlus PermSet
		PermsMinus: r.PermsMinus.Clone(), // Clone the PermsMinus PermSet
	}
}

// Roles represents a slice of Role pointers.
type Roles []*Role

func (rs Roles) Validate() error {
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
func (rs Roles) GetFirst() *Role {
	if rs == nil || len(rs) == 0 {
		return nil
	}
	return rs[0]
}

// SetRole adds a new role to the Roles slice or updates an existing one by name.
func (rs *Roles) SetRole(role *Role) {
	if role == nil || role.Name.IsEmpty() {
		return // Ignore invalid roles
	}

	for i, existingRole := range *rs {
		if existingRole.Name == role.Name {
			// Update the existing role
			(*rs)[i] = role
			return
		}
	}

	// Role does not exist, so add it
	*rs = append(*rs, role)
}

// RemoveRole removes a role from the Roles slice by its index.
func (rs *Roles) RemoveRole(index int) error {
	if index < 0 || index >= len(*rs) {
		return fmt.Errorf("index out of range")
	}
	*rs = append((*rs)[:index], (*rs)[index+1:]...)
	return nil
}

// FindRoleByName finds a role by its name in the Roles slice.
func (rs Roles) FindRoleByName(name RoleName) *Role {
	for _, role := range rs {
		if role.Name == name {
			return role
		}
	}
	return nil
}

type RoleMap map[RoleName]*Role

// Set adds or updates a role in the RoleMap.
// If a role with the same name already exists, it is replaced.
func (rm RoleMap) Set(role *Role) {
	if role == nil || role.Name.IsEmpty() {
		return // Ignore nil or invalid roles
	}
	if role.PermsPlus == nil {
		role.PermsPlus = PermSet{}
	}
	if role.PermsMinus == nil {
		role.PermsMinus = PermSet{}
	}
	rm[role.Name] = role
}

// Remove deletes a role from the RoleMap by its name.
func (rm RoleMap) Remove(roleName RoleName) {
	delete(rm, roleName)
}

// Get retrieves a role from the RoleMap by its name.
func (rm RoleMap) Get(roleName RoleName) *Role {
	return rm[roleName]
}

// Contains checks if a role with the given name exists in the RoleMap.
func (rm RoleMap) Contains(roleName RoleName) bool {
	_, exists := rm[roleName]
	return exists
}

// Merge merges another RoleMap into the current RoleMap.
// Existing roles are updated, and new roles are added.
func (rm RoleMap) Merge(other RoleMap) {
	for _, role := range other {
		rm.Set(role)
	}
}

// Clone creates a deep copy of the RoleMap.
func (rm RoleMap) Clone() RoleMap {
	clone := RoleMap{}
	for roleName, role := range rm {
		clone[roleName] = role.Clone()
	}
	return clone
}
