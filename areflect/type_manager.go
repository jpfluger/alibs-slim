package areflect

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

/*
The code is part of a reflection-based type management system in Go, which could be used in a plugin architecture. Here's an analysis of its use cases and alternative design patterns:

**Use Cases:**
1. **Dynamic Type Resolution:** The system can dynamically resolve types based on string identifiers, which is useful in scenarios where types need to be determined at runtime, such as loading plugins or components that are not known at compile time.
2. **Serialization/Deserialization:** In systems that serialize and deserialize data, especially when dealing with various data types or schemas, this type manager can help map type names to actual Go types.
3. **Dependency Injection:** The type manager can act as a registry for service implementations, allowing for dynamic dependency injection based on type names.
4. **Extensibility:** It allows an application to be extended with new functionality without recompiling the entire application. New types can be registered at runtime, and existing ones can be overridden or removed.

**Alternative Design Patterns for Plugin Architecture:**
1. **Interface-based Plugins:** Instead of using reflection, plugins implement a known interface. The main application uses these interfaces to interact with the plugins, ensuring compile-time type safety.
2. **Module Pattern:** Using Go's built-in support for modules, plugins can be separate Go modules that the main application imports. This approach also leverages interfaces for interaction.
3. **RPC or IPC:** Plugins run as separate processes and communicate with the main application through Remote Procedure Calls (RPC) or Inter-Process Communication (IPC). This pattern can improve security and stability since plugins are isolated from the main application's memory space.
4. **Middleware Pattern:** In web applications, plugins can be implemented as middleware that hooks into the request/response lifecycle to add or modify functionality.
5. **Event-driven Plugins:** An event bus within the application allows plugins to subscribe to and emit events. Plugins perform actions in response to events without direct coupling to the main application's code.
6. **Scripting Language Plugins:** Embedding a scripting language interpreter (like Lua or JavaScript) into the application allows writing plugins in that scripting language, which the main application executes at runtime.

Each design pattern has its trade-offs. Reflection-based systems like the one in the code are flexible and powerful but can lead to runtime errors if types are not managed carefully. Interface-based and module patterns provide more safety and are easier to reason about but require knowing all possible types at compile time. RPC/IPC and event-driven patterns offer good extensibility and isolation but may introduce complexity in communication and event management. Scripting language plugins offer great flexibility and ease of writing but may have performance implications and require embedding an interpreter.
*/

// ReturnReflectType is a function type that returns a reflect.Type based on a given type name.
type ReturnReflectType func(typeName string) (reflect.Type, error)

// ReturnReflectTypes is a map that associates a string key with a ReturnReflectType function.
type ReturnReflectTypes map[string]ReturnReflectType

// typeManagerFunctions holds a map of ReturnReflectTypes and provides methods to manage them.
type typeManagerFunctions struct {
	returnReflectTypes map[string]ReturnReflectTypes
	mu                 sync.RWMutex
}

// typeManager is a singleton instance of typeManagerFunctions.
var typeManager *typeManagerFunctions

// once is used to initialize the singleton instance only once.
var once sync.Once

// TypeManager returns the singleton instance of typeManagerFunctions, creating it if necessary.
func TypeManager() *typeManagerFunctions {
	once.Do(func() {
		typeManager = &typeManagerFunctions{
			returnReflectTypes: make(map[string]ReturnReflectTypes),
		}
	})
	return typeManager
}

// Get retrieves the ReturnReflectTypes associated with a key.
func (gf *typeManagerFunctions) Get(key string) (ReturnReflectTypes, error) {
	gf.mu.RLock()
	defer gf.mu.RUnlock()

	key = strings.TrimSpace(key)
	if key == "" {
		return nil, fmt.Errorf("key is empty")
	}

	fns, exists := gf.returnReflectTypes[key]
	if !exists {
		fns = ReturnReflectTypes{}
		gf.returnReflectTypes[key] = fns
	}

	return fns, nil
}

