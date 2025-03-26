package aclient_http

import "github.com/jpfluger/alibs-slim/aconns"

// To-do: Safe wrapper embedding HTTP functionality with panic handler.

// ISBAdapterHTTP is for sandboxed adapters with HTTP capability.
type ISBAdapterHTTP interface {
	aconns.ISBAdapter

	// From AClientHTTP
}
