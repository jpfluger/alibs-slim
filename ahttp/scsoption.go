package ahttp

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
)

// SCSOptions is a configuration struct for session cookies used with scs.SessionManager.
// It provides session management customization options and ensures secure and structured
// handling of cookies and sessions. Use in conjunction with the alexedwards/scs library.
// See: https://github.com/alexedwards/scs/blob/master/session.go
// The MIT License (MIT); Copyright (c) 2016 Alex Edwards
type SCSOptions struct {
	// Name specifies the name of the session cookie, ensuring uniqueness if multiple sessions are used.
	// Default: "session"
	Name string `json:"name"`

	// Domain specifies the 'Domain' attribute on the session cookie.
	// Default: the domain that issued the cookie.
	Domain string `json:"domain"`

	// HttpOnly sets the 'HttpOnly' attribute on the session cookie.
	// Default: true
	HttpOnly bool `json:"httpOnly"`

	// Path sets the 'Path' attribute on the session cookie.
	// Default: "/"
	Path string `json:"path"`

	// Persist indicates whether the session cookie should be retained after the browser is closed (sets Max-Age).
	// Default: true
	Persist bool `json:"persist"`

	// SameSite controls the 'SameSite' attribute on the session cookie.
	// Default: http.SameSiteLaxMode
	SameSite http.SameSite `json:"sameSite"`

	// Secure indicates whether the 'Secure' attribute on the session cookie is set (recommended for production).
	// Default: false (set to true in prod)
	Secure bool `json:"secure"`

	// Partitioned sets the 'Partitioned' attribute on the session cookie for enhanced privacy.
	// Default: false
	Partitioned bool `json:"partitioned"`

	// IdleTimeoutMinutes specifies the maximum session inactivity period before expiration, in minutes.
	// Set to 0 to disable. Default: 0 (disabled)
	IdleTimeoutMinutes int `json:"idleTimeoutMinutes"`

	// LifetimeMinutes sets the session's absolute maximum duration, independent of activity.
	// Must be >0; default: 1440 minutes (24 hours) if <=0.
	LifetimeMinutes int `json:"lifetimeMinutes"`

	// HashTokenInStore controls whether to hash the session token (SHA-256) before storing it.
	// Enhances security for stores like Redis. Default: false
	HashTokenInStore bool `json:"hashTokenInStore"`
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

	// Use local variables to avoid mutating the input struct.
	cookieName := strings.TrimSpace(ss.Name)
	if cookieName == "" {
		cookieName = "SID" // default
	}

	idleTimeoutMinutes := ss.IdleTimeoutMinutes
	if idleTimeoutMinutes < 0 {
		idleTimeoutMinutes = 0 // Default: disabled
	}

	lifetimeMinutes := ss.LifetimeMinutes
	if lifetimeMinutes <= 0 {
		lifetimeMinutes = 1440 // Default: 24 hours
	}

	// Basic validation: Lifetime should exceed IdleTimeout if set.
	if idleTimeoutMinutes > 0 && lifetimeMinutes <= idleTimeoutMinutes {
		return nil, fmt.Errorf("LifetimeMinutes (%d) must exceed IdleTimeoutMinutes (%d) if set", lifetimeMinutes, idleTimeoutMinutes)
	}

	// Initialize the session manager with the provided scs store.
	sessionManager := scs.New()
	sessionManager.Store = scsStore

	// Register types with gob for use in session data encoding.
	for _, gobType := range gobRegister {
		gob.Register(gobType)
	}

	// Set session cookie attributes based on SCSOptions configuration.
	sessionManager.Cookie.Name = cookieName
	if ss.Domain != "" {
		sessionManager.Cookie.Domain = ss.Domain
	}
	if ss.Path != "" {
		sessionManager.Cookie.Path = ss.Path
	} else {
		sessionManager.Cookie.Path = "/" // Default
	}
	sessionManager.Cookie.HttpOnly = ss.HttpOnly // Zero value is false, but default to true in code if needed
	sessionManager.Cookie.Secure = ss.Secure
	sessionManager.Cookie.SameSite = ss.SameSite
	sessionManager.Cookie.Persist = ss.Persist
	sessionManager.Cookie.Partitioned = ss.Partitioned

	// Set session timeouts and lifetime.
	sessionManager.IdleTimeout = time.Duration(idleTimeoutMinutes) * time.Minute
	sessionManager.Lifetime = time.Duration(lifetimeMinutes) * time.Minute

	// Set hashing option.
	sessionManager.HashTokenInStore = ss.HashTokenInStore

	if addToGlobalSCS {
		InitializeSCS(sessionManager)
	}

	return sessionManager, nil
}
