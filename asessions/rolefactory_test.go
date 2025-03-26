package asessions

import (
	"reflect"
	"testing"
)

func TestRoles_Merge(t *testing.T) {
	roleFactory := RoleFactory{
		"admin": NewPermSetByString([]string{"self:XLCRUD", "dash:XLCRUD", "bills:XLCRUD"}),
		"user":  NewPermSetByString([]string{"self:XR", "dash:XLR", "bills:"}),
		"api":   NewPermSetByString([]string{"self:XR", "dash:R", "bills:R"}),
	}

	tests := []struct {
		roleNames   RoleNames
		permsPlus   PermSet
		permsMinus  PermSet
		expected    PermSet
		description string
	}{
		{
			roleNames:   RoleNames{"admin"},
			permsPlus:   NewPermSetByString([]string{"reports:CRUD"}),
			permsMinus:  NewPermSetByString([]string{"dash:L"}),
			expected:    NewPermSetByString([]string{"self:XLCRUD", "dash:XCRUD", "bills:XLCRUD", "reports:CRUD"}),
			description: "admin role with additional permissions for reports and minus dash:L",
		},
		{
			roleNames:   RoleNames{"user", "api"},
			permsPlus:   NewPermSetByString([]string{"files:XR"}),
			permsMinus:  NewPermSetByString([]string{"self:R"}),
			expected:    NewPermSetByString([]string{"self:X", "dash:XLR", "bills:R", "files:XR"}),
			description: "user and api roles with additional files:XR permission and removal of self:R",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			// Step 1: Initialize a PermSet for merging
			mergedPermSet := PermSet{}

			// Step 2: Merge permissions from RoleFactory for each role name
			for _, roleName := range test.roleNames {
				if perms, ok := roleFactory[roleName]; ok {
					mergedPermSet.MergeByPermSet(perms)
				}
			}

			// Step 3: Apply PermPlus to add permissions
			mergedPermSet.MergeByPermSet(test.permsPlus)

			// Step 4: Apply PermMinus to subtract permissions
			mergedPermSet.SubtractByPermSet(test.permsMinus)

			// Step 5: Compare the merged permissions with the expected result
			if !reflect.DeepEqual(mergedPermSet, test.expected) {
				t.Errorf("failed for %s: expected %v, got %v", test.description, test.expected.ToStringArray(), mergedPermSet.ToStringArray())
			}
		})
	}
}

func TestRoleFactory_BuildPermSet(t *testing.T) {
	roleFactory := RoleFactory{
		"admin": NewPermSetByString([]string{"self:XLCRUD", "dash:XLCRUD", "bills:XLCRUD"}),
		"user":  NewPermSetByString([]string{"self:XR", "dash:XLR", "bills:"}),
		"api":   NewPermSetByString([]string{"self:XR", "dash:R", "bills:R"}),
	}

	role := &RoleMulti{
		Names:      RoleNames{"admin", "user"},
		PermsPlus:  NewPermSetByString([]string{"reports:CRUD"}),
		PermsMinus: NewPermSetByString([]string{"dash:L"}),
	}

	expectedPermSet := NewPermSetByString([]string{"self:XLCRUD", "dash:XCRUD", "bills:XLCRUD", "reports:CRUD"})

	// Call the BuildPermSet function
	resultPermSet := roleFactory.BuildPermSetMulti(role)

	// Validate the result against the expected PermSet
	if !reflect.DeepEqual(resultPermSet, expectedPermSet) {
		t.Errorf("BuildPermSet failed: expected %v, got %v", expectedPermSet.ToStringArray(), resultPermSet.ToStringArray())
	}
}

func TestRoleFactory_BuildPermSetWithLimit(t *testing.T) {
	roleFactory := RoleFactory{
		"admin": NewPermSetByString([]string{"self:XLCRUD", "dash:XLCRUD", "bills:XLCRUD"}),
		"user":  NewPermSetByString([]string{"self:XR", "dash:XLR", "bills:"}),
	}

	role := &RoleMulti{
		Names:      RoleNames{"admin", "user"},
		PermsPlus:  NewPermSetByString([]string{"reports:CRUD"}),
		PermsMinus: NewPermSetByString([]string{"dash:L"}),
	}

	defaultPermSet := NewPermSetByString([]string{"self:XR", "dash:XR", "bills:XL", "reports:R"})

	expectedPermSet := NewPermSetByString([]string{"self:XR", "dash:XR", "bills:XL", "reports:R"})

	// Call the BuildPermSetWithLimit function
	resultPermSet := roleFactory.BuildPermSetWithLimitMulti(role, defaultPermSet)

	// Validate the result against the expected PermSet
	if !reflect.DeepEqual(resultPermSet, expectedPermSet) {
		t.Errorf("BuildPermSetWithLimit failed: expected %v, got %v", expectedPermSet.ToStringArray(), resultPermSet.ToStringArray())
	}
}
