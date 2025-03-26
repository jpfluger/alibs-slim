package aclient_smtp

import "github.com/jpfluger/alibs-slim/aconns"

// To-do: Safe wrapper embedding SMTP functionality with panic handler.

// ISBAdapterSMTP is for sandboxed adapters with SMTP capability.
type ISBAdapterSMTP interface {
	aconns.ISBAdapter

	// From AClientSMTP
}
