package aclient_redis

import (
	"sync"

	"github.com/alexedwards/scs/redisstore"
)

type RedisStoreMap struct {
	mu    sync.RWMutex
	store map[string]*redisstore.RedisStore
}

// NewRedisStoreMap creates a new RedisStoreMap.
func NewRedisStoreMap() *RedisStoreMap {
	return &RedisStoreMap{
		store: make(map[string]*redisstore.RedisStore),
	}
}

// Get safely retrieves a RedisStore by its prefix.
func (rsm *RedisStoreMap) Get(prefix string) (*redisstore.RedisStore, bool) {
	rsm.mu.RLock()
	defer rsm.mu.RUnlock()
	store, exists := rsm.store[prefix]
	return store, exists
}

// Set safely sets a RedisStore with the given prefix.
func (rsm *RedisStoreMap) Set(prefix string, store *redisstore.RedisStore) {
	rsm.mu.Lock()
	defer rsm.mu.Unlock()
	rsm.store[prefix] = store
}
