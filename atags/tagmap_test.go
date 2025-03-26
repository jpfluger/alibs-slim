package atags

import (
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestFind verifies the Find method of TagMap.
func TestFind(t *testing.T) {
	tmap := TagMap{
		"type:id": "value",
	}
	if tmap.Find("type:id") != "value" {
		t.Errorf("Find method did not return the correct value")
	}
}

// TestHasValue verifies the HasValue method of TagMap.
func TestHasValue(t *testing.T) {
	tmap := TagMap{
		"type:id": "value",
	}
	if !tmap.HasValue("type:id") {
		t.Errorf("HasValue method should return true for existing key")
	}
}

// TestAdd verifies the Add method of TagMap.
func TestAdd(t *testing.T) {
	tmap := make(TagMap)
	tmap.Add("type:id", "value")
	if tmap["type:id"] != "value" {
		t.Errorf("Add method did not insert the value correctly")
	}
}

// TestRemove verifies the Remove method of TagMap.
func TestRemove(t *testing.T) {
	tmap := TagMap{
		"type:id": "value",
	}
	tmap.Remove("type:id")
	if _, exists := tmap["type:id"]; exists {
		t.Errorf("Remove method did not delete the key")
	}
}

// TestGetValue verifies the GetValue method of TagMap.
func TestGetValue(t *testing.T) {
	tmap := TagMap{
		"type:id": "value",
	}
	val, err := tmap.GetValue("type:id")
	if err != nil || val != "value" {
		t.Errorf("GetValue method did not return the correct value")
	}
}

// TestGetValueAsString verifies the GetValueAsString method of TagMap.
func TestGetValueAsString(t *testing.T) {
	tmap := TagMap{
		"type:id": "value",
	}
	val, err := tmap.GetValueAsString("type:id")
	if err != nil || val != "value" {
		t.Errorf("GetValueAsString method did not return the correct value")
	}
}

// TestGetValueAsInt verifies the GetValueAsInt method of TagMap.
func TestGetValueAsInt(t *testing.T) {
	tmap := TagMap{
		"type:id": 42,
	}
	val, err := tmap.GetValueAsInt("type:id")
	if err != nil || val != 42 {
		t.Errorf("GetValueAsInt method did not return the correct value")
	}
}

// TestGetValueAsFloat verifies the GetValueAsFloat method of TagMap.
func TestGetValueAsFloat(t *testing.T) {
	tmap := TagMap{
		"type:id": 3.14,
	}
	val, err := tmap.GetValueAsFloat("type:id")
	if err != nil || val != 3.14 {
		t.Errorf("GetValueAsFloat method did not return the correct value")
	}
}

// TestGetValueAsBool verifies the GetValueAsBool method of TagMap.
func TestGetValueAsBool(t *testing.T) {
	tmap := TagMap{
		"type:id": true,
	}
	val, err := tmap.GetValueAsBool("type:id")
	if err != nil || !val {
		t.Errorf("GetValueAsBool method did not return the correct value")
	}
}

// TestGetValueAsUUID verifies the GetValueAsUUID method of TagMap.
func TestGetValueAsUUID(t *testing.T) {
	testUUID, _ := uuid.NewV4()
	tmap := TagMap{
		"type:id": testUUID.String(),
	}
	val, err := tmap.GetValueAsUUID("type:id")
	if err != nil || val != testUUID {
		t.Errorf("GetValueAsUUID method did not return the correct value")
	}
}

// TestToArray verifies the ToArray method of TagMap.
func TestToArray(t *testing.T) {
	tmap := TagMap{
		"type:id": "value",
	}
	tarr := tmap.ToArray()
	if len(tarr) != 1 || tarr[0].Key != "type:id" || tarr[0].Value != "value" {
		t.Errorf("ToArray method did not return the correct array")
	}
}

// TestMergeFrom verifies the MergeFrom method of TagMap.
func TestMergeFrom(t *testing.T) {
	source := TagMap{
		"type:id": "sourceValue",
	}
	target := TagMap{
		"type:id": "targetValue",
	}
	merged := source.MergeFrom(target)
	if merged["type:id"] != "targetValue" {
		t.Errorf("MergeFrom method did not merge the maps correctly")
	}
}

// testStruct is a struct used for testing purposes.
type testStruct struct {
	PropertyA string
}

// TestTagMap tests various functionalities of the TagMap.
func TestTagMap(t *testing.T) {
	tm := TagMap{}
	tm.Add("vInt", 100)

	// Test conversion to string.
	vString, err := tm.GetValueAsString("vInt")
	assert.NoError(t, err)
	assert.Equal(t, "100", vString)

	// Test conversion to int.
	vInt, err := tm.GetValueAsInt("vInt")
	assert.NoError(t, err)
	assert.Equal(t, 100, vInt)

	// Test conversion to float.
	vFloat, err := tm.GetValueAsFloat("vInt")
	assert.NoError(t, err)
	assert.Equal(t, float64(100), vFloat)

	// Test conversion to bool (true because != 0).
	vBool, err := tm.GetValueAsBool("vInt")
	assert.NoError(t, err)
	assert.True(t, vBool)

	// Test retrieval of a pointer to a struct.
	tm.Add("vPStruct", &testStruct{PropertyA: "ValueA"})
	v, err := tm.GetValue("vPStruct")
	assert.NoError(t, err)
	vPStruct, ok := v.(*testStruct)
	assert.True(t, ok)
	assert.Equal(t, "ValueA", vPStruct.PropertyA)

	// Test that the value is not a non-pointer struct.
	_, ok = v.(testStruct)
	assert.False(t, ok)

	// Test retrieval of a struct.
	tm.Add("vStruct", testStruct{PropertyA: "ValueA"})
	v, err = tm.GetValue("vStruct")
	assert.NoError(t, err)

	// Test that the value is not a pointer to a struct.
	vPStruct, ok = v.(*testStruct)
	assert.False(t, ok)

	// Test that the value is a non-pointer struct.
	vStruct, ok := v.(testStruct)
	assert.True(t, ok)
	assert.Equal(t, "ValueA", vStruct.PropertyA)
}
