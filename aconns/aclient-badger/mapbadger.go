// mapbadger.go
package aclient_badger

import (
	"sync"

	"github.com/jpfluger/alibs-slim/aconns"
)

type MapBadger struct {
	mu    sync.RWMutex
	store map[aconns.AdapterName]*AClientBadger
}

// NewMapBadger creates a new MapBadger.
func NewMapBadger() *MapBadger {
	return &MapBadger{
		store: make(map[aconns.AdapterName]*AClientBadger),
	}
}

// Get safely retrieves an AClientBadger by its AdapterName.
func (mb *MapBadger) Get(name aconns.AdapterName) (*AClientBadger, bool) {
	mb.mu.RLock()
	defer mb.mu.RUnlock()
	client, exists := mb.store[name]
	return client, exists
}

// Set safely sets an AClientBadger with the given AdapterName.
func (mb *MapBadger) Set(name aconns.AdapterName, client *AClientBadger) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.store[name] = client
}

// Delete safely removes an AClientBadger by its AdapterName.
func (mb *MapBadger) Delete(name aconns.AdapterName) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	delete(mb.store, name)
}
