package aconns

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/jpfluger/alibs-slim/alog"
	"github.com/jpfluger/alibs-slim/areflect"
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

	TenantInfo *ConnTenantInfo `json:"tenantInfo,omitempty"` // Tenant metadata including region, tenant ID, and priority

	AuthScopes AuthScopes `json:"authScopes,omitempty"` // e.g. ["domain", "module"]
	AuthUsages AuthUsages `json:"authUsages,omitempty"` // e.g. ["mfa", "primary"]

	IsValidated bool // Set to true if Validate() succeeds; used to track config validity without disabling.

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
// Sets IsValidated to true on success.
func (c *Conn) Validate() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.validate()
}

// validate is the internal, lock-free version of Validate.
func (c *Conn) validate() error {
	if c.Ignore {
		return nil
	}
	if c.Adapter == nil {
		return fmt.Errorf("adapter is nil")
	}
	if err := c.Adapter.Validate(); err != nil {
		return fmt.Errorf("failed adapter validate; %v", err)
	}

	// Ensure ID is set
	if c.Id.IsNil() {
		c.Id = NewConnId()
	}

	// Ensure TenantInfo is not nil (backward compatible safety)
	if c.TenantInfo == nil {
		c.TenantInfo = &ConnTenantInfo{}
	}

	c.IsValidated = true

	return nil
}

// Test attempts to establish the connection and handles failures based on IsRequired and IsBootstrap flags.
func (c *Conn) Test() (bool, TestStatus, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Ignore {
		return true, TESTSTATUS_INITIALIZED, nil
	}

	if !c.IsValidated {
		if err := c.validate(); err != nil {
			return false, TESTSTATUS_FAILED, err
		}
	}
	if c.Adapter == nil {
		return false, TESTSTATUS_FAILED, fmt.Errorf("adapter is nil")
	}
	return c.Adapter.Test() // Delegate to adapter
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

	if !c.IsValidated {
		if err := c.validate(); err != nil {
			return false, TESTSTATUS_FAILED, err
		}
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

// Refresh reconnects the underlying adapter.
func (c *Conn) Refresh() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Adapter == nil {
		return fmt.Errorf("no adapter")
	}
	return c.Adapter.Refresh()
}

// GetHealth returns the health from the adapter.
func (c *Conn) GetHealth() *HealthCheck {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.Adapter == nil {
		return NewHealthCheck(HEALTHSTATUS_UNKNOWN)
	}
	return c.Adapter.GetHealth()
}

// StartMonitor periodically checks and refreshes if not validated or unhealthy.
func (c *Conn) StartMonitor(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			<-ticker.C
			c.mu.RLock()
			health := c.GetHealth()
			needsCheck := !c.IsValidated || health.IsStale(interval) || health.IsFailed()
			c.mu.RUnlock()
			if needsCheck {
				if _, _, err := c.Test(); err != nil {
					// Log warn; optional backoff
				} else if !c.IsValidated {
					c.Validate() // Retry validate
				}
			}
		}
	}()
}

// HealthCallback is a function type for notifications on health changes.
// It receives the Conn, updated HealthCheck, and any error from the check.
type HealthCallback func(c *Conn, hc *HealthCheck, err error)

// StartMonitorWithCallback starts a background goroutine to periodically monitor the connection's health.
// It uses the provided interval for checks. If cb is non-nil, it invokes the callback after each check.
// The monitor runs indefinitely until the program exits.
func (c *Conn) StartMonitorWithCallback(interval time.Duration, cb HealthCallback) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Ignore {
		return // No monitoring for ignored conns.
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			c.mu.RLock()
			needsCheck := !c.IsValidated || c.GetHealth().IsStale(interval) || c.GetHealth().IsFailed()
			c.mu.RUnlock()

			if needsCheck {
				_, status, err := c.Test()
				hc := c.GetHealth() // Updated in Test().

				if cb != nil {
					// Invoke callback safely with recovery.
					func() {
						defer func() {
							if r := recover(); r != nil {
								// Log panic (use your logger).
								alog.LOGGER(alog.LOGGER_APP).Error().Msgf("Callback panic recovered: %v", r)
							}
						}()
						cb(c, hc, err)
					}()
				}

				if err != nil {
					// Log warn with status.
					alog.LOGGER(alog.LOGGER_APP).Warn().Err(err).Msgf("Conn test failed with status: %s", status)
					// Optional backoff can be added here.
				} else if !c.IsValidated {
					c.Validate() // Retry validate.
				}
			}
		}
	}()
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

// GetAuthScopes returns the scopes associated with the connection.
// These scopes determine where in the system (e.g., master, domain) this connection is valid.
func (c *Conn) GetAuthScopes() AuthScopes {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.AuthScopes
}

// SetAuthScopes sets the scopes for the connection in a thread-safe way.
func (c *Conn) SetAuthScopes(scopes AuthScopes) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.AuthScopes = scopes
}

// GetTenantInfo returns the tenant info associated with the connection.
func (c *Conn) GetTenantInfo() *ConnTenantInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.TenantInfo == nil {
		// Don't mutate, just return a blank safe value
		return &ConnTenantInfo{}
	}
	return c.TenantInfo
}

// SetTenantInfo sets the tenant info for the connection in a thread-safe way.
func (c *Conn) SetTenantInfo(info *ConnTenantInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.TenantInfo = info
}

// GetAuthUsages returns the list of auth usages this connection supports.
// Usages like "primary", "mfa", or "sspr" define how this connection is invoked in the auth pipeline.
func (c *Conn) GetAuthUsages() AuthUsages {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.AuthUsages
}

// SetAuthUsages updates the list of auth usages for this connection.
func (c *Conn) SetAuthUsages(usages AuthUsages) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.AuthUsages = usages
}

// GetIsValidated returns if this Conn has been validated successfully.
func (c *Conn) GetIsValidated() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.IsValidated
}

// Clone creates a deep copy of the Conn using JSON marshaling/unmarshaling.
func (c *Conn) Clone() *Conn {
	if c == nil {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	data, err := json.Marshal(c)
	if err != nil {
		// Handle error: for simplicity, return nil; in production, log or panic as needed.
		return nil
	}
	clone := &Conn{}
	if err := json.Unmarshal(data, clone); err != nil {
		return nil
	}
	return clone
}

func (c *Conn) CloneWithError() (*Conn, error) {
	if c == nil {
		return nil, fmt.Errorf("nil conn")
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	data, err := json.Marshal(c)
	if err != nil {
		// Handle error: for simplicity, return nil; in production, log or panic as needed.
		return nil, err
	}
	clone := &Conn{}
	if err := json.Unmarshal(data, clone); err != nil {
		return nil, err
	}
	return clone, nil
}
