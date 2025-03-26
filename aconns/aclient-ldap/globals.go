package aclient_ldap

import (
	"fmt"
	"sync"

	"github.com/jpfluger/alibs-slim/aconns"
)

// LDAP_MASTER is the AdapterName for the master LDAP connection.
const LDAP_MASTER aconns.AdapterName = "ldap:master"

// connMap is the global connection map.
var connMap *connMapGlobal
var muCM sync.RWMutex

// connMapGlobal holds the global map of LDAP connections.
type connMapGlobal struct {
	Map *MapLDAP
	mu  sync.RWMutex
}

func init() {
	connMap = &connMapGlobal{Map: NewMapLDAP()}
}

// LDAP returns the global connection map.
func LDAP() *connMapGlobal {
	muCM.RLock()
	defer muCM.RUnlock()
	return connMap
}

// Get retrieves an AClientLDAP by its AdapterName.
func (cg *connMapGlobal) Get(name aconns.AdapterName) *AClientLDAP {
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

// Set adds or updates an AClientLDAP in the connection map.
func (cg *connMapGlobal) Set(cn *AClientLDAP) error {
	if cn == nil {
		return fmt.Errorf("connMapGlobal is nil")
	}
	if cn.GetName().IsEmpty() {
		return fmt.Errorf("name is empty")
	}

	cg.mu.Lock()
	defer cg.mu.Unlock()

	cg.Map.Set(cn.GetName(), cn)
	return nil
}

// Remove deletes an AClientLDAP from the connection map by its AdapterName.
func (cg *connMapGlobal) Remove(name aconns.AdapterName) {
	if name.IsEmpty() {
		return
	}

	cg.mu.Lock()
	defer cg.mu.Unlock()

	cg.Map.Delete(name)
}