// HasReflectTypeFunction checks if a specific function ID is registered under a key.
func (gf *typeManagerFunctions) HasReflectTypeFunction(key string, fnId string) bool {
	key = strings.TrimSpace(key)
	fnId = strings.TrimSpace(fnId)
	if key == "" || fnId == "" {
		return false
	}

	fnMap, err := gf.Get(key)
	if err != nil {
		return false
	}

	fn, ok := fnMap[fnId]
	return ok && fn != nil
}

// Register adds or replaces a ReturnReflectType function under a specific key and function ID.
func (gf *typeManagerFunctions) Register(key string, fnId string, rft ReturnReflectType) error {
	key = strings.TrimSpace(key)
	fnId = strings.TrimSpace(fnId)
	if key == "" {
		return fmt.Errorf("key is empty")
	}
	if fnId == "" {
		return fmt.Errorf("function id is not defined")
	}
	if rft == nil {
		return fmt.Errorf("function parameter is nil")
	}

	gf.mu.Lock()
	defer gf.mu.Unlock()

	fns, exists := gf.returnReflectTypes[key]
	if !exists {
		fns = ReturnReflectTypes{}
		gf.returnReflectTypes[key] = fns
	}
	fns[fnId] = rft

	return nil
}

// FindReflectType searches for a reflect.Type based on a key and type name.
func (gf *typeManagerFunctions) FindReflectType(key string, typeName string) (reflect.Type, error) {
	return gf.FindReflectTypeWithOptions(key, typeName, false)
}

// FindReflectTypeWithOptions searches for a reflect.Type based on a key and type name with an option to suppress the final error.
func (gf *typeManagerFunctions) FindReflectTypeWithOptions(key string, typeName string, noFinalError bool) (reflect.Type, error) {
	fnMap, err := gf.Get(key)
	if err != nil {
		return nil, err
	}

	gf.mu.RLock()
	defer gf.mu.RUnlock()

	for _, fn := range fnMap {
		if fn != nil {
			if reflectType, err := fn(typeName); err != nil {
				return nil, err
			} else if reflectType != nil {
				return reflectType, nil
			}
		}
	}

	if noFinalError {
		return nil, nil
	}

	return nil, fmt.Errorf("failed to find type '%s'", typeName)
}

// Remove deletes all functions registered under a key.
func (gf *typeManagerFunctions) Remove(key string) error {
	gf.mu.Lock()
	defer gf.mu.Unlock()

	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("key is empty")
	}

	delete(gf.returnReflectTypes, key)

	return nil
}

// RemoveAll deletes all registered functions.
func (gf *typeManagerFunctions) RemoveAll() {
	gf.mu.Lock()
	defer gf.mu.Unlock()

	for k := range gf.returnReflectTypes {
		delete(gf.returnReflectTypes, k)
	}
}

// RemoveByFunctionId deletes a specific function registered under a key and function ID.
func (gf *typeManagerFunctions) RemoveByFunctionId(key string, fnId string) error {
	gf.mu.Lock()
	defer gf.mu.Unlock()

	key = strings.TrimSpace(key)
	fnId = strings.TrimSpace(fnId)
	if key == "" {
		return fmt.Errorf("key is empty")
	}
	if fnId == "" {
		return fmt.Errorf("function id is empty")
	}

	fnMap, exists := gf.returnReflectTypes[key]
	if !exists {
		return nil
	}

	delete(fnMap, fnId)

	return nil
}

// Count returns the number of keys registered in the type manager.
func (gf *typeManagerFunctions) Count() int {
	gf.mu.RLock()
	defer gf.mu.RUnlock()

	return len(gf.returnReflectTypes)
}

// CountByKey returns the number of functions registered under a specific key.
func (gf *typeManagerFunctions) CountByKey(key string) int {
	gf.mu.RLock()
	defer gf.mu.RUnlock()

	key = strings.TrimSpace(key)
	if key == "" {
		return -1
	}

	fnMap, exists := gf.returnReflectTypes[key]
	if !exists {
		return 0
	}

	return len(fnMap)
}
