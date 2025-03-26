package aclient_ldap

import (
	"sync"

	"github.com/jpfluger/alibs-slim/aconns"
)

// MapLDAP is a thread-safe map for storing AClientLDAP instances by AdapterName.
type MapLDAP struct {
	mu    sync.RWMutex
	store map[aconns.AdapterName]*AClientLDAP
}

// NewMapLDAP creates a new MapLDAP.
func NewMapLDAP() *MapLDAP {
	return &MapLDAP{
		store: make(map[aconns.AdapterName]*AClientLDAP),
	}
}

// Get safely retrieves an AClientLDAP instance by its AdapterName.
func (ml *MapLDAP) Get(name aconns.AdapterName) (*AClientLDAP, bool) {
	ml.mu.RLock()
	defer ml.mu.RUnlock()
	client, exists := ml.store[name]
	return client, exists
}

// Set safely sets an AClientLDAP instance with the given AdapterName.
func (ml *MapLDAP) Set(name aconns.AdapterName, client *AClientLDAP) {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	ml.store[name] = client
}

// Delete safely removes an AClientLDAP instance by its AdapterName.
func (ml *MapLDAP) Delete(name aconns.AdapterName) {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	delete(ml.store, name)
}
