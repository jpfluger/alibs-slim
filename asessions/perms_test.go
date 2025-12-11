package asessions

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPerm(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		perm, err := NewPerm("admin:CRUD")
		assert.NoError(t, err)
		assert.NotNil(t, perm)
		assert.Equal(t, "admin", perm.Key())
		assert.Equal(t, "CRUD", perm.Value()) // Assuming bitwise conversion to string
	})

	t.Run("InvalidFormat_NoColon", func(t *testing.T) {
		perm, err := NewPerm("nocolon")
		assert.Error(t, err)
		assert.Nil(t, perm)
		assert.Contains(t, err.Error(), "invalid format") // Based on ParsePerm error
	})

	t.Run("InvalidValueChars", func(t *testing.T) {
		perm, err := NewPerm("admin:INVALID")
		assert.Error(t, err)
		assert.Nil(t, perm)
		assert.Contains(t, err.Error(), "invalid value") // Based on ParsePerm error
	})

	t.Run("EmptyInput", func(t *testing.T) {
		perm, err := NewPerm("")
		assert.Error(t, err)
		assert.Nil(t, perm)
		assert.Contains(t, err.Error(), "invalid format") // Based on ParsePerm
	})

	t.Run("NullValue", func(t *testing.T) {
		perm, err := NewPerm("admin:null")
		assert.NoError(t, err) // Assuming ParsePerm handles "null" as empty but valid Perm
		assert.NotNil(t, perm)
		assert.Equal(t, "admin", perm.Key())
		assert.Equal(t, "", perm.Value()) // Or check IsValueEmpty()
	})
}

func TestNewPermSetByString(t *testing.T) {
	t.Run("ValidInputs", func(t *testing.T) {
		ps, err := NewPermSetByString([]string{"admin:CRUD", "user:RU"})
		assert.NoError(t, err)
		assert.NotNil(t, ps)
		assert.Len(t, ps, 2)
		assert.True(t, ps.HasPermS("admin:CRUD"))
		assert.True(t, ps.HasPermS("user:RU"))
	})

	t.Run("MixedValidAndInvalid", func(t *testing.T) {
		ps, err := NewPermSetByString([]string{"good:CRUD", "bad:INVALID"})
		assert.Error(t, err)
		assert.Nil(t, ps)
		assert.Contains(t, err.Error(), "invalid perm string \"bad:INVALID\"")
		assert.Contains(t, err.Error(), "invalid value") // Wrapped from ParsePerm
	})

	t.Run("EmptySlice", func(t *testing.T) {
		ps, err := NewPermSetByString([]string{})
		assert.NoError(t, err)
		assert.NotNil(t, ps)
		assert.Len(t, ps, 0)
	})

	t.Run("AllInvalid", func(t *testing.T) {
		ps, err := NewPermSetByString([]string{"invalid1:FOO", "invalid2:BAR"})
		assert.Error(t, err)
		assert.Nil(t, ps)
		assert.Contains(t, err.Error(), "invalid perm string \"invalid1:FOO\"") // Errors on first invalid
	})

	t.Run("DuplicateKeys", func(t *testing.T) {
		ps, err := NewPermSetByString([]string{"admin:CR", "admin:UD"}) // Assuming SetPerm merges
		assert.NoError(t, err)
		assert.NotNil(t, ps)
		assert.Len(t, ps, 1)
		assert.True(t, ps.HasPermS("admin:CRUD")) // Merged via existing logic
	})
}

// TestMustNewPerm tests the creation of a new Perm with a key-value string.
func TestMustNewPerm(t *testing.T) {
	perm := MustNewPerm("self.test:XCRUD")

	if perm.Key() != "self.test" {
		t.Errorf("perm.Key() expected 'self.test' but has '%s'", perm.Key())
	}

	if perm.Value() != "XCRUD" {
		t.Errorf("perm.Value() expected 'XCRUD' but has '%s'", perm.Value())
	}
}

// TestMustNewPermByPair tests the creation of a new Perm with separate key and value.
func TestMustNewPermByPair(t *testing.T) {
	perm := MustNewPermByPair("self.test", "XCRUD")

	if perm.Key() != "self.test" {
		t.Errorf("perm.Key() expected 'self.test' but has '%s'", perm.Key())
	}

	if perm.Value() != "XCRUD" {
		t.Errorf("perm.Value() expected 'XCRUD' but has '%s'", perm.Value())
	}
}

