package asessions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPermSet_SetPerm(t *testing.T) {
	ps := PermSet{}

	// Add a new permission
	perm1 := MustNewPerm("admin:CRUD")
	ps.SetPerm(perm1)
	assert.Equal(t, ps["admin"], perm1, "Expected 'admin' permission to be set")

	// Add another permission
	perm2 := MustNewPerm("user:RU")
	ps.SetPerm(perm2)
	assert.Equal(t, ps["user"], perm2, "Expected 'user' permission to be set")

	// Replace an existing permission
	perm3 := MustNewPerm("admin:R")
	ps.SetPerm(perm3)
	assert.Equal(t, ps["admin"], perm3, "Expected 'admin' permission to be replaced with 'R'")
}

// TestPermSet_MergePerm tests merging a permission into a PermSet.
func TestPermSet_MergePerm(t *testing.T) {
	ps := PermSet{}
	perm1 := MustNewPerm("admin:CR")
	perm2 := MustNewPerm("admin:UD")

	ps.SetPerm(perm1)
	ps.MergePerm(perm2)

	assert.Equal(t, 1, len(ps))
	assert.Equal(t, "CRUD", ps["admin"].Value())
}

// TestPermSet_SubtractPerm tests subtracting a permission from a PermSet.
func TestPermSet_SubtractPerm(t *testing.T) {
	ps := PermSet{}
	perm1 := MustNewPerm("admin:CRUD")
	perm2 := MustNewPerm("admin:CR")

	ps.SetPerm(perm1)
	ps.SubtractPerm(perm2)

	assert.Equal(t, 1, len(ps))
	assert.Equal(t, "UD", ps["admin"].Value())
}

// TestPermSet_MatchesPerm tests matching a permission in a PermSet.
func TestPermSet_MatchesPerm(t *testing.T) {
	ps := PermSet{}
	perm := MustNewPerm("admin:CRUD")
	ps.SetPerm(perm)

	matchPerm := MustNewPerm("admin:CR")
	noMatchPerm := MustNewPerm("admin:XL")

	assert.True(t, ps.MatchesPerm(matchPerm))
	assert.False(t, ps.MatchesPerm(noMatchPerm))
}

// TestPermSet_ToStringArray tests converting a PermSet to a string array.
func TestPermSet_ToStringArray(t *testing.T) {
	ps := PermSet{}
	ps.SetPerm(MustNewPerm("admin:CRUD"))
	ps.SetPerm(MustNewPerm("user:CR"))

	arr := ps.ToStringArray()
	assert.ElementsMatch(t, []string{"admin:CRUD", "user:CR"}, arr)
}

// TestPermSet_FromStringArray tests creating a PermSet from a string array.
func TestPermSet_FromStringArray(t *testing.T) {
	perms := []string{"admin:CRUD", "user:CR"}
	ps := FromStringArray(perms)

	assert.Equal(t, 2, len(ps))
	assert.Equal(t, "CRUD", ps["admin"].Value())
	assert.Equal(t, "CR", ps["user"].Value())
}

// TestPermSet_MergeByPermSet tests merging two PermSets.
func TestPermSet_MergeByPermSet(t *testing.T) {
	ps1 := PermSet{}
	ps1.SetPerm(MustNewPerm("admin:CR"))

	ps2 := PermSet{}
	ps2.SetPerm(MustNewPerm("admin:UD"))
	ps2.SetPerm(MustNewPerm("user:CRUD"))

	ps1.MergeByPermSet(ps2)

	assert.Equal(t, 2, len(ps1))
	assert.Equal(t, "CRUD", ps1["admin"].Value())
	assert.Equal(t, "CRUD", ps1["user"].Value())
}

// TestPermSet_SubtractByPermSet tests subtracting one PermSet from another.
func TestPermSet_SubtractByPermSet(t *testing.T) {
	ps1 := PermSet{}
	ps1.SetPerm(MustNewPerm("admin:CRUD"))
	ps1.SetPerm(MustNewPerm("user:CRUD"))

	ps2 := PermSet{}
	ps2.SetPerm(MustNewPerm("admin:CR"))
	ps2.SetPerm(MustNewPerm("user:U"))

	ps1.SubtractByPermSet(ps2)

	assert.Equal(t, 2, len(ps1))
	assert.Equal(t, "UD", ps1["admin"].Value())
	assert.Equal(t, "CRD", ps1["user"].Value())
}

