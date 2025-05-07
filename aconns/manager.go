package aconns

import (
	"fmt"
)

// Manager is used as the basis for interacting with Connections.
// Use cases can vary. In an SaaS application, consider initializing a
// single Manager as a global, while any client-specific conns are
// initialized separately with their own Manager instance. Depending
// upon your app requirements, you may need to copy and modify this
// structure to suit your needs. This struct is designed to be compliant
// for applications that need to run scripts in a sandboxed environment
// with options to limit client access.
type Manager struct {
	// Conns
	Conns IConns `json:"conns,omitempty"`

	// In a sandbox situation, set LimitAccess to true, which indicates
	// to the manager to only allow operations that can be run safely and
	// specifically inside the sandbox.
	LimitAccess bool `json:"limitAccess,omitempty"`
}

// NewManager creates a new Manager instance.
func NewManager() *Manager {
	return &Manager{}
}

// HasConns checks if the Manager has any connections.
func (m *Manager) HasConns() bool {
	return m != nil && len(m.Conns) > 0
}

// Validate checks the syntax of all connections in the queue.
func (m *Manager) Validate() error {
	return m.ValidateWithOptions(false, false)
}

// ValidateBootstrap ensures expected connections can be tried early in a SaaS-style application.
func (m *Manager) ValidateBootstrap() error {
	return m.ValidateWithOptions(true, false)
}

// ValidateRequired ensures required connections can be tried during normal boot.
func (m *Manager) ValidateRequired() error {
	return m.ValidateWithOptions(false, true)
}

// ValidateWithOptions validates connections based on bootstrap and required flags.
func (m *Manager) ValidateWithOptions(checkIsBootstrap bool, checkIsRequired bool) error {
	if m == nil {
		return fmt.Errorf("manager is nil")
	}
	if len(m.Conns) == 0 {
		if checkIsBootstrap || checkIsRequired {
			return fmt.Errorf("manager has no connections")
		}
		return nil // nothing to validate
	}
	var countBootstrap int
	var countRequired int

	if checkIsBootstrap || checkIsRequired {
		var arrToVal []int

		if checkIsBootstrap {
			for ii, conn := range m.Conns {
				if !conn.GetIsBootstrap() {
					continue
				}
				countBootstrap++
				arrToVal = append(arrToVal, ii)
				if err := conn.Validate(); err != nil {
					return fmt.Errorf("failed conns validate checkIsBootstrap at index %d: %v", ii, err)
				}
			}
		}

		if checkIsRequired {
			for ii, conn := range m.Conns {
				if !conn.GetIsRequired() {
					continue
				}
				countRequired++
				for _, val := range arrToVal {
					if val == ii {
						continue // already in array
					}
				}
				arrToVal = append(arrToVal, ii)
				if err := conn.Validate(); err != nil {
					return fmt.Errorf("failed conns validate checkIsRequired at index %d: %v", ii, err)
				}
			}
		}

		if checkIsBootstrap && countBootstrap == 0 {
			return fmt.Errorf("manager has no bootstrap connections")
		}
		if checkIsRequired && countRequired == 0 {
			return fmt.Errorf("manager has no required connections")
		}

		return nil
	}

	for ii, conn := range m.Conns {
		// Validate the connection
		if err := conn.Validate(); err != nil {
			return fmt.Errorf("failed conns validate at index %d: %v", ii, err)
		}
	}

	return nil
}

// FindConn finds an IConn by its UUID.
func (m *Manager) FindConn(id ConnId) IConn {
	if m == nil || m.Conns == nil {
		return nil
	}
	conn, ok := m.Conns.FindByConnId(id)
	if !ok {
		return nil
	}
	return conn
}

// FindAdapter finds an IAdapter by its name.
func (m *Manager) FindAdapter(name AdapterName) IAdapter {
	if m == nil || m.Conns == nil {
		return nil
	}
	adapter, ok := m.Conns.FindByAdapterName(name)
	if !ok {
		return nil
	}
	return adapter
}

// Test both validates and opens a connection, testing it.
// If an adapter uses connection pools, they are initialized here.
func (m *Manager) Test(failQuiet bool) error {
	if m == nil || m.Conns == nil {
		return nil
	}
	for _, conn := range m.Conns {
		if conn.DoIgnore() {
			continue
		}
		_, _, err := conn.GetAdapter().Test()
		if err != nil {
			if !failQuiet {
				return fmt.Errorf("failed test for adapter %s: %v", conn.GetAdapter().GetName(), err)
			}
		}
	}
	return nil
}

// ToTenantManager transforms a Manager into a structured TenantManager by grouping connections by their roles,
// and initializing the authentication pipeline from the auth-capable connections.
func (m *Manager) ToTenantManager() *TenantManager {
	auths := IConns{}
	masters := IConns{}
	tenants := IConns{}

	for _, conn := range m.Conns {
		roles := conn.GetRoles()

		if roles.HasRole(CONNROLE_AUTH) {
			auths = append(auths, conn)
		}
		if roles.HasRole(CONNROLE_MASTER) {
			masters = append(masters, conn)
		}
		if roles.HasRole(CONNROLE_TENANT) {
			tenants = append(tenants, conn)
		}
	}

	tm := &TenantManager{
		Auths:   auths,
		Masters: masters,
		Tenants: tenants,
	}

	// Build the auth pipeline from the Auths list.
	tm.AuthFlows = tm.BuildAuthPipeline()

	return tm
}
