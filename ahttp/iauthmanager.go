package ahttp

import (
	"github.com/jpfluger/alibs-slim/asessions"
	"github.com/jpfluger/alibs-slim/auser"
	"sync"
)

type IAuthManager interface {
	//NewLoginUserSession(username asessions.Username, password string) (*asessions.UserSessionPerm, error)
	NewLoginUserSession(username auser.Username, authType string, secret string) (*asessions.UserSessionPerm, error)
}

var (
	authManagerInstance   IAuthManager
	muAuthManagerInstance sync.RWMutex
)

// SetAuthManager initializes the global instance of IAuthManager.
func SetAuthManager(authManager IAuthManager) {
	muAuthManagerInstance.Lock()
	defer muAuthManagerInstance.Unlock()
	authManagerInstance = authManager
}

// AUTHMANAGER returns the global instance of IAuthManager.
func AUTHMANAGER() IAuthManager {
	muAuthManagerInstance.RLock()
	defer muAuthManagerInstance.RUnlock()
	return authManagerInstance
}
