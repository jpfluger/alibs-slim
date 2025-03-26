package g_aconns

import (
	"sync"

	"github.com/jpfluger/alibs-slim/aconns"
)

var connMap *connMapGlobal
var muCM sync.RWMutex

type connMapGlobal struct {
	Map   aconns.IConnMap
	Index map[aconns.AdapterName]aconns.ConnId
	mu    sync.RWMutex
}

func init() {
	connMap = &connMapGlobal{Map: aconns.IConnMap{}, Index: map[aconns.AdapterName]aconns.ConnId{}}
}

func CONNS() *connMapGlobal {
	muCM.RLock()
	defer muCM.RUnlock()
	return connMap
}

// Get retrieves an IConn by its UUID.
func (cmg *connMapGlobal) Get(id aconns.ConnId) (aconns.IConn, bool) {
	cmg.mu.RLock()
	defer cmg.mu.RUnlock()
	conn, exists := cmg.Map[id]
	return conn, exists
}

// Set adds or updates an IConn in the map and updates the index.
func (cmg *connMapGlobal) Set(conn aconns.IConn) {
	if conn != nil && !conn.GetId().IsNil() && conn.GetAdapter() != nil {
		cmg.mu.Lock()
		defer cmg.mu.Unlock()
		cmg.Map[conn.GetId()] = conn
		cmg.Index[conn.GetAdapter().GetName()] = conn.GetId()
	}
}

// Remove deletes an IConn from the map by its UUID and updates the index.
func (cmg *connMapGlobal) Remove(id aconns.ConnId) {
	cmg.mu.Lock()
	defer cmg.mu.Unlock()
	if conn, exists := cmg.Map[id]; exists {
		delete(cmg.Index, conn.GetAdapter().GetName())
		delete(cmg.Map, id)
	}
}

// FindByAdapterName finds an IAdapter by its name.
func (cmg *connMapGlobal) FindByAdapterName(name aconns.AdapterName) (aconns.IAdapter, bool) {
	cmg.mu.RLock()
	defer cmg.mu.RUnlock()
	if id, exists := cmg.Index[name]; exists {
		if conn, exists := cmg.Map[id]; exists {
			return conn.GetAdapter(), true
		}
	}
	return nil, false
}

// SetByIConns initializes the global map and index from a slice of IConns.
func (cmg *connMapGlobal) SetByIConns(conns aconns.IConns) {
	cmg.mu.Lock()
	defer cmg.mu.Unlock()
	for _, conn := range conns {
		if conn != nil && !conn.GetId().IsNil() && conn.GetAdapter() != nil {
			cmg.Map[conn.GetId()] = conn
			cmg.Index[conn.GetAdapter().GetName()] = conn.GetId()
		}
	}
}

// Reset reinitializes the global connMap.
func (cmg *connMapGlobal) Reset() {
	cmg.mu.Lock()
	defer cmg.mu.Unlock()
	cmg.Map = aconns.IConnMap{}
	cmg.Index = map[aconns.AdapterName]aconns.ConnId{}
}
