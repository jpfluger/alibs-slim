package asessions

// LoginSessionStatusType defines the type for status of a login session.
type LoginSessionStatusType uint32

const (
	// LOGIN_SESSION_STATUS_NONE indicates no status or uninitialized status.
	LOGIN_SESSION_STATUS_NONE LoginSessionStatusType = iota

	// LOGIN_SESSION_STATUS_OK indicates a successful login session.
	LOGIN_SESSION_STATUS_OK

	// LOGIN_SESSION_STATUS_ACTIONS indicates a login session that requires further actions.
	LOGIN_SESSION_STATUS_ACTIONS
)
