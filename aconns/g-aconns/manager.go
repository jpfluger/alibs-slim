package g_aconns

import (
	"github.com/jpfluger/alibs-slim/aconns"
)

// Manager is used as the basis for interacting with Connections.
// Use cases can vary. In an SaaS application, consider initializing a
// single Manager as a global, while any client-specific conns are
// initialized separately with their own Manager instance. Depending
// upon your app requirements, you may need to copy and modify this
// structure to suit your needs. This struct is designed to be compliant
// for applications that need to run scripts in a sandboxed environment
// with options to limit client access.
type Manager struct {
	aconns.Manager

	// Implement a custom manager.
}

// NewManager creates a new Manager instance.
func NewManager() *Manager {
	return &Manager{}
}
