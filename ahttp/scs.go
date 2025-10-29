package ahttp

import (
	"github.com/alexedwards/scs/v2"
	"sync"
)

// The initialized SCS manager with associated data store.
var scsInstance *scs.SessionManager
var muSCS sync.RWMutex

// InitializeSCS initializes the global instance of *scs.SessionManager.
func InitializeSCS(mySCS *scs.SessionManager) {
	muSCS.Lock()
	defer muSCS.Unlock()
	scsInstance = mySCS
}

// SCS returns the global instance of *scs.SessionManager.
// It uses read locking for concurrent access safety.
func SCS() *scs.SessionManager {
	muSCS.RLock()
	defer muSCS.RUnlock()
	if scsInstance == nil {
		panic("scsInstance is not initialized")
	}
	return scsInstance
}