package asessions

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestRoleMultis_SetRole(t *testing.T) {
	role1 := &RoleMulti{Names: RoleNames{RoleName("admin")}}
	role2 := &RoleMulti{Names: RoleNames{RoleName("user")}}
	role3 := &RoleMulti{Names: RoleNames{RoleName("admin")}} // Update existing role

	roles := RoleMultis{}

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

func TestRoleMultis_RemoveRole(t *testing.T) {
	role1 := &RoleMulti{Names: RoleNames{RoleName("admin")}}
	role2 := &RoleMulti{Names: RoleNames{RoleName("user")}}

	roles := RoleMultis{role1, role2}

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

func TestRoleMultis_FindRoleByName(t *testing.T) {
	role1 := &RoleMulti{Names: RoleNames{RoleName("admin")}}
	role2 := &RoleMulti{Names: RoleNames{RoleName("user")}}

	roles := RoleMultis{role1, role2}

	// Find an existing role
	foundRole := roles.FindRoleByName("admin")
	if foundRole == nil || !foundRole.Names.Contains("admin") {
		t.Errorf("expected to find role 'admin', got %v", foundRole)
	}

	// Try to find a non-existing role
	foundRole = roles.FindRoleByName("guest")
	if foundRole != nil {
		t.Errorf("expected to find no role, got %v", foundRole)
	}
}

func TestRoleMulti_Validate(t *testing.T) {
	tests := []struct {
		role        *RoleMulti
		expectError bool
		description string
	}{
		{&RoleMulti{Names: RoleNames{"admin"}}, false, "valid role with names"},
		{&RoleMulti{}, true, "role with empty names"},
		{nil, true, "nil role"},
		{&RoleMulti{PermsPlus: MustNewPermSetByString([]string{"self:CRUD"})}, true, "role with permsPlus but empty names"},
	}

	for _, test := range tests {
		err := test.role.Validate()
		if (err != nil) != test.expectError {
			t.Errorf("Validate() for %s failed: expected error %v, got %v", test.description, test.expectError, err)
		}
	}
}

func TestRoleMulti_BuildPermSetByFactorWithLimit(t *testing.T) {
	roleFactory := RoleFactory{
		"admin": MustNewPermSetByString([]string{"self:XLCRUD", "dash:XLCRUD", "bills:XLCRUD"}),
		"user":  MustNewPermSetByString([]string{"self:XR", "dash:XLR", "bills:"}),
	}

	defaultPermSet := MustNewPermSetByString([]string{"self:XR", "dash:XR", "bills:XL", "reports:R"})

	tests := []struct {
		role        *RoleMulti
		expectError bool
		expected    PermSet
		description string
	}{
		{
			role: &RoleMulti{
				Names:      RoleNames{"admin"},
				PermsPlus:  MustNewPermSetByString([]string{"reports:CRUD"}),
				PermsMinus: MustNewPermSetByString([]string{"dash:L"}),
			},
			expectError: false,
			expected:    MustNewPermSetByString([]string{"self:XR", "dash:XR", "bills:XL", "reports:R"}),
			description: "valid role with admin permissions constrained by defaultPermSet",
		},
		{
			role:        nil,
			expectError: true,
			description: "nil role",
		},
		{
			role: &RoleMulti{
				Names: RoleNames{"nonexistent"},
			},
			expectError: false,
			expected:    PermSet{},
			description: "role with no matching names in factory",
		},
		{
			role: &RoleMulti{
				PermsPlus: MustNewPermSetByString([]string{"extra:XR"}),
			},
			expectError: true,
			description: "role with permsPlus but empty names",
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

func TestRoleMultiMap(t *testing.T) {
	roleMultiMap := RoleMultiMap{}

	role1 := &RoleMulti{
		Names:      RoleNames{"admin", "domain-admin"},
		PermsPlus:  MustNewPermSetByString([]string{"self:CRUD"}),
		PermsMinus: MustNewPermSetByString([]string{"reports:R"}),
	}

	role2 := &RoleMulti{
		Names:      RoleNames{"user"},
		PermsPlus:  MustNewPermSetByString([]string{"dash:XR"}),
		PermsMinus: MustNewPermSetByString([]string{"bills:L"}),
	}

	// Test Set and Get
	roleMultiMap.Set(role1)
	roleMultiMap.Set(role2)

	retrievedRole := roleMultiMap.Get("admin")
	assert.NotNil(t, retrievedRole, "Expected to retrieve a role for 'admin'")
	assert.Equal(t, role1.PermsPlus, retrievedRole.PermsPlus, "PermsPlus mismatch")
	assert.Equal(t, role1.PermsMinus, retrievedRole.PermsMinus, "PermsMinus mismatch")

	// Test Exists
	assert.True(t, roleMultiMap.Exists("admin"), "Expected 'admin' to exist")
	assert.False(t, roleMultiMap.Exists("nonexistent"), "Expected 'nonexistent' to not exist")

	// Test Remove
	roleMultiMap.Remove("admin")
	assert.False(t, roleMultiMap.Exists("admin"), "Expected 'admin' to be removed")

	// Test Clone
	clonedMap := roleMultiMap.Clone()
	assert.NotSame(t, &roleMultiMap, &clonedMap, "Cloned map should be a different instance")
	assert.Equal(t, roleMultiMap.Get("user"), clonedMap.Get("user"), "Cloned map content mismatch")
}
