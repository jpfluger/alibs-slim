package aclient_ldap

import "github.com/jpfluger/alibs-slim/aconns"

// To-do: Safe wrapper embedding LDAP functionality with panic handler.

// ISBAdapterLDAP is for sandboxed adapters with LDAP capability.
type ISBAdapterLDAP interface {
	aconns.ISBAdapter

	// From AClientLDAP
}
