package aclient_badger

import (
	"sync"

	"github.com/jpfluger/alibs-slim/aconns/aclient-badger/badgerstore"
)

type BadgerStoreMap struct {
	mu    sync.RWMutex
	store map[string]*badgerstore.BadgerStore
}

// NewBadgerStoreMap creates a new BadgerStoreMap.
func NewBadgerStoreMap() *BadgerStoreMap {
	return &BadgerStoreMap{
		store: make(map[string]*badgerstore.BadgerStore),
	}
}

// Get safely retrieves a BadgerStore by its prefix.
func (bsm *BadgerStoreMap) Get(prefix string) (*badgerstore.BadgerStore, bool) {
	bsm.mu.RLock()
	defer bsm.mu.RUnlock()
	store, exists := bsm.store[prefix]
	return store, exists
}

// Set safely sets a BadgerStore with the given prefix.
func (bsm *BadgerStoreMap) Set(prefix string, store *badgerstore.BadgerStore) {
	bsm.mu.Lock()
	defer bsm.mu.Unlock()
	bsm.store[prefix] = store
}
