package areflect

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// iteststruct is an interface used for testing, requiring a GetType method.
type iteststruct interface {
	GetType() string
}

// testStructA is a struct used for testing, with a type and additional properties.
type testStructA struct {
	Type      string `json:"type"`
	PropertyA string `json:"propertyA"`
	PropertyB string `json:"propertyB"`
}

// GetType returns the type of testStructA.
func (s *testStructA) GetType() string {
	return "structA"
}

// testStructB embeds testStructA and adds an additional property.
type testStructB struct {
	*testStructA
	PropertyC string `json:"propertyC"`
}

// GetType returns the type of testStructB.
func (s *testStructB) GetType() string {
	return "structB"
}

// returnInitPByType is a helper function that returns the reflect.Type based on a type name.
func returnInitPByType(typeName string) (reflect.Type, error) {
	switch typeName {
	case "structA":
		return reflect.TypeOf(testStructA{}), nil
	case "structB":
		return reflect.TypeOf(testStructB{}), nil
	default:
		return nil, nil
	}
}

// TestHasReflectTypeFunction tests the HasReflectTypeFunction method of the type manager.
func TestHasReflectTypeFunction(t *testing.T) {
	tm := TypeManager()
	key := "test"
	fnIdOne := "one"
	fnIdTwo := "two"

	// Ensure the type manager is empty before starting the test.
	assert.NoError(t, tm.Remove(key))
	assert.Equal(t, 0, tm.Count())

	// Register a function and check if it's registered.
	assert.NoError(t, tm.Register(key, fnIdOne, returnInitPByType))
	assert.Equal(t, 1, tm.Count())
	assert.Equal(t, 1, tm.CountByKey(key))
	assert.True(t, tm.HasReflectTypeFunction(key, fnIdOne))

	// Re-register the same function and check the count.
	assert.NoError(t, tm.Register(key, fnIdOne, returnInitPByType))
	assert.Equal(t, 1, tm.CountByKey(key))

	// Register a second function and check the count.
	assert.NoError(t, tm.Register(key, fnIdTwo, returnInitPByType))
	assert.Equal(t, 2, tm.CountByKey(key))
	assert.True(t, tm.HasReflectTypeFunction(key, fnIdTwo))

	// Remove functions by function ID and check the count.
	assert.NoError(t, tm.RemoveByFunctionId(key, fnIdOne))
	assert.Equal(t, 1, tm.CountByKey(key))
	assert.NoError(t, tm.RemoveByFunctionId(key, fnIdTwo))
	assert.Equal(t, 0, tm.CountByKey(key))

	// Remove the key and check the count.
	assert.NoError(t, tm.Remove(key))
	assert.Equal(t, 0, tm.Count())
}

// TestFindReflectType tests the FindReflectType method of the type manager.
func TestFindReflectType(t *testing.T) {
	tm := TypeManager()
	key := "test"

	// Register functions for testing.
	assert.NoError(t, tm.Register(key, "one", returnInitPByType))

	// Find and assert the types.
	rtypeA, err := tm.FindReflectType(key, "structA")
	assert.NoError(t, err)
	assert.Equal(t, "testStructA", rtypeA.Name())

	rtypeB, err := tm.FindReflectType(key, "structB")
	assert.NoError(t, err)
	assert.Equal(t, "testStructB", rtypeB.Name())

	// Attempt to find an unknown type and expect an error.
	_, err = tm.FindReflectType(key, "unknown")
	assert.Error(t, err)

	// Create instances using reflect.New and assert they implement the iteststruct interface.
	instanceA, ok := reflect.New(rtypeA).Interface().(iteststruct)
	assert.True(t, ok)
	assert.Equal(t, "structA", instanceA.GetType())

	instanceB, ok := reflect.New(rtypeB).Interface().(iteststruct)
	assert.True(t, ok)
	assert.Equal(t, "structB", instanceB.GetType())

	// Clean up by removing the registered functions.
	assert.NoError(t, tm.Remove(key))
}

// Plugin is the common interface that all plugins must implement.
type Plugin interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	GetType() string
}

// testPluginA is a concrete implementation of the Plugin interface.
type testPluginA struct {
	Type      string `json:"type"`
	PropertyA string `json:"propertyA"`
}

func (p *testPluginA) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *testPluginA) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p *testPluginA) GetType() string {
	return p.Type
}

// testPluginB is another concrete implementation of the Plugin interface.
type testPluginB struct {
	Type      string `json:"type"`
	PropertyB string `json:"propertyB"`
}

func (p *testPluginB) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *testPluginB) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p *testPluginB) GetType() string {
	return p.Type
}

// registerTestPlugins registers the test plugins with the type manager.
func registerTestPlugins() {
	TypeManager().RemoveAll()
	tm := TypeManager()
	tm.Register("plugins", "testPluginA", func(typeName string) (reflect.Type, error) {
		if typeName == "testPluginA" {
			return reflect.TypeOf(testPluginA{}), nil
		}
		return nil, nil
	})
	tm.Register("plugins", "testPluginB", func(typeName string) (reflect.Type, error) {
		if typeName == "testPluginB" {
			return reflect.TypeOf(testPluginB{}), nil
		}
		return nil, nil
	})
}

