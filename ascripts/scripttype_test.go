package ascripts

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestScriptType has various tests on ScriptType.
func TestScriptType(t *testing.T) {
	// Testing conversion from extension to ScriptType
	assert.Equal(t, SCRIPTTYPE_GO.String(), ExtToScriptType(".go").String(), "Extension to ScriptType conversion failed for '.go'")
	assert.Equal(t, SCRIPTTYPE_GO.String(), ExtToScriptType("go").String(), "Extension to ScriptType conversion failed for 'go'")
	assert.Equal(t, "part", ExtToScriptType(".custom.ext.4.part").String(), "Extension to ScriptType conversion failed for '.custom.ext.4.part'")
	assert.Equal(t, "part", ExtToScriptType("custom.ext.4.part").String(), "Extension to ScriptType conversion failed for 'custom.ext.4.part'")

	// Testing conversion from full file path to ScriptType
	assert.Equal(t, SCRIPTTYPE_GO.String(), FilePathToScriptType("/this/is/a/path/to/file.go").String(), "File path to ScriptType conversion failed for '/this/is/a/path/to/file.go'")
	assert.Equal(t, "part", FilePathToScriptType("/this/is/a/path/to/file.custom.ext.4.part").String(), "File path to ScriptType conversion failed for '/this/is/a/path/to/file.custom.ext.4.part'")

	// Testing conversion from relative file path to ScriptType
	assert.Equal(t, SCRIPTTYPE_GO.String(), FilePathToScriptType("this/is/a/path/to/file.go").String(), "Relative path to ScriptType conversion failed for 'this/is/a/path/to/file.go'")
	assert.Equal(t, "part", FilePathToScriptType("this/is/a/path/to/file.custom.ext.4.part").String(), "Relative path to ScriptType conversion failed for 'this/is/a/path/to/file.custom.ext.4.part'")
	assert.Equal(t, SCRIPTTYPE_GO.String(), FilePathToScriptType("../this/is/a/path/to/file.go").String(), "Relative path to ScriptType conversion failed for '../this/is/a/path/to/file.go'")
	assert.Equal(t, "part", FilePathToScriptType("../this/is/a/path/to/file.custom.ext.4.part").String(), "Relative path to ScriptType conversion failed for '../this/is/a/path/to/file.custom.ext.4.part'")
}

// TestScriptType_IsEmpty checks if IsEmpty method correctly identifies empty ScriptType
func TestScriptType_IsEmpty(t *testing.T) {
	tests := []struct {
		sType   ScriptType
		isEmpty bool
	}{
		{SCRIPTTYPE_GO, false},
		{ScriptType(" "), true},
		{ScriptType(""), true},
	}

	for _, test := range tests {
		if test.sType.IsEmpty() != test.isEmpty {
			t.Errorf("Expected %v for IsEmpty with input '%v'", test.isEmpty, test.sType)
		}
	}
}

// TestScriptType_TrimSpace checks if TrimSpace method correctly trims spaces from ScriptType
func TestScriptType_TrimSpace(t *testing.T) {
	tests := []struct {
		sType      ScriptType
		trimmedStr ScriptType
	}{
		{ScriptType(" go "), SCRIPTTYPE_GO},
		{ScriptType(" html "), SCRIPTTYPE_HTML},
		{ScriptType(" "), ScriptType("")},
	}

	for _, test := range tests {
		if test.sType.TrimSpace() != test.trimmedStr {
			t.Errorf("Expected '%v' for TrimSpace with input '%v'", test.trimmedStr, test.sType)
		}
	}
}

// TestScriptType_String checks if String method correctly converts ScriptType to string
func TestScriptType_String(t *testing.T) {
	if SCRIPTTYPE_GO.String() != "go" {
		t.Errorf("Expected 'go' for String method of SCRIPTTYPE_GO")
	}
}

// TestScriptType_TrimSpaceToLower checks if TrimSpaceToLower method correctly trims spaces and converts to lower case
func TestScriptType_TrimSpaceToLower(t *testing.T) {
	if ScriptType(" Go ").TrimSpaceToLower() != SCRIPTTYPE_GO {
		t.Errorf("Expected 'go' for TrimSpaceToLower with input ' Go '")
	}
}

// TestScriptType_GetExt checks if GetExt method returns correct file extension for ScriptType
func TestScriptType_GetExt(t *testing.T) {
	if SCRIPTTYPE_GO.GetExt() != ".go" {
		t.Errorf("Expected '.go' for GetExt method of SCRIPTTYPE_GO")
	}
}

// TestScriptType_FilePathToScriptType checks if FilePathToScriptType function converts file path to correct ScriptType
func TestScriptType_FilePathToScriptType(t *testing.T) {
	if FilePathToScriptType("test.go") != SCRIPTTYPE_GO {
		t.Errorf("Expected SCRIPTTYPE_GO for FilePathToScriptType with input 'test.go'")
	}
}

// TestScriptType_ExtToScriptType checks if ExtToScriptType function converts file extension to correct ScriptType
func TestScriptType_ExtToScriptType(t *testing.T) {
	if ExtToScriptType("go") != SCRIPTTYPE_GO {
		t.Errorf("Expected SCRIPTTYPE_GO for ExtToScriptType with input 'go'")
	}
}

// TestScriptType_HasMatch checks if HasMatch method correctly identifies if ScriptType is in the slice
func TestScriptType_HasMatch(t *testing.T) {
	sts := ScriptTypes{SCRIPTTYPE_GO, SCRIPTTYPE_HTML}
	if !sts.HasMatch(SCRIPTTYPE_GO) {
		t.Errorf("Expected true for HasMatch with input SCRIPTTYPE_GO")
	}
}