func TestSetPermSetByBits(t *testing.T) {
	ps := PermSet{}

	// Add permissions for "admin"
	ps = SetPermSetByBits(ps, "admin", PERM_C|PERM_R)
	assert.Equal(t, "CR", ps["admin"].Value())

	// Add more permissions for "admin"
	ps = SetPermSetByBits(ps, "admin", PERM_U|PERM_D)
	assert.Equal(t, "CRUD", ps["admin"].Value())

	// Add permissions for "user"
	ps = SetPermSetByBits(ps, "user", PERM_X|PERM_L)
	assert.Equal(t, "XL", ps["user"].Value())

	// Ensure no changes on invalid input
	ps = SetPermSetByBits(ps, "", PERM_C)
	assert.Equal(t, 2, len(ps)) // "admin" and "user" keys remain
	ps = SetPermSetByBits(ps, "user", 0)
	assert.Equal(t, "XL", ps["user"].Value()) // No changes
}

func TestMustNewPermSetByBits(t *testing.T) {
	// Create a PermSet with valid input
	ps := MustNewPermSetByBits("admin", PERM_C|PERM_R)
	assert.Equal(t, 1, len(ps))
	assert.Equal(t, "CR", ps["admin"].Value())

	// Test with empty key
	ps = MustNewPermSetByBits("", PERM_C|PERM_R)
	assert.Equal(t, 0, len(ps)) // Should return an empty PermSet

	// Test with zero bits
	ps = MustNewPermSetByBits("user", 0)
	assert.Equal(t, 0, len(ps)) // Should return an empty PermSet
}

// TestPermSet_MarshalAsInt tests the MarshalAsInt method of PermSet.
func TestPermSet_MarshalAsInt(t *testing.T) {
	ps := PermSet{
		"admin": MustNewPerm("admin:CRUD"),
		"user":  MustNewPerm("user:RU"),
		"guest": nil, // Ensure nil entries are handled correctly
	}

	data, err := ps.MarshalAsInt()
	if err != nil {
		t.Fatalf("unexpected error during MarshalAsInt: %v", err)
	}

	expectedJSON := `[
		"admin:15",
		"user:6"
	]`

	assert.JSONEq(t, expectedJSON, string(data), "MarshalAsInt did not produce expected JSON")
}

func TestPermSet_IsSubsetOf(t *testing.T) {
	tests := []struct {
		current  PermSet
		target   PermSet
		expected bool
		desc     string
	}{
		{
			current: PermSet{
				"admin": MustNewPerm("admin:CR"),
				"user":  MustNewPerm("user:R"),
			},
			target: PermSet{
				"admin": MustNewPerm("admin:CRUD"),
				"user":  MustNewPerm("user:RU"),
			},
			expected: true,
			desc:     "current is a subset of target",
		},
		{
			current: PermSet{
				"admin": MustNewPerm("admin:CRUD"),
				"user":  MustNewPerm("user:R"),
			},
			target: PermSet{
				"admin": MustNewPerm("admin:CR"),
				"user":  MustNewPerm("user:R"),
			},
			expected: false,
			desc:     "admin has excessive permissions",
		},
		{
			current: PermSet{
				"admin": MustNewPerm("admin:CR"),
				"user":  MustNewPerm("user:R"),
			},
			target:   PermSet{},
			expected: false,
			desc:     "target is empty",
		},
		{
			current: PermSet{
				"admin": MustNewPerm("admin:CR"),
			},
			target: PermSet{
				"admin": MustNewPerm("admin:CR"),
				"user":  MustNewPerm("user:RU"),
			},
			expected: true,
			desc:     "admin matches and no extra keys",
		},
		{
			current: PermSet{
				"admin": MustNewPerm("admin:CR"),
				"guest": MustNewPerm("guest:L"),
			},
			target: PermSet{
				"admin": MustNewPerm("admin:CRUD"),
			},
			expected: false,
			desc:     "guest is missing in target",
		},
	}

	for _, test := range tests {
		result := test.current.IsSubsetOf(test.target)
		assert.Equal(t, test.expected, result, test.desc)
	}
}

