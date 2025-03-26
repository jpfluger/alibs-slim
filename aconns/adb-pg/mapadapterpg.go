package adb_pg

import (
	"github.com/jpfluger/alibs-slim/aconns"
	"sync"
)

// MapAdapterPG is a thread-safe map for storing ADBPG instances by AdapterName.
type MapAdapterPG struct {
	mu    sync.RWMutex
	store map[aconns.AdapterName]*ADBPG
}

// NewMapAdapterPG creates a new MapAdapterPG.
func NewMapAdapterPG() *MapAdapterPG {
	return &MapAdapterPG{
		store: make(map[aconns.AdapterName]*ADBPG),
	}
}

// Get safely retrieves an ADBPG instance by its AdapterName.
func (mpg *MapAdapterPG) Get(name aconns.AdapterName) (*ADBPG, bool) {
	mpg.mu.RLock()
	defer mpg.mu.RUnlock()
	client, exists := mpg.store[name]
	return client, exists
}

// Set safely sets an ADBPG instance with the given AdapterName.
func (mpg *MapAdapterPG) Set(name aconns.AdapterName, client *ADBPG) {
	mpg.mu.Lock()
	defer mpg.mu.Unlock()
	mpg.store[name] = client
}

// Delete safely removes an ADBPG instance by its AdapterName.
func (mpg *MapAdapterPG) Delete(name aconns.AdapterName) {
	mpg.mu.Lock()
	defer mpg.mu.Unlock()
	delete(mpg.store, name)
}
