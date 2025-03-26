package ajson

import (
	"github.com/jpfluger/alibs-slim/autils"
	"os"
	"path"
	"reflect"
	"testing"
)

type TestHelpers struct {
	ValueBool   bool   `json:"vbool,omitempty"`
	ValueString string `json:"vstring,omitempty"`
}

func TestConvertFromMapOfInterfaces2(t *testing.T) {
	dataMap := map[string]interface{}{
		"vbool":   true,
		"vstring": "test-string",
	}

	target := &TestHelpers{}

	if err := ConvertFromMapOfInterfaces(dataMap, target); err != nil {
		t.Error(err)
	}

	if target.ValueBool != dataMap["vbool"] {
		t.Errorf("Expected ValueBool to be %v, got %v", dataMap["vbool"], target.ValueBool)
	}
	if target.ValueString != dataMap["vstring"] {
		t.Errorf("Expected ValueString to be %v, got %v", dataMap["vstring"], target.ValueString)
	}
}

func TestConvertToMapOfInterfaces2(t *testing.T) {
	dataFrom := &TestHelpers{
		ValueBool:   true,
		ValueString: "test-string",
	}

	dataMap, err := ConvertToMapOfInterfaces(dataFrom)
	if err != nil {
		t.Error(err)
	}

	if dataMap["vbool"] != dataFrom.ValueBool {
		t.Errorf("Expected vbool to be %v, got %v", dataFrom.ValueBool, dataMap["vbool"])
	}
	if dataMap["vstring"] != dataFrom.ValueString {
		t.Errorf("Expected vstring to be %v, got %v", dataFrom.ValueString, dataMap["vstring"])
	}
}

// TestUnmarshalFromFile tests the UnmarshalFromFile function.
func TestUnmarshalFromFile(t *testing.T) {
	dir, err := autils.CreateTempDir()
	if err != nil {
		t.Fatalf("cannot create temp directory; %v", err)
	}
	defer os.RemoveAll(dir)

	// Setup
	var testFile = path.Join(dir, "test.json")
	testData := `{"key": "value"}`

	os.WriteFile(testFile, []byte(testData), 0644)
	defer os.Remove(testFile)

	var result map[string]string

	// Test
	err = UnmarshalFromFile(testFile, &result)
	if err != nil {
		t.Errorf("UnmarshalFromFile failed: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("Expected value 'value', got '%s'", result["key"])
	}
}

// TestUnmarshalBytes tests the UnmarshalBytes function.
func TestUnmarshalBytes(t *testing.T) {
	testData := []byte(`{"key": "value"}`)
	var result map[string]string

	// Test
	err := UnmarshalBytes(testData, &result)
	if err != nil {
		t.Errorf("UnmarshalBytes failed: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("Expected value 'value', got '%s'", result["key"])
	}
}

// TestMarshal tests the Marshal function.
func TestMarshal(t *testing.T) {
	testData := map[string]string{"key": "value"}
	expected := `{"key":"value"}`

	// Test
	result, err := Marshal(testData)
	if err != nil {
		t.Errorf("Marshal failed: %v", err)
	}
	if string(result) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// TestMarshalIndent tests the MarshalIndent function.
func TestMarshalIndent(t *testing.T) {
	testData := map[string]string{"key": "value"}
	expected := "{\n  \"key\": \"value\"\n}"

	// Test
	result, err := MarshalIndent(testData)
	if err != nil {
		t.Errorf("MarshalIndent failed: %v", err)
	}
	if string(result) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// TestConvertToMapOfInterfaces tests the ConvertToMapOfInterfaces function.
func TestConvertToMapOfInterfaces(t *testing.T) {
	testData := struct {
		Key string `json:"key"`
	}{"value"}
	expected := map[string]interface{}{"key": "value"}

	// Test
	result, err := ConvertToMapOfInterfaces(testData)
	if err != nil {
		t.Errorf("ConvertToMapOfInterfaces failed: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected '%v', got '%v'", expected, result)
	}
}

// TestConvertFromMapOfInterfaces tests the ConvertFromMapOfInterfaces function.
func TestConvertFromMapOfInterfaces(t *testing.T) {
	testData := map[string]interface{}{"key": "value"}
	expected := struct {
		Key string `json:"key"`
	}{"value"}
	var result struct {
		Key string `json:"key"`
	}

	// Test
	err := ConvertFromMapOfInterfaces(testData, &result)
	if err != nil {
		t.Errorf("ConvertFromMapOfInterfaces failed: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected '%v', got '%v'", expected, result)
	}
}
