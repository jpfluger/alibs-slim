package ahttp

import (
	"github.com/alexedwards/scs/v2"
)

//// The initialized SCS manager with associated data store.
//var scsInstance *scs.SessionManager
//var muSCSInstance sync.RWMutex
//
//func SCS() *scs.SessionManager {
//	muSCSInstance.RLock()
//	defer muSCSInstance.RUnlock()
//	return scsInstance
//}
//
//func SetSCS(mySCS *scs.SessionManager) {
//	muSCSInstance.Lock()
//	defer muSCSInstance.Unlock()
//	scsInstance = mySCS
//}

// Define a global instance for scs.SessionManager
var scsInstance *scs.SessionManager

// InitializeSCS initializes the global instance of IPageSessionController.
// This function should be called once at program startup.
func InitializeSCS(mySCS *scs.SessionManager) {
	if scsInstance != nil {
		panic("scsInstance already initialized")
	}
	scsInstance = mySCS
}

// SCS returns the global instance of *scs.SessionManager.
func SCS() *scs.SessionManager {
	if scsInstance == nil {
		panic("scsInstance is not initialized")
	}
	return scsInstance
}