// TestPermMergeByChar tests the merging of permissions by characters.
func TestPermMergeByChar(t *testing.T) {
	p1 := MustNewPerm("self.test:C")
	p1.MergePermsByChars("R")
	if p1.Value() != "CR" {
		t.Errorf("MergePermsByChars: perm.Value() expected 'CR' but has '%s'", p1.Value())
	}

	p1.MergePermsByChars("UD")
	if p1.Value() != "CRUD" {
		t.Errorf("MergePermsByChars: perm.Value() expected 'CRUD' but has '%s'", p1.Value())
	}

	p1.MergePermsByChars("X")
	if p1.Value() != "XCRUD" {
		t.Errorf("MergePermsByChars: perm.Value() expected 'XCRUD' but has '%s'", p1.Value())
	}
}

// TestPermSubtractByChar tests the subtraction of permissions by characters.
func TestPermSubtractByChar(t *testing.T) {
	p1 := MustNewPerm("self.test:XCRUD")

	p1.SubtractPermsByChars("D")
	if p1.Value() != "XCRU" {
		t.Errorf("SubtractPermsByChars: perm.Value() expected 'XCRU' but has '%s'", p1.Value())
	}

	p1.SubtractPermsByChars("CR")
	if p1.Value() != "XU" {
		t.Errorf("SubtractPermsByChars: perm.Value() expected 'XU' but has '%s'", p1.Value())
	}

	p1.SubtractPermsByChars("XU")
	if p1.Value() != "" {
		t.Errorf("SubtractPermsByChars: perm.Value() expected '' but has '%s'", p1.Value())
	}
}

// TestPermJSON tests the JSON marshaling and unmarshaling of Perm.
func TestPermJSON(t *testing.T) {
	perm := MustNewPerm("admin:XCRUD")

	b, err := json.Marshal(perm)
	if err != nil {
		t.Fatal(err)
	}

	var p2 Perm
	if err := json.Unmarshal(b, &p2); err != nil {
		t.Fatal(err)
	}

	if perm.value.value != p2.value.value {
		t.Errorf("value of perm '%d' does not match p2 '%d'", perm.value.value, p2.value.value)
	}

	if perm.key != p2.key {
		t.Errorf("key of perm '%s' does not match p2 '%s'", perm.key, p2.key)
	}
}

// TestPermJSONSlice tests the JSON marshaling and unmarshaling of a slice of Perms.
func TestPermJSONSlice(t *testing.T) {
	type Test struct {
		IDs  []string `json:"ids,omitempty"`
		Perm Perm     `json:"perm,omitempty"`
	}

	tss := []*Test{
		{IDs: []string{"unique-id-1"}, Perm: *MustNewPerm("admin:XCRUD")},
		{IDs: []string{"unique-id-2"}, Perm: *MustNewPerm("admin:XCRUD")},
	}

	b, err := json.Marshal(tss)
	if err != nil {
		t.Fatal(err)
	}

	var tss2 []*Test
	if err := json.Unmarshal(b, &tss2); err != nil {
		t.Fatal(err)
	}

	if tss[0].IDs[0] != tss2[0].IDs[0] {
		t.Errorf("IDs do not match: '%s' vs '%s'", tss[0].IDs[0], tss2[0].IDs[0])
	}

	if tss[0].Perm.key != tss2[0].Perm.key {
		t.Errorf("Keys do not match: '%s' vs '%s'", tss[0].Perm.key, tss2[0].Perm.key)
	}

	if tss[0].Perm.value.value != tss2[0].Perm.value.value {
		t.Errorf("Values do not match: '%d' vs '%d'", tss[0].Perm.value.value, tss2[0].Perm.value.value)
	}
}

// TestPermsMarshalJSON tests the MarshalJSON method.
func TestPermsMarshalJSON(t *testing.T) {
	perm := MustNewPerm("admin:CRUD")
	jsonBytes, err := perm.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON failed, got error: %s", err)
	}
	if string(jsonBytes) != "\"admin:CRUD\"" {
		t.Errorf("MarshalJSON failed, expected \"admin:CRUD\", got %s", string(jsonBytes))
	}
}

