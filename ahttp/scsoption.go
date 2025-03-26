package ahttp

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"net/http"
	"strings"
	"sync"
	"time"
)

// SCSOptions is a configuration struct for session cookies used with scs.SessionManager.
// It provides session management customization options and ensures secure and structured
// handling of cookies and sessions. Use in conjunction with the alexedwards/scs library.
// See: https://github.com/alexedwards/scs/blob/master/session.go
// The MIT License (MIT); Copyright (c) 2016 Alex Edwards
type SCSOptions struct {
	// Name specifies the name of the session cookie, ensuring uniqueness if multiple sessions are used.
	Name string `json:"name"`

	// Domain specifies the 'Domain' attribute on the session cookie, defaulting to the domain that issued the cookie.
	Domain string `json:"domain"`

	// HttpOnly sets the 'HttpOnly' attribute on the session cookie, defaulting to true.
	HttpOnly bool `json:"httpOnly"`

	// Path sets the 'Path' attribute on the session cookie, defaulting to "/" (root).
	Path string `json:"path"`

	// Persist indicates whether the session cookie should be retained after the browser is closed.
	Persist bool `json:"persist"`

	// SameSite controls the 'SameSite' attribute on the session cookie, with a default value of Lax.
	SameSite http.SameSite `json:"sameSite"`

	// Secure indicates whether the 'Secure' attribute on the session cookie is set (recommended for production).
	Secure bool `json:"secure"`

	// IdleTimeoutMinutes specifies the maximum session inactivity period before expiration, in minutes. Default is 1440.
	IdleTimeoutMinutes int `json:"idleTimeoutMinutes"`

	// LifetimeMinutes sets the session's absolute maximum duration, independent of activity. Default is 1440 minutes (24 hours).
	LifetimeMinutes int `json:"lifetimeMinutes"`

	// IsTokenRequired specifies whether a token is generated on LoadCheck for all requests.
	IsTokenRequired bool `json:"isTokenRequired"`

	// Mutex for concurrent access protection.
	mu sync.RWMutex
}

// Initialize configures a new scs.SessionManager based on SCSOptions settings and an scs.Store for session storage.
// Registers types with gob for encoding/decoding session data if needed.
// If addToGlobalSCS is true, then the scs.SessionManager is auto-added to scsInstance.
// Returns a configured session manager or an error if initialization fails.
func (ss *SCSOptions) Initialize(scsStore scs.Store, gobRegister []interface{}, addToGlobalSCS bool) (*scs.SessionManager, error) {
	if ss == nil {
		return nil, fmt.Errorf("SCSOptions is nil")
	}

	if scsStore == nil {
		return nil, fmt.Errorf("scsStore parameter is nil")
	}

	// Ensure session cookie name is not empty; default to "SID".
	ss.Name = strings.TrimSpace(ss.Name)
	if ss.Name == "" {
		ss.Name = "SID"
	}

	// Ensure IdleTimeout and Lifetime are set to valid values.
	if ss.IdleTimeoutMinutes < 0 {
		ss.IdleTimeoutMinutes = 0
	}
	if ss.LifetimeMinutes <= 0 {
		// Default Lifetime to 1440 minutes (24 hours) if unset or zero.
		ss.LifetimeMinutes = 1440
	}

	// Initialize the session manager with the provided scs store.
	sessionManager := scs.New()
	sessionManager.Store = scsStore

	// Register types with gob for use in session data encoding.
	for _, gobType := range gobRegister {
		gob.Register(gobType)
	}

	// Set session cookie attributes based on SCSOptions configuration.
	sessionManager.Cookie.Name = ss.Name
	if ss.Domain != "" {
		sessionManager.Cookie.Domain = ss.Domain
	}
	if ss.Path != "" {
		sessionManager.Cookie.Path = ss.Path
	}
	if ss.SameSite >= http.SameSiteDefaultMode && ss.SameSite <= http.SameSiteNoneMode {
		sessionManager.Cookie.SameSite = ss.SameSite
	}
	sessionManager.Cookie.HttpOnly = ss.HttpOnly
	sessionManager.Cookie.Secure = ss.Secure

	// Set session timeouts and lifetime.
	if ss.IdleTimeoutMinutes > 0 {
		sessionManager.IdleTimeout = time.Duration(ss.IdleTimeoutMinutes) * time.Minute
	}
	sessionManager.Lifetime = time.Duration(ss.LifetimeMinutes) * time.Minute

	if addToGlobalSCS {
		InitializeSCS(sessionManager)
	}

	return sessionManager, nil
}
