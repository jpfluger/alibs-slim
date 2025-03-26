package aclient_redis

import "github.com/jpfluger/alibs-slim/aconns"

// To-do: Safe wrapper embedding REDIS functionality with panic handler.

// ISBAdapterREDIS is for sandboxed adapters with REDIS capability.
type ISBAdapterREDIS interface {
	aconns.ISBAdapter

	// From AClientREDIS
}
