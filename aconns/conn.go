package aconns

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/areflect"
	"reflect"
	"sync"
)

// Conn struct represents a connection with an adapter.
// It encapsulates metadata, role designation, and runtime behavior for a specific integration,
// such as a database, LDAP service, API endpoint, etc.
type Conn struct {
	Ignore      bool     `json:"ignore,omitempty"`      // Ignore flag for the connection.
	Id          ConnId   `json:"id,omitempty"`          // Unique identifier for the connection.
	Description string   `json:"description,omitempty"` // Optional description for the connection.
	Adapter     IAdapter `json:"adapter,omitempty"`     // Adapter associated with the connection.

	// Optional behavior flags
	IsRequired  bool `json:"isRequired,omitempty"`  // If true, then this adapter is required.
	IsBootstrap bool `json:"isBootstrap,omitempty"` // If true, then this adapter is required during boot.

	Roles      ConnRoles      `json:"roles,omitempty"`      // Roles assigned to this connection (e.g., master, auth, tenant)
	TenantInfo ConnTenantInfo `json:"tenantInfo,omitempty"` // Tenant metadata including region, tenant ID, and priority

	// Declares which authentication methods this connection supports (e.g., primary, MFA).
	// These are used to construct TenantManager.AuthFlows, which determines ordered auth execution.
	AuthMethods AuthMethods `json:"authMethods,omitempty"`

	mu sync.RWMutex // Protects access to the fields for concurrent usage.
}

// DoIgnore returns the Ignore flag of the connection.
func (c *Conn) DoIgnore() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Ignore
}

// GetId returns the Id of the connection.
func (c *Conn) GetId() ConnId {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Id
}

// GetAdapter returns the Adapter of the connection.
func (c *Conn) GetAdapter() IAdapter {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Adapter
}

// GetIsRequired returns the IsRequired flag of the connection.
func (c *Conn) GetIsRequired() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.IsRequired
}

// GetIsBootstrap returns the IsBootstrap flag of the connection.
func (c *Conn) GetIsBootstrap() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.IsBootstrap
}

// Validate checks if the Conn is valid.
// It ensures the adapter is not nil and validates the adapter.
// If the Id is not set, it auto-assigns a new UUID.
func (c *Conn) Validate() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Ignore {
		return nil
	}
	if c.Adapter == nil {
		return fmt.Errorf("adapter is nil")
	}
	if err := c.Adapter.Validate(); err != nil {
		return fmt.Errorf("failed adapter validate; %v", err)
	}
	// Auto-set the ID if it is not already set.
	if c.Id.IsNil() {
		c.Id = NewConnId()
	}
	return nil
}

// Test attempts to establish the connection and handles failures based on IsRequired and IsBootstrap flags.
func (c *Conn) Test() (bool, TestStatus, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.test()
}

// TestBootstrap attempts to establish the connection if it is marked as bootstrap.
func (c *Conn) TestBootstrap() (bool, TestStatus, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.IsBootstrap {
		return c.test()
	}
	return true, TESTSTATUS_INITIALIZED, nil
}

// TestRequired attempts to establish the connection if it is marked as required.
func (c *Conn) TestRequired() (bool, TestStatus, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.IsRequired {
		return c.test()
	}
	return true, TESTSTATUS_INITIALIZED, nil
}

// test attempts to establish the connection and handles failures based on IsRequired and IsBootstrap flags.
func (c *Conn) test() (bool, TestStatus, error) {
	if c.Ignore {
		return true, TESTSTATUS_INITIALIZED, nil
	}

	// Attempt to test the adapter
	ok, status, err := c.Adapter.Test()
	if err != nil {
		if c.IsRequired {
			return ok, status, fmt.Errorf("required adapter test failed: %v", err)
		}
		if c.IsBootstrap {
			return ok, status, fmt.Errorf("bootstrap adapter test failed: %v", err)
		}
		return ok, status, fmt.Errorf("adapter test failed: %v", err)
	}

	return ok, status, nil
}

// UnmarshalJSON is a custom unmarshaller for Conn that handles IAdapter.
// It unmarshals the JSON data into the Conn struct and its Adapter.
func (c *Conn) UnmarshalJSON(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	type Alias Conn
	aux := &struct {
		Adapter json.RawMessage `json:"adapter"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("failed to unmarshal Conn: %v", err)
	}

	if aux.Adapter == nil || len(aux.Adapter) == 0 {
		return fmt.Errorf("empty adapter")
	}

	var rawmap map[string]interface{}
	if err := json.Unmarshal(aux.Adapter, &rawmap); err != nil {
		return fmt.Errorf("failed to unmarshal IAdapter: %v", err)
	}

	rawType, ok := rawmap["type"].(string)
	if !ok {
		return fmt.Errorf("type field not found or is not a string in IAdapter")
	}

	rtype, err := areflect.TypeManager().FindReflectType(TYPEMANAGER_CONNADAPTERS, rawType)
	if err != nil {
		return fmt.Errorf("cannot find type struct '%s': %v", rawType, err)
	}

	obj := reflect.New(rtype).Interface()
	if err = json.Unmarshal(aux.Adapter, obj); err != nil {
		return fmt.Errorf("failed to unmarshal IAdapter where type is '%s': %v", rawType, err)
	}

	iAdapter, ok := obj.(IAdapter)
	if !ok {
		return fmt.Errorf("created object does not implement IAdapter where type is '%s'", rawType)
	}
	c.Adapter = iAdapter

	return nil
}

// GetRoles returns the roles associated with the connection.
func (c *Conn) GetRoles() ConnRoles {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Roles
}

// GetTenantInfo returns the tenant info associated with the connection.
func (c *Conn) GetTenantInfo() ConnTenantInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.TenantInfo
}

// SetRoles sets the roles for the connection in a thread-safe way.
func (c *Conn) SetRoles(roles ConnRoles) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Roles = roles
}

// SetTenantInfo sets the tenant info for the connection in a thread-safe way.
func (c *Conn) SetTenantInfo(info ConnTenantInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.TenantInfo = info
}

// GetAuthMethods returns the list of authentication methods supported by this connection.
// These methods (e.g., primary, mfa, sspr) are used to build the ordered authentication pipeline.
func (c *Conn) GetAuthMethods() AuthMethods {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.AuthMethods
}

// SetAuthMethods updates the list of authentication methods for this connection.
// Use this when modifying auth configuration dynamically or during loading.
func (c *Conn) SetAuthMethods(methods AuthMethods) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.AuthMethods = methods
}
