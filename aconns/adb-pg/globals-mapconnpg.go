package adb_pg

import (
	"fmt"
	"sync"

	"github.com/jpfluger/alibs-slim/aconns"
)

// gConnMap is the global instance of gConnMapGlobal, which manages PostgreSQL connections.
var gConnMap *gConnMapGlobal

// muCM is a mutex to ensure safe initialization and access to the global gConnMap.
var muCM sync.RWMutex

// gConnMapGlobal encapsulates the global map of PostgreSQL connections and provides thread-safe access.
type gConnMapGlobal struct {
	// Map stores the actual thread-safe mapping of connection IDs to ADBPG instances.
	Map *MapConnPG
	// mu ensures thread-safe operations on the gConnMapGlobal structure itself.
	mu sync.RWMutex
}

// init initializes the global gConnMap instance with a new MapConnPG.
func init() {
	gConnMap = &gConnMapGlobal{Map: NewMapConnPG()}
}

// PGCONNS provides thread-safe access to the global gConnMap instance.
func PGCONNS() *gConnMapGlobal {
	muCM.RLock()
	defer muCM.RUnlock()
	return gConnMap
}

// Get retrieves an ADBPG instance from the connection map using the provided connection ID.
// If the ID is nil or no instance exists, it returns nil.
func (cg *gConnMapGlobal) Get(id aconns.ConnId) *ADBPG {
	if id.IsNil() {
		return nil
	}

	cg.mu.RLock()
	defer cg.mu.RUnlock()

	client, exists := cg.Map.Get(id)
	if !exists {
		return nil
	}
	return client
}

// Set adds or updates an ADBPG instance in the connection map with the given connection ID.
// It returns an error if the instance is nil or has an empty name.
func (cg *gConnMapGlobal) Set(id aconns.ConnId, cn *ADBPG) error {
	if cn == nil {
		return fmt.Errorf("ADBPG instance is nil")
	}
	if cn.GetName().IsEmpty() {
		return fmt.Errorf("connection name is empty")
	}

	cg.mu.Lock()
	defer cg.mu.Unlock()

	cg.Map.Set(id, cn)
	return nil
}

// Remove deletes an ADBPG instance from the connection map by its connection ID.
// If the provided ID is nil, the method does nothing.
func (cg *gConnMapGlobal) Remove(id aconns.ConnId) {
	if id.IsNil() {
		return
	}

	cg.mu.Lock()
	defer cg.mu.Unlock()

	cg.Map.Delete(id)
}