// TestPermsUnmarshalJSON tests the UnmarshalJSON method.
func TestPermsUnmarshalJSON(t *testing.T) {
	perm := &Perm{}
	err := perm.UnmarshalJSON([]byte("\"admin:CRUD\""))
	if err != nil {
		t.Errorf("UnmarshalJSON failed, got error: %s", err)
	}
	if perm.Value() != "CRUD" {
		t.Errorf("UnmarshalJSON failed, expected CRUD, got %s", perm.Value())
	}
}

// TestPermsMustNewPerm tests the creation of a new Perm.
func TestPermsMustNewPerm(t *testing.T) {
	perm := MustNewPerm("admin:XCRUD")
	if perm.Key() != "admin" || perm.Value() != "XCRUD" {
		t.Errorf("MustNewPerm failed, got: %s:%s", perm.Key(), perm.Value())
	}
}

// TestPermsIsValid tests the IsValid method.
func TestPermsIsValid(t *testing.T) {
	perm := MustNewPerm("admin:XCRUD")
	if !perm.IsValid() {
		t.Errorf("IsValid failed, expected true, got false")
	}
}

// TestPermsIsValueEmpty tests the IsValueEmpty method.
func TestPermsIsValueEmpty(t *testing.T) {
	perm := MustNewPerm("admin:")
	if !perm.IsValueEmpty() {
		t.Errorf("IsValueEmpty failed, expected true, got false")
	}
}

// TestPermsCanCreate tests the CanCreate method.
func TestPermsCanCreate(t *testing.T) {
	perm := MustNewPerm("admin:C")
	if !perm.CanCreate() {
		t.Errorf("CanCreate failed, expected true, got false")
	}
}

// TestPermsCanRead tests the CanRead method.
func TestPermsCanRead(t *testing.T) {
	perm := MustNewPerm("admin:R")
	if !perm.CanRead() {
		t.Errorf("CanRead failed, expected true, got false")
	}
}

// TestPermsCanUpdate tests the CanUpdate method.
func TestPermsCanUpdate(t *testing.T) {
	perm := MustNewPerm("admin:U")
	if !perm.CanUpdate() {
		t.Errorf("CanUpdate failed, expected true, got false")
	}
}

// TestPermsCanDelete tests the CanDelete method.
func TestPermsCanDelete(t *testing.T) {
	perm := MustNewPerm("admin:D")
	if !perm.CanDelete() {
		t.Errorf("CanDelete failed, expected true, got false")
	}
}

// TestPermsCanExecute tests the CanExecute method.
func TestPermsCanExecute(t *testing.T) {
	perm := MustNewPerm("admin:X")
	if !perm.CanExecute() {
		t.Errorf("CanExecute failed, expected true, got false")
	}
}

// TestPermsMatchOne tests the MatchOne method.
func TestPermsMatchOne(t *testing.T) {
	perm := MustNewPerm("admin:XCRUD")
	if !perm.MatchOne("X") {
		t.Errorf("MatchOne failed, expected true, got false")
	}
}

// TestPermsMergePermsByChars tests the MergePermsByChars method.
func TestPermsMergePermsByChars(t *testing.T) {
	perm := MustNewPerm("admin:C")
	perm.MergePermsByChars("R")
	if perm.Value() != "CR" {
		t.Errorf("MergePermsByChars failed, expected CR, got %s", perm.Value())
	}
}

// TestPermsSubtractPermsByChars tests the SubtractPermsByChars method.
func TestPermsSubtractPermsByChars(t *testing.T) {
	perm := MustNewPerm("admin:CRUD")
	perm.SubtractPermsByChars("CR")
	if perm.Value() != "UD" {
		t.Errorf("SubtractPermsByChars failed, expected UD, got %s", perm.Value())
	}
}

// TestPermsHasExcessivePermsByChars tests the HasExcessivePermsByChars method.
func TestPermsHasExcessivePermsByChars(t *testing.T) {
	perm := MustNewPerm("admin:CRUD")
	if !perm.HasExcessivePermsByChars("CR") {
		t.Errorf("HasExcessivePermsByChars failed, expected true, got false")
	}
}

// TestPermsReplaceExcessivePermsByChars tests the ReplaceExcessivePermsByChars method.
func TestPermsReplaceExcessivePermsByChars(t *testing.T) {
	perm := MustNewPerm("admin:CRUD")
	perm.ReplaceExcessivePermsByChars("CR")
	if perm.Value() != "CR" {
		t.Errorf("ReplaceExcessivePermsByChars failed, expected CR, got %s", perm.Value())
	}
}

