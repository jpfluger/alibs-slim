package aconns

// TestStatus represents the status of the connection test.
type TestStatus string

const (
	TESTSTATUS_INITIALIZED            TestStatus = "initialized"
	TESTSTATUS_INITIALIZED_SUCCESSFUL TestStatus = "initialized+test-successful"
	TESTSTATUS_FAILED                 TestStatus = "test-failed"
)
