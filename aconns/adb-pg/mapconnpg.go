package adb_pg

import (
	"github.com/jpfluger/alibs-slim/aconns"
	"sync"
)

// MapConnPG is a thread-safe map for storing ADBPG instances by ConnId.
type MapConnPG struct {
	mu    sync.RWMutex
	store map[aconns.ConnId]*ADBPG
}

// NewMapConnPG creates a new MapConnPG.
func NewMapConnPG() *MapConnPG {
	return &MapConnPG{
		store: make(map[aconns.ConnId]*ADBPG),
	}
}

// Get safely retrieves an ADBPG instance by its AdapterName.
func (mpg *MapConnPG) Get(id aconns.ConnId) (*ADBPG, bool) {
	mpg.mu.RLock()
	defer mpg.mu.RUnlock()
	if id.IsNil() {
		return nil, false
	}
	client, exists := mpg.store[id]
	return client, exists
}

// Set safely sets an ADBPG instance with the given AdapterName.
func (mpg *MapConnPG) Set(id aconns.ConnId, client *ADBPG) {
	mpg.mu.Lock()
	defer mpg.mu.Unlock()
	if id.IsNil() {
		return
	}
	mpg.store[id] = client
}

// Delete safely removes an ADBPG instance by its AdapterName.
func (mpg *MapConnPG) Delete(id aconns.ConnId) {
	mpg.mu.Lock()
	defer mpg.mu.Unlock()
	if id.IsNil() {
		return
	}
	delete(mpg.store, id)
}
