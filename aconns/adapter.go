package aconns

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrHostIsEmpty = errors.New("host is empty")

// Adapter satisfies IAdapter.
// Use it to compose connector structs.
type Adapter struct {
	Type AdapterType `json:"type,omitempty"`
	Name AdapterName `json:"name,omitempty"`

	//Network string `json:"network,omitempty"`
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`

	health HealthCheck // Runtime health

	mu sync.RWMutex
}

// GetType returns the AdapterType of the Adapter.
func (a *Adapter) GetType() AdapterType {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Type
}

// GetName returns the AdapterName of the Adapter.
func (a *Adapter) GetName() AdapterName {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Name
}

// GetHost returns the host of the database connection.
func (a *Adapter) GetHost() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Host
}

// GetPort returns the port of the database connection.
func (a *Adapter) GetPort() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Port
}

// Validate checks if the Adapter is valid.
func (a *Adapter) Validate() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.validate()
}

// validate checks if the Adapter is valid.
func (a *Adapter) validate() error {
	a.Type = a.Type.TrimSpace()
	if a.Type.IsEmpty() {
		a.health.Update(HEALTHSTATUS_VALIDATE_FAILED)
		return fmt.Errorf("type is empty")
	}
	a.Name = a.Name.TrimSpace()
	if a.Name.IsEmpty() {
		a.health.Update(HEALTHSTATUS_VALIDATE_FAILED)
		return fmt.Errorf("name is empty")
	}
	a.Host = strings.TrimSpace(a.Host)
	if a.Host == "" {
		a.health.Update(HEALTHSTATUS_VALIDATE_FAILED)
		return ErrHostIsEmpty
	}
	if a.Port < 0 {
		a.Port = 0
	}

	a.health.Update(HEALTHSTATUS_HEALTHY)
	return nil
}

// Test attempts to validate the Adapter and returns the test status and error if any.
// It uses a write lock to ensure thread-safe access.
func (a *Adapter) Test() (bool, TestStatus, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.validate(); err != nil {
		return false, TESTSTATUS_FAILED, err
	}
	// Base has no real test; specific adapters override
	return true, TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

// Refresh refreshes the connection in the Adapter.
// Base implementation is a no-op; specific adapters override with reconnect logic.
func (a *Adapter) Refresh() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	// No-op for base; adapters implement actual refresh
	return nil
}

// GetHealth returns a copy of the Adapter's health check status.
func (a *Adapter) GetHealth() *HealthCheck {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return &a.health
}

// UpdateHealth updates the health status of the Adapter in a thread-safe manner.
func (a *Adapter) UpdateHealth(status HealthStatus) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.health.Update(status)
}

// IsHealthy returns the health status of the HealthCheck.
func (a *Adapter) IsHealthy() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.health.IsHealthy
}
