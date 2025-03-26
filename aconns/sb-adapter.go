package aconns

// ISBAdapter interface defines the sandboxed interface implementation of IAdapter.
type ISBAdapter interface {
	GetType() AdapterType
	GetName() AdapterName
	GetHost() string
}
