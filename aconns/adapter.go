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

	//Network  string `json:"network,omitempty"`
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`

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
		return fmt.Errorf("type is empty")
	}
	a.Name = a.Name.TrimSpace()
	if a.Name.IsEmpty() {
		return fmt.Errorf("name is empty")
	}
	a.Host = strings.TrimSpace(a.Host)
	if a.Host == "" {
		return ErrHostIsEmpty
	}
	if a.Port < 0 {
		a.Port = 0
	}

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
	return false, TESTSTATUS_FAILED, fmt.Errorf("not implemented")
}
