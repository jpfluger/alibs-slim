package asessions

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermValuePermissionsStringToBit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Single permission - X", "X", PERM_X},
		{"Single permission - L", "L", PERM_L},
		{"Single permission - C", "C", PERM_C},
		{"Single permission - R", "R", PERM_R},
		{"Single permission - U", "U", PERM_U},
		{"Single permission - D", "D", PERM_D},
		{"Multiple permissions - XCR", "XCR", PERM_X | PERM_C | PERM_R},
		{"Mixed case input", "xCr", PERM_X | PERM_C | PERM_R},
		{"Invalid characters ignored", "XZY", PERM_X},
		{"Empty string", "", 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := PermissionsStringToBit(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestPermValuePermissionsBitToString(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"Single permission - X", PERM_X, "X"},
		{"Single permission - L", PERM_L, "L"},
		{"Single permission - C", PERM_C, "C"},
		{"Single permission - R", PERM_R, "R"},
		{"Single permission - U", PERM_U, "U"},
		{"Single permission - D", PERM_D, "D"},
		{"Multiple permissions - XCR", PERM_X | PERM_C | PERM_R, "XCR"},
		{"No permissions", 0, ""},
		{"Unordered permissions", PERM_R | PERM_X | PERM_C, "XCR"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := PermissionsBitToString(test.input)
			// Ensure the result is a permutation of the expected string
			assert.ElementsMatch(t, strings.Split(test.expected, ""), strings.Split(result, ""))
		})
	}
}

func TestPermValueRoundTripConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Single permission", "X", "X"},
		{"Multiple permissions", "XCR", "XCR"},
		{"Mixed case input", "xCr", "XCR"},
		{"Invalid characters ignored", "XZY", "X"},
		{"Empty string", "", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bitValue := PermissionsStringToBit(test.input)
			result := PermissionsBitToString(bitValue)
			// Ensure the result is a permutation of the expected string
			assert.ElementsMatch(t, strings.Split(test.expected, ""), strings.Split(result, ""))
		})
	}
}

// TestMustNewPermValue checks the creation of a new PermValue with the provided permissions.
func TestMustNewPermValue(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	assert.True(t, pv.value&PERM_C != 0)
	assert.True(t, pv.value&PERM_R != 0)
	assert.True(t, pv.value&PERM_U != 0)
	assert.True(t, pv.value&PERM_D != 0)
	assert.False(t, pv.value&PERM_X != 0)
	assert.False(t, pv.value&PERM_L != 0)
}

// TestValues verifies the string representation of the permissions in PermValue.
func TestValues(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	assert.Equal(t, "CRUD", pv.Values())
}

// TestSetValues confirms that permissions can be set correctly in PermValue.
func TestSetValues(t *testing.T) {
	pv := &PermValue{}
	pv.SetValues("CRUD")
	assert.True(t, pv.value&PERM_C != 0)
	assert.True(t, pv.value&PERM_R != 0)
	assert.True(t, pv.value&PERM_U != 0)
	assert.True(t, pv.value&PERM_D != 0)
	assert.False(t, pv.value&PERM_X != 0)
	assert.False(t, pv.value&PERM_L != 0)
}

// TestIsEmptyValue checks if PermValue correctly identifies an empty set of permissions.
func TestIsEmptyValue(t *testing.T) {
	pv := &PermValue{}
	assert.True(t, pv.IsEmptyValue())
	pv = MustNewPermValue("CRUD")
	assert.False(t, pv.IsEmptyValue())
}

// TestHasValue verifies that PermValue correctly identifies when it has at least one permission set.
func TestHasValue(t *testing.T) {
	pv := &PermValue{}
	assert.False(t, pv.HasValue())
	pv = MustNewPermValue("CRUD")
	assert.True(t, pv.HasValue())
}

// TestMatchOne checks if PermValue correctly identifies a match with at least one permission character.
func TestMatchOne(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	assert.True(t, pv.MatchOne("C"))
	assert.False(t, pv.MatchOne("L"))
}

// TestMatchOneByPerm verifies that PermValue correctly identifies a match with at least one permission from another PermValue.
func TestMatchOneByPerm(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	target := MustNewPermValue("CU")
	assert.True(t, pv.MatchOneByPerm(target))
	target = MustNewPermValue("LX")
	assert.False(t, pv.MatchOneByPerm(target))
}

// TestMergePermsByChars confirms that permissions are correctly merged into PermValue.
func TestMergePermsByChars(t *testing.T) {
	pv := MustNewPermValue("CR")
	pv.MergePermsByChars("UD")
	assert.True(t, pv.value&PERM_U != 0)
	assert.True(t, pv.value&PERM_D != 0)
}

// TestHasExcessiveChars checks if PermValue correctly identifies excessive permissions.
func TestHasExcessiveChars(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	assert.False(t, pv.HasExcessiveChars("CRUDLX"))
	assert.True(t, pv.HasExcessiveChars("CR"))
}

// TestReplaceExcessiveChars verifies that PermValue correctly replaces excessive permissions.
func TestReplaceExcessiveChars(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	pv.ReplaceExcessiveChars("CR")
	assert.False(t, pv.value&PERM_U != 0)
	assert.False(t, pv.value&PERM_D != 0)
}

// TestSubtractPermsByChars confirms that permissions are correctly subtracted from PermValue.
func TestSubtractPermsByChars(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	pv.SubtractPermsByChars("CR")
	assert.False(t, pv.value&PERM_C != 0)
	assert.False(t, pv.value&PERM_R != 0)
	assert.True(t, pv.value&PERM_U != 0)
	assert.True(t, pv.value&PERM_D != 0)
}

