package aclient_ftp

import "github.com/jpfluger/alibs-slim/aconns"

// To-do: Safe wrapper embedding FTP functionality with panic handler.

// ISBAdapterFTP is for sandboxed adapters with FTP capability.
type ISBAdapterFTP interface {
	aconns.ISBAdapter

	// From AClientFTP
	Open() error
	Close() error

	// Wrapper for "github.com/jlaffaye/ftp"
	//ChangeDir(path string) error
	//ChangeDirToParent() error
	//CurrentDir() (string, error)
	//FileSize(path string) (int64, error)
	//GetTime(path string) (time.Time, error)
	//IsGetTimeSupported() bool
}