// TestPluginMarshalling tests the marshalling and unmarshalling of plugins.
func TestPluginMarshalling(t *testing.T) {
	registerTestPlugins()

	// Create instances of the plugins.
	pluginA := &testPluginA{Type: "testPluginA", PropertyA: "valueA"}
	pluginB := &testPluginB{Type: "testPluginB", PropertyB: "valueB"}

	// Marshal the plugins.
	dataA, err := pluginA.Marshal()
	assert.NoError(t, err)
	dataB, err := pluginB.Marshal()
	assert.NoError(t, err)

	// Unmarshal the plugins using the type manager to find the correct type.
	tm := TypeManager()
	rtypeA, err := tm.FindReflectType("plugins", "testPluginA")
	assert.NoError(t, err)
	rtypeB, err := tm.FindReflectType("plugins", "testPluginB")
	assert.NoError(t, err)

	instanceA := reflect.New(rtypeA).Interface().(Plugin)
	err = instanceA.Unmarshal(dataA)
	assert.NoError(t, err)
	assert.Equal(t, "testPluginA", instanceA.GetType())

	instanceB := reflect.New(rtypeB).Interface().(Plugin)
	err = instanceB.Unmarshal(dataB)
	assert.NoError(t, err)
	assert.Equal(t, "testPluginB", instanceB.GetType())
}

// Define the interface that your plugins will implement.
type IPlugin interface {
	PluginType() string
}

// registerTestPluginAB registers the test plugins with the type manager.
func registerTestPluginAB() {
	TypeManager().RemoveAll()
	tm := TypeManager()
	tm.Register("plugins", "PluginA", func(typeName string) (reflect.Type, error) {
		if typeName == "PluginA" {
			return reflect.TypeOf(PluginA{}), nil
		}
		return nil, nil
	})
	tm.Register("plugins", "PluginB", func(typeName string) (reflect.Type, error) {
		if typeName == "PluginB" {
			return reflect.TypeOf(PluginB{}), nil
		}
		return nil, nil
	})
}

// Concrete implementations of the IPlugin interface.
type PluginA struct {
	Name string `json:"name"`
}

func (p *PluginA) PluginType() string {
	return "PluginA"
}

type PluginB struct {
	Description string `json:"description"`
}

func (p *PluginB) PluginType() string {
	return "PluginB"
}

// MarshalPlugin marshals a Plugin into JSON, including the type information.
func MarshalPlugin(p IPlugin) ([]byte, error) {
	typeName := p.PluginType()
	wrapper := map[string]interface{}{
		"type": typeName,
		"data": p,
	}
	return json.Marshal(wrapper)
}

// UnmarshalPlugin unmarshals JSON into a IPlugin, using the type information to create the correct type.
func UnmarshalPlugin(data []byte) (IPlugin, error) {
	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}

	var typeName string
	if err := json.Unmarshal(wrapper["type"], &typeName); err != nil {
		return nil, err
	}

	tm := TypeManager()
	pluginType, err := tm.FindReflectType("plugins", typeName)
	if err != nil {
		return nil, err
	}

	plugin := reflect.New(pluginType).Interface().(IPlugin)
	if err := json.Unmarshal(wrapper["data"], plugin); err != nil {
		return nil, err
	}

	return plugin, nil
}

// TestMarshalUnmarshalPlugin tests the marshalling and unmarshalling of IPlugins.
func TestMarshalUnmarshalPlugin(t *testing.T) {
	registerTestPluginAB()

	pluginA := &PluginA{Name: "Test Plugin A"}
	pluginB := &PluginB{Description: "Test Plugin B"}

	// Test marshalling and unmarshalling for PluginA.
	dataA, err := MarshalPlugin(pluginA)
	if err != nil {
		t.Fatalf("MarshalPlugin() error = %v", err)
	}

	unmarshalledPluginA, err := UnmarshalPlugin(dataA)
	if err != nil {
		t.Fatalf("UnmarshalPlugin() error = %v", err)
	}

	assertEqualPlugin(t, pluginA, unmarshalledPluginA)

	// Test marshalling and unmarshalling for PluginB.
	dataB, err := MarshalPlugin(pluginB)
	if err != nil {
		t.Fatalf("MarshalPlugin() error = %v", err)
	}

	unmarshalledPluginB, err := UnmarshalPlugin(dataB)
	if err != nil {
		t.Fatalf("UnmarshalPlugin() error = %v", err)
	}

	assertEqualPlugin(t, pluginB, unmarshalledPluginB)
}

// assertEqualPlugin is a helper function to assert that two Plugins are equal.
func assertEqualPlugin(t *testing.T, expected IPlugin, actual IPlugin) {
	registerTestPluginAB()

	if expected.PluginType() != actual.PluginType() {
		t.Errorf("PluginType mismatch: expected %v, got %v", expected.PluginType(), actual.PluginType())
	}

	switch e := expected.(type) {
	case *PluginA:
		a, ok := actual.(*PluginA)
		if !ok || e.Name != a.Name {
			t.Errorf("PluginA mismatch: expected %+v, got %+v", e, a)
		}
	case *PluginB:
		b, ok := actual.(*PluginB)
		if !ok || e.Description != b.Description {
			t.Errorf("PluginB mismatch: expected %+v, got %+v", e, b)
		}
	default:
		t.Errorf("Unknown Plugin type: %v", reflect.TypeOf(expected))
	}
}
