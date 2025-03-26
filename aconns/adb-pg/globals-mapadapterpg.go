package adb_pg

import (
	"fmt"
	"sync"

	"github.com/jpfluger/alibs-slim/aconns"
)

// PG_MASTER is the AdapterName for the master PostgreSQL connection.
const PG_MASTER aconns.AdapterName = "pg:master"

// gAdapterMap is the global connection map.
var gAdapterMap *adapterMapGlobal
var gMUAdapterMap sync.RWMutex

// adapterMapGlobal holds the global map of PostgreSQL connections.
type adapterMapGlobal struct {
	Map *MapAdapterPG
	mu  sync.RWMutex
}

func init() {
	gAdapterMap = &adapterMapGlobal{Map: NewMapAdapterPG()}
}

// PGADAPTERS returns the global adapters map.
func PGADAPTERS() *adapterMapGlobal {
	gMUAdapterMap.RLock()
	defer gMUAdapterMap.RUnlock()
	return gAdapterMap
}

// Get retrieves an ADBPG by its AdapterName.
func (cg *adapterMapGlobal) Get(name aconns.AdapterName) *ADBPG {
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

// Set adds or updates an ADBPG in the connection map.
func (cg *adapterMapGlobal) Set(cn *ADBPG) error {
	if cn == nil {
		return fmt.Errorf("adapterMapGlobal is nil")
	}
	if cn.GetName().IsEmpty() {
		return fmt.Errorf("name is empty")
	}

	cg.mu.Lock()
	defer cg.mu.Unlock()

	cg.Map.Set(cn.GetName(), cn)
	return nil
}

// Remove deletes an ADBPG from the connection map by its AdapterName.
func (cg *adapterMapGlobal) Remove(name aconns.AdapterName) {
	if name.IsEmpty() {
		return
	}

	cg.mu.Lock()
	defer cg.mu.Unlock()

	cg.Map.Delete(name)
}
