package ahttp

import (
	"github.com/labstack/echo/v4"
	"sync"
)

// AuthenticateProvisioner holds URLs for redirection in case of authentication issues
// and provides thread-safe access to these URLs.
type AuthenticateProvisioner struct {
	UrlNoLogin      string       // URL to redirect when no login is detected
	UrlInvalidPerms string       // URL to redirect when invalid permissions are detected
	mu              sync.RWMutex // Mutex to ensure thread-safe access
}

// GetUrlNoLogin safely returns the URL to redirect when no login is detected.
func (ap *AuthenticateProvisioner) GetUrlNoLogin() string {
	ap.mu.RLock()         // Acquire read lock
	defer ap.mu.RUnlock() // Defer the unlocking to the end of the method
	return ap.UrlNoLogin
}

// GetUrlInvalidPerms safely returns the URL to redirect when invalid permissions are detected.
func (ap *AuthenticateProvisioner) GetUrlInvalidPerms() string {
	ap.mu.RLock()         // Acquire read lock
	defer ap.mu.RUnlock() // Defer the unlocking to the end of the method
	return ap.UrlInvalidPerms
}

// LogAuthError logs an authentication error. This method is a placeholder and should be implemented.
func (ap *AuthenticateProvisioner) LogAuthError(c echo.Context, err error) {
	// TODO: Implement the logic to log the authentication error
	// panic("not implemented")
}
