package aclient_redis

import (
	"fmt"
	"github.com/alexedwards/scs/redisstore"
	"github.com/jpfluger/alibs-slim/aconns"
	"sync"
)

const (
	REDISSTORE_SCS         = "scs:"
	REDISSTORE_FORGOTLOGIN = "forgotlogin:"

	REDIS_MASTER aconns.AdapterName = "redis:master"
)

// connMap is the global connection map.
var connMap *connMapGlobal
var muCM sync.RWMutex

// connMapGlobal holds the global map of Redis connections.
type connMapGlobal struct {
	Map *MapRedis
	mu  sync.RWMutex
}

func init() {
	connMap = &connMapGlobal{Map: NewMapRedis()}
}

// REDIS returns the global connection map.
func REDIS() *connMapGlobal {
	muCM.RLock()
	defer muCM.RUnlock()
	return connMap
}

// Get retrieves an AClientRedis by its AdapterName.
func (cg *connMapGlobal) Get(name aconns.AdapterName) *AClientRedis {
	if name.IsEmpty() {
		return nil
	}

	cg.mu.RLock()
	defer cg.mu.RUnlock()

	client, exists := cg.Map.Get(name)
	if !exists {
		return nil
	}
	return client
}

// Set adds or updates an AClientRedis in the connection map.
func (cg *connMapGlobal) Set(cn *AClientRedis) error {
	if cn == nil {
		return fmt.Errorf("connMapGlobal is nil")
	}
	if cn.GetName().IsEmpty() {
		return fmt.Errorf("name is empty")
	}
	if cn.Pool() == nil {
		return fmt.Errorf("pool is nil")
	}

	cg.mu.Lock()
	defer cg.mu.Unlock()

	cg.Map.Set(cn.GetName(), cn)
	return nil
}

// Remove deletes an AClientRedis from the connection map by its AdapterName.
func (cg *connMapGlobal) Remove(name aconns.AdapterName) {
	if name.IsEmpty() {
		return
	}

	cg.mu.Lock()
	defer cg.mu.Unlock()

	cg.Map.Delete(name)
}

// GetRedisStore retrieves a RedisStore by its prefix from the specified AClientRedis.
func (cg *connMapGlobal) GetRedisStore(name aconns.AdapterName, prefix string) *redisstore.RedisStore {
	db := cg.Get(name)
	if db == nil {
		return nil
	}
	return db.GetRedisStore(prefix)
}

// NewRedisStoreSCS returns a RedisStore for session management.
func NewRedisStoreSCS() *redisstore.RedisStore {
	store := REDIS().GetRedisStore(REDIS_MASTER, REDISSTORE_SCS)
	if store == nil {
		return nil
	}
	return store
}
