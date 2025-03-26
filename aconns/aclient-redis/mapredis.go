package aclient_redis

import (
	"sync"

	"github.com/jpfluger/alibs-slim/aconns"
)

type MapRedis struct {
	mu    sync.RWMutex
	store map[aconns.AdapterName]*AClientRedis
}

// NewMapRedis creates a new MapRedis.
func NewMapRedis() *MapRedis {
	return &MapRedis{
		store: make(map[aconns.AdapterName]*AClientRedis),
	}
}

// Get safely retrieves an AClientRedis by its AdapterName.
func (mr *MapRedis) Get(name aconns.AdapterName) (*AClientRedis, bool) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()
	client, exists := mr.store[name]
	return client, exists
}

// Set safely sets an AClientRedis with the given AdapterName.
func (mr *MapRedis) Set(name aconns.AdapterName, client *AClientRedis) {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	mr.store[name] = client
}

// Delete safely removes an AClientRedis by its AdapterName.
func (mr *MapRedis) Delete(name aconns.AdapterName) {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	delete(mr.store, name)
}