// TestClone checks if PermValue can be cloned correctly.
func TestClone(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	clone := pv.Clone()
	assert.Equal(t, pv.Values(), clone.Values())
}

// TestPermValue performs a series of assertions on PermValue to verify its functionality.
func TestPermValue(t *testing.T) {
	pv := MustNewPermValue("XLCRUD")
	assert.Equal(t, "XLCRUD", pv.Values())

	pv.SetValues("X")
	assert.Equal(t, "X", pv.Values())
	assert.True(t, pv.HasValue())
	assert.True(t, pv.MatchOne("X"))

	target := MustNewPermValue("X")
	assert.True(t, pv.MatchOneByPerm(target))

	pv.MergePermsByChars("LCRUD")
	assert.Equal(t, "XLCRUD", pv.Values())

	pv.SubtractPermsByChars("LCRUD")
	assert.Equal(t, "X", pv.Values())

	pv.SubtractPermsByChars("X")
	assert.Empty(t, pv.Values())

	assert.False(t, pv.HasExcessiveChars("X"))

	pv.ReplaceExcessiveChars("X")
	assert.Empty(t, pv.Values())

	clone := pv.Clone()
	assert.Equal(t, pv.Values(), clone.Values())
}

func TestMarshalJSON(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	jsonData, err := json.Marshal(pv)
	assert.NoError(t, err)
	assert.JSONEq(t, `"CRUD"`, string(jsonData))

	pv = MustNewPermValue("XLCRUD")
	jsonData, err = json.Marshal(pv)
	assert.NoError(t, err)
	assert.JSONEq(t, `"XLCRUD"`, string(jsonData))
}

func TestUnmarshalJSON(t *testing.T) {
	var pv PermValue

	err := json.Unmarshal([]byte(`"CRUD"`), &pv)
	assert.NoError(t, err)
	assert.Equal(t, "CRUD", pv.Values())

	err = json.Unmarshal([]byte(`"XLCRUD"`), &pv)
	assert.NoError(t, err)
	assert.Equal(t, "XLCRUD", pv.Values())

	err = json.Unmarshal([]byte(`""`), &pv)
	assert.NoError(t, err)
	assert.True(t, pv.IsEmptyValue())
}

// Unit tests for MarshalJSONAsInt and UnmarshalJSON
func TestPermValue_MarshalJSONAsInt(t *testing.T) {
	pv := &PermValue{value: PERM_C | PERM_R | PERM_U} // CRU
	data, err := pv.MarshalJSONAsInt()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result int
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal data back to int: %v", err)
	}

	if result != pv.value {
		t.Errorf("expected %d, got %d", pv.value, result)
	}
}

func TestPermValue_UnmarshalJSON_String(t *testing.T) {
	data := []byte(`"CRU"`)
	pv := &PermValue{}
	if err := pv.UnmarshalJSON(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := PERM_C | PERM_R | PERM_U
	if pv.value != expected {
		t.Errorf("expected %d, got %d", expected, pv.value)
	}
}

func TestPermValue_UnmarshalJSON_Int(t *testing.T) {
	data := []byte(`7`) // Binary: 0111 (CRU)
	pv := &PermValue{}
	if err := pv.UnmarshalJSON(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := PERM_C | PERM_R | PERM_U
	if pv.value != expected {
		t.Errorf("expected %d, got %d", expected, pv.value)
	}
}

func TestPermValue_UnmarshalJSON_Invalid(t *testing.T) {
	data := []byte(`true`)
	pv := &PermValue{}
	if err := pv.UnmarshalJSON(data); err == nil {
		t.Fatalf("expected error but got none")
	}
}

// TestPermValue_MergePermsByBits tests merging permissions using bitwise parameters.
func TestPermValue_MergePermsByBits(t *testing.T) {
	pv := MustNewPermValue("CR")
	pv.MergePermsByBits(PERM_U | PERM_D)
	assert.Equal(t, PERM_C|PERM_R|PERM_U|PERM_D, pv.value)
}

// TestPermValue_HasExcessiveBits tests checking for excessive permissions using bitwise parameters.
func TestPermValue_HasExcessiveBits(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	assert.False(t, pv.HasExcessiveBits(PERM_C|PERM_R|PERM_U|PERM_D|PERM_X|PERM_L))
	assert.True(t, pv.HasExcessiveBits(PERM_C|PERM_R))
}

// TestPermValue_ReplaceExcessiveBits tests replacing excessive permissions using bitwise parameters.
func TestPermValue_ReplaceExcessiveBits(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	pv.ReplaceExcessiveBits(PERM_C | PERM_R)
	assert.Equal(t, PERM_C|PERM_R, pv.value)
}

// TestPermValue_SubtractPermsByBits tests subtracting permissions using bitwise parameters.
func TestPermValue_SubtractPermsByBits(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	pv.SubtractPermsByBits(PERM_C | PERM_R)
	assert.Equal(t, PERM_U|PERM_D, pv.value)
}

// TestPermValue_MatchOneByBit tests matching at least one bitwise parameter.
func TestPermValue_MatchOneByBit(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	assert.True(t, pv.MatchOneByBit(PERM_U|PERM_X))
	assert.False(t, pv.MatchOneByBit(PERM_L|PERM_X))
}
