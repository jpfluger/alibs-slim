package asessions

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestRoles_SetRole(t *testing.T) {
	role1 := &Role{Name: RoleName("admin")}
	role2 := &Role{Name: RoleName("user")}
	role3 := &Role{Name: RoleName("admin")} // Update existing role

	roles := Roles{}

	// Add the first role
	roles.SetRole(role1)
	if len(roles) != 1 {
		t.Fatalf("expected 1 role, got %d", len(roles))
	}
	if !reflect.DeepEqual(roles[0], role1) {
		t.Errorf("expected role1 to be added, got %+v", roles[0])
	}

	// Add a second role
	roles.SetRole(role2)
	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	}
	if !reflect.DeepEqual(roles[1], role2) {
		t.Errorf("expected role2 to be added, got %+v", roles[1])
	}

	// Update an existing role (role1)
	roles.SetRole(role3)
	if len(roles) != 2 {
		t.Fatalf("expected 2 roles after update, got %d", len(roles))
	}
	if !reflect.DeepEqual(roles[0], role3) {
		t.Errorf("expected role1 to be updated, got %+v", roles[0])
	}
}

func TestRoles_RemoveRole(t *testing.T) {
	role1 := &Role{Name: RoleName("admin")}
	role2 := &Role{Name: RoleName("user")}

	roles := Roles{role1, role2}

	// Remove the second role
	err := roles.RemoveRole(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 1 {
		t.Fatalf("expected 1 role, got %d", len(roles))
	}
	if roles[0] != role1 {
		t.Errorf("expected role1 to remain, got %+v", roles[0])
	}

	// Attempt to remove an out-of-range index
	err = roles.RemoveRole(5)
	if err == nil {
		t.Errorf("expected error for out-of-range index, got nil")
	}
}

func TestRoles_FindRoleByName(t *testing.T) {
	role1 := &Role{Name: RoleName("admin")}
	role2 := &Role{Name: RoleName("user")}

	roles := Roles{role1, role2}

	// Find an existing role
	foundRole := roles.FindRoleByName("admin")
	if foundRole == nil || foundRole.Name != "admin" {
		t.Errorf("expected to find role 'admin', got %v", foundRole)
	}

	// Try to find a non-existing role
	foundRole = roles.FindRoleByName("guest")
	if foundRole != nil {
		t.Errorf("expected to find no role, got %v", foundRole)
	}
}

func TestRole_Validate(t *testing.T) {
	tests := []struct {
		role        *Role
		expectError bool
		description string
	}{
		{&Role{Name: "admin"}, false, "valid role with name"},
		{&Role{}, true, "role with empty name"},
		{nil, true, "nil role"},
		{&Role{PermsPlus: NewPermSetByString([]string{"self:CRUD"})}, true, "role with permsPlus but empty name"},
	}

	for _, test := range tests {
		err := test.role.Validate()
		if (err != nil) != test.expectError {
			t.Errorf("Validate() for %s failed: expected error %v, got %v", test.description, test.expectError, err)
		}
	}
}

func TestRole_BuildPermSetByFactorWithLimit(t *testing.T) {
	roleFactory := RoleFactory{
		"admin": NewPermSetByString([]string{"self:XLCRUD", "dash:XLCRUD", "bills:XLCRUD"}),
		"user":  NewPermSetByString([]string{"self:XR", "dash:XLR", "bills:"}),
	}

	defaultPermSet := NewPermSetByString([]string{"self:XR", "dash:XR", "bills:XL", "reports:R"})

	tests := []struct {
		role        *Role
		expectError bool
		expected    PermSet
		description string
	}{
		{
			role: &Role{
				Name:       "admin",
				PermsPlus:  NewPermSetByString([]string{"reports:CRUD"}),
				PermsMinus: NewPermSetByString([]string{"dash:L"}),
			},
			expectError: false,
			expected:    NewPermSetByString([]string{"self:XR", "dash:XR", "bills:XL", "reports:R"}),
			description: "valid role with admin permissions constrained by defaultPermSet",
		},
		{
			role:        nil,
			expectError: true,
			description: "nil role",
		},
		{
			role: &Role{
				Name: "nonexistent",
			},
			expectError: false,
			expected:    PermSet{},
			description: "role with no matching names in factory",
		},
		{
			role: &Role{
				PermsPlus: NewPermSetByString([]string{"extra:XR"}),
			},
			expectError: true,
			description: "role with permsPlus but empty name",
		},
	}

	for _, test := range tests {
		result, err := test.role.BuildLimitedPermSet(roleFactory, defaultPermSet)
		if (err != nil) != test.expectError {
			t.Errorf("BuildPermSetByFactorWithLimit() for %s failed: expected error %v, got %v", test.description, test.expectError, err)
		}
		if err == nil && !reflect.DeepEqual(result, test.expected) {
			t.Errorf("BuildPermSetByFactorWithLimit() for %s failed: expected %v, got %v", test.description, test.expected.ToStringArray(), result.ToStringArray())
		}
	}
}

func TestRoleMap(t *testing.T) {
	// Create roles for testing
	role1 := &Role{
		Name:      "admin",
		PermsPlus: NewPermSetByString([]string{"g.self:XRUD"}),
	}
	role2 := &Role{
		Name:      "user",
		PermsPlus: NewPermSetByString([]string{"g.self:XR"}),
	}
	role3 := &Role{
		Name:      "manager",
		PermsPlus: NewPermSetByString([]string{"g.reports:XL"}),
	}

	t.Run("Set and Get Role", func(t *testing.T) {
		rm := RoleMap{}
		rm.Set(role1)

		retrieved := rm.Get("admin")
		assert.NotNil(t, retrieved, "Expected role to be set")
		assert.Equal(t, role1, retrieved, "Expected role to match")
	})

	t.Run("Update Role", func(t *testing.T) {
		rm := RoleMap{}
		rm.Set(role1)

		// Update the same role
		updatedRole := &Role{
			Name:      "admin",
			PermsPlus: NewPermSetByString([]string{"g.self:X"}),
		}
		rm.Set(updatedRole)

		retrieved := rm.Get("admin")
		assert.NotNil(t, retrieved, "Expected updated role to be set")
		assert.Equal(t, updatedRole, retrieved, "Expected role to be updated")
		assert.NotEqual(t, role1, retrieved, "Expected role to be replaced")
	})

	t.Run("Remove Role", func(t *testing.T) {
		rm := RoleMap{}
		rm.Set(role1)
		rm.Remove("admin")

		retrieved := rm.Get("admin")
		assert.Nil(t, retrieved, "Expected role to be removed")
	})

	t.Run("Contains Role", func(t *testing.T) {
		rm := RoleMap{}
		rm.Set(role1)

		assert.True(t, rm.Contains("admin"), "Expected RoleMap to contain 'admin'")
		assert.False(t, rm.Contains("user"), "Expected RoleMap not to contain 'user'")
	})

	t.Run("Merge RoleMap", func(t *testing.T) {
		rm1 := RoleMap{}
		rm1.Set(role1)
		rm1.Set(role2)

		rm2 := RoleMap{}
		rm2.Set(role3)

		rm1.Merge(rm2)

		assert.Equal(t, 3, len(rm1), "Expected RoleMap to have 3 roles after merge")
		assert.NotNil(t, rm1.Get("manager"), "Expected 'manager' role to be added")
		assert.NotNil(t, rm1.Get("admin"), "Expected 'admin' role to remain")
		assert.NotNil(t, rm1.Get("user"), "Expected 'user' role to remain")
	})

	t.Run("Clone RoleMap", func(t *testing.T) {
		rm := RoleMap{}
		rm.Set(role1)
		rm.Set(role2)

		clone := rm.Clone()

		assert.Equal(t, len(rm), len(clone), "Expected cloned RoleMap to have the same number of roles")
		assert.Equal(t, rm["admin"], clone["admin"], "Expected roles to match in cloned RoleMap")
		assert.NotSame(t, rm["admin"], clone["admin"], "Expected roles to be deep copies")
	})
}
