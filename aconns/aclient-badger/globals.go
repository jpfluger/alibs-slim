package aclient_badger

import (
	"fmt"
	"github.com/alexedwards/scs/badgerstore"
	"github.com/jpfluger/alibs-slim/aconns"
	"sync"
)

const (
	BADGERSTORE_SCS         = "scs:"
	BADGERSTORE_FORGOTLOGIN = "forgotlogin:"

	BADGER_MASTER aconns.AdapterName = "badger:master"
)

// connMap is the global connection map.
var connMap *connMapGlobal
var muCM sync.RWMutex

// connMapGlobal holds the global map of Badger connections.
type connMapGlobal struct {
	Map *MapBadger
	mu  sync.RWMutex
}

func init() {
	connMap = &connMapGlobal{Map: NewMapBadger()}
}

// BADGER returns the global connection map.
func BADGER() *connMapGlobal {
	muCM.RLock()
	defer muCM.RUnlock()
	return connMap
}

// Get retrieves an AClientBadger by its AdapterName.
func (cg *connMapGlobal) Get(name aconns.AdapterName) *AClientBadger {
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

// Set adds or updates an AClientBadger in the connection map.
func (cg *connMapGlobal) Set(cn *AClientBadger) error {
	if cn == nil {
		return fmt.Errorf("connMapGlobal is nil")
	}
	if cn.GetName().IsEmpty() {
		return fmt.Errorf("name is empty")
	}
	if cn.DB() == nil {
		return fmt.Errorf("db is nil")
	}

	cg.mu.Lock()
	defer cg.mu.Unlock()

	cg.Map.Set(cn.GetName(), cn)
	return nil
}

// Remove deletes an AClientBadger from the connection map by its AdapterName.
func (cg *connMapGlobal) Remove(name aconns.AdapterName) {
	if name.IsEmpty() {
		return
	}

	cg.mu.Lock()
	defer cg.mu.Unlock()

	cg.Map.Delete(name)
}

// GetBadgerStore retrieves a BadgerStore by its prefix from the specified AClientBadger.
func (cg *connMapGlobal) GetBadgerStore(name aconns.AdapterName, prefix string) *badgerstore.BadgerStore {
	db := cg.Get(name)
	if db == nil {
		return nil
	}
	return db.GetBadgerStore(prefix)
}

// NewBadgerStoreSCS returns a BadgerStore for session management.
func NewBadgerStoreSCS() *badgerstore.BadgerStore {
	store := BADGER().GetBadgerStore(BADGER_MASTER, BADGERSTORE_SCS)
	if store == nil {
		return nil
	}
	return store
}