// TestPermsClone tests the Clone method.
func TestPermsClone(t *testing.T) {
	perm := MustNewPerm("admin:CRUD")
	clone := perm.Clone()
	if clone.Value() != perm.Value() {
		t.Errorf("Clone failed, expected %s, got %s", perm.Value(), clone.Value())
	}
}

// Unit tests for SingleAsInt, MarshalJSONAsInt, and UnmarshalJSON
func TestPerm_SingleAsInt(t *testing.T) {
	pv := &PermValue{value: PERM_C | PERM_R | PERM_U} // CRU
	perm := Perm{
		key:   "admin",
		value: pv,
	}
	expected := "admin:7" // Binary: 0111 (CRU)
	result := perm.SingleAsInt()
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestPerm_MarshalJSONAsInt(t *testing.T) {
	pv := &PermValue{value: PERM_C | PERM_R | PERM_U} // CRU
	perm := Perm{
		key:   "admin",
		value: pv,
	}
	data, err := perm.MarshalJSONAsInt()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `"admin:7"`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestPerm_UnmarshalJSON_String(t *testing.T) {
	data := []byte(`"admin:CRU"`)
	perm := &Perm{}
	if err := perm.UnmarshalJSON(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedKey := "admin"
	expectedValue := PERM_C | PERM_R | PERM_U

	if perm.key != expectedKey {
		t.Errorf("expected key %s, got %s", expectedKey, perm.key)
	}
	if perm.value == nil || perm.value.value != expectedValue {
		t.Errorf("expected value %d, got %d", expectedValue, perm.value.value)
	}
}

func TestPerm_UnmarshalJSON_Int(t *testing.T) {
	data := []byte(`"admin:7"`)
	perm := &Perm{}
	if err := perm.UnmarshalJSON(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedKey := "admin"
	expectedValue := PERM_C | PERM_R | PERM_U

	if perm.key != expectedKey {
		t.Errorf("expected key %s, got %s", expectedKey, perm.key)
	}
	if perm.value == nil || perm.value.value != expectedValue {
		t.Errorf("expected value %d, got %d", expectedValue, perm.value.value)
	}
}

func TestPerm_UnmarshalJSON_Invalid(t *testing.T) {
	// Test case with "admin:null"
	data := []byte(`"admin:null"`)
	perm := &Perm{}
	if err := perm.UnmarshalJSON(data); err != nil {
		t.Fatalf("unexpected error for 'admin:null', which gets converted to admin:0")
	}

	// Test case with invalid characters on the right side
	data = []byte(`"admin:!@#"`)
	perm = &Perm{}
	if err := perm.UnmarshalJSON(data); err == nil {
		t.Fatalf("expected error but got none for 'admin:!@#'")
	} else {
		assert.Contains(t, err.Error(), "invalid value")
	}

	// Test case with valid integer value
	data = []byte(`"admin:42"`)
	perm = &Perm{}
	if err := perm.UnmarshalJSON(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, "admin", perm.key)
	assert.Equal(t, 42, perm.value.value)

	// Test case with valid human-readable string
	data = []byte(`"admin:CRUD"`)
	perm = &Perm{}
	if err := perm.UnmarshalJSON(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, "admin", perm.key)
	assert.False(t, perm.value.IsEmptyValue())
}

// TestPermsPermValueSetValues tests setting values for PermValue using bitwise operations.
func TestPermsPermValueSetValues(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	expected := PERM_C | PERM_R | PERM_U | PERM_D

	if pv.value != expected {
		t.Errorf("SetValues failed, expected %b, got %b", expected, pv.value)
	}
}

// TestPermsPermValueIsEmptyValue tests the IsEmptyValue method for PermValue.
func TestPermsPermValueIsEmptyValue(t *testing.T) {
	pv := MustNewPermValue("")
	if !pv.IsEmptyValue() {
		t.Errorf("IsEmptyValue failed, expected true, got false")
	}
}

// TestPermsPermValueHasValue tests the HasValue method for PermValue.
func TestPermsPermValueHasValue(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	if !pv.HasValue() {
		t.Errorf("HasValue failed, expected true, got false")
	}
}

// TestPermsPermValueMatchOne tests the MatchOne method for PermValue.
func TestPermsPermValueMatchOne(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	if !pv.MatchOne("C") {
		t.Errorf("MatchOne failed, expected true, got false")
	}
}

// TestPermsPermValueMatchOneByPerm tests the MatchOneByPerm method for PermValue.
func TestPermsPermValueMatchOneByPerm(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	target := MustNewPermValue("C")
	if !pv.MatchOneByPerm(target) {
		t.Errorf("MatchOneByPerm failed, expected true, got false")
	}
}

// TestPermsPermValueMergePermsByChars tests the MergePermsByChars method for PermValue.
func TestPermsPermValueMergePermsByChars(t *testing.T) {
	pv := MustNewPermValue("CR")
	pv.MergePermsByChars("UD")
	if pv.Values() != "CRUD" {
		t.Errorf("MergePermsByChars failed, expected CRUD, got %s", pv.Values())
	}
}

// TestPermsPermValueHasExcessiveChars tests the HasExcessiveChars method for PermValue.
func TestPermsPermValueHasExcessiveChars(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	if !pv.HasExcessiveChars("CR") {
		t.Errorf("HasExcessiveChars failed, expected true, got false")
	}
}

// TestPermsPermValueReplaceExcessiveChars tests the ReplaceExcessiveChars method for PermValue.
func TestPermsPermValueReplaceExcessiveChars(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	pv.ReplaceExcessiveChars("CR")
	if pv.Values() != "CR" {
		t.Errorf("ReplaceExcessiveChars failed, expected CR, got %s", pv.Values())
	}
}

// TestPermValueSubtractPermsByChars tests the SubtractPermsByChars method for PermValue.
func TestPermValueSubtractPermsByChars(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	pv.SubtractPermsByChars("CR")
	if pv.Values() != "UD" {
		t.Errorf("SubtractPermsByChars failed, expected UD, got %s", pv.Values())
	}
}

// TestPermValueClone tests the Clone method for PermValue.
func TestPermValueClone(t *testing.T) {
	pv := MustNewPermValue("CRUD")
	clone := pv.Clone()
	if clone.Values() != pv.Values() {
		t.Errorf("Clone failed, expected %s, got %s", pv.Values(), clone.Values())
	}
}

// TestPerm_MergePermsByBits tests merging permissions on Perm using bitwise parameters.
func TestPerm_MergePermsByBits(t *testing.T) {
	perm := MustNewPerm("admin:CR")
	perm.MergePermsByBits(PERM_U | PERM_D)
	assert.Equal(t, PERM_C|PERM_R|PERM_U|PERM_D, perm.value.value)
}

// TestPerm_HasExcessiveBits tests checking for excessive permissions on Perm using bitwise parameters.
func TestPerm_HasExcessiveBits(t *testing.T) {
	perm := MustNewPerm("admin:CRUD")
	assert.False(t, perm.HasExcessiveBits(PERM_C|PERM_R|PERM_U|PERM_D|PERM_X|PERM_L))
	assert.True(t, perm.HasExcessiveBits(PERM_C|PERM_R))
}

// TestPerm_ReplaceExcessiveBits tests replacing excessive permissions on Perm using bitwise parameters.
func TestPerm_ReplaceExcessiveBits(t *testing.T) {
	perm := MustNewPerm("admin:CRUD")
	perm.ReplaceExcessiveBits(PERM_C | PERM_R)
	assert.Equal(t, PERM_C|PERM_R, perm.value.value)
}

// TestPerm_SubtractPermsByBits tests subtracting permissions on Perm using bitwise parameters.
func TestPerm_SubtractPermsByBits(t *testing.T) {
	perm := MustNewPerm("admin:CRUD")
	perm.SubtractPermsByBits(PERM_C | PERM_R)
	assert.Equal(t, PERM_U|PERM_D, perm.value.value)
}

// TestPerm_MatchOneByBit tests matching at least one bitwise parameter on Perm.
func TestPerm_MatchOneByBit(t *testing.T) {
	perm := MustNewPerm("admin:CRUD")
	assert.True(t, perm.MatchOneByBit(PERM_U|PERM_X))
	assert.False(t, perm.MatchOneByBit(PERM_L|PERM_X))
}