func TestPermSet_ReplaceExcessivePermSet(t *testing.T) {
	// Create the master PermSet
	master := PermSet{
		"read":  MustNewPermByBitValue("read", PERM_R|PERM_U),
		"write": MustNewPermByBitValue("write", PERM_C),
	}

	// Create the target PermSet that may exceed the master
	ps := PermSet{
		"read":  MustNewPermByBitValue("read", PERM_R|PERM_U|PERM_X), // Exceeds master
		"write": MustNewPermByBitValue("write", PERM_C|PERM_U),       // Exceeds master
		"admin": MustNewPermByBitValue("admin", PERM_X),              // Not in master
	}

	// Replace excessive permissions
	ps.ReplaceExcessivePermSet(master)

	// Expected results
	expected := PermSet{
		"read":  MustNewPermByBitValue("read", PERM_R|PERM_U),
		"write": MustNewPermByBitValue("write", PERM_C),
	}

	// Verify the result
	assert.Equal(t, expected.ToStringArray(), ps.ToStringArray(), "Permissions should match the expected after replacement")
	assert.NotContains(t, ps, "admin", "Permission 'admin' should be removed as it is not in master")
}

func TestPermSet_Clone(t *testing.T) {
	// Create a sample PermSet
	original := PermSet{
		"admin": &Perm{
			key:      "admin",
			value:    &PermValue{value: 15}, // CRUD permissions
			category: "management",
		},
		"user": &Perm{
			key:      "user",
			value:    &PermValue{value: 6}, // RU permissions
			category: "general",
		},
		"guest": nil, // Ensure nil values are handled
	}

	// Clone the PermSet
	cloned := original.Clone()

	// Verify that the cloned PermSet is not the same as the original
	assert.NotSame(t, &original, &cloned, "Cloned PermSet should be a different instance")

	// Verify that the cloned PermSet has fewer keys due to nil-skipping
	expectedLength := len(original) - 1 // One nil key (guest) should be omitted
	assert.Equal(t, expectedLength, len(cloned), "Cloned PermSet should have fewer entries than the original due to nil-skipping")

	// Check that all non-nil entries in the original are present and correctly cloned in the clone
	for key, originalPerm := range original {
		// If the original key has a nil value, it should not exist in the cloned PermSet
		if originalPerm == nil {
			_, exists := cloned[key]
			assert.False(t, exists, "Key '%s' with nil value in original should not exist in the cloned PermSet", key)
			continue
		}

		// Validate the presence of the key in the cloned PermSet
		clonedPerm, exists := cloned[key]
		assert.True(t, exists, "Key '%s' should exist in the cloned PermSet", key)

		// Ensure cloned Perm is not the same instance as the original Perm
		assert.NotSame(t, originalPerm, clonedPerm, "Perm for key '%s' should be a different instance", key)

		// Verify that the fields match
		assert.Equal(t, originalPerm.key, clonedPerm.key, "Key for Perm '%s' should match in the cloned PermSet", key)
		assert.Equal(t, originalPerm.category, clonedPerm.category, "Category for Perm '%s' should match in the cloned PermSet", key)

		// Ensure PermValue is deeply cloned
		if originalPerm.value != nil {
			assert.NotSame(t, originalPerm.value, clonedPerm.value, "PermValue for key '%s' should be a different instance", key)
			assert.Equal(t, originalPerm.value.value, clonedPerm.value.value, "PermValue for key '%s' should match in the cloned PermSet", key)
		} else {
			assert.Nil(t, clonedPerm.value, "PermValue for key '%s' in the cloned PermSet should also be nil", key)
		}
	}

	// Modify the cloned PermSet and ensure it does not affect the original
	cloned["admin"].key = "superadmin"
	cloned["admin"].value.value = 31
	cloned["user"] = nil

	assert.NotEqual(t, original["admin"].key, cloned["admin"].key, "Modifying the cloned PermSet should not affect the original")
	assert.NotEqual(t, original["admin"].value.value, cloned["admin"].value.value, "Modifying the cloned PermSet should not affect the original PermValue")
	assert.NotNil(t, original["user"], "Original PermSet should remain unchanged after modifying the clone")
	assert.Nil(t, cloned["user"], "Cloned PermSet should reflect changes independently of the original")
}
