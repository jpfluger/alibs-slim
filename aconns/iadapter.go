package aconns

// IAdapter interface defines the methods that an adapter should implement.
type IAdapter interface {
	GetType() AdapterType
	GetName() AdapterName
	GetHost() string
	GetPort() int
	Validate() error
	Test() (ok bool, testStatus TestStatus, err error)
}

// IAdapters is a slice of IAdapter.
type IAdapters []IAdapter

// ToMap converts the IAdapters to an IAdapterMap.
func (adapters IAdapters) ToMap() IAdapterMap {
	adapterMap := make(IAdapterMap)
	for _, adapter := range adapters {
		if adapter != nil {
			adapterMap[adapter.GetName()] = adapter
		}
	}
	return adapterMap
}

// IAdapterMap is a map of IAdapter with AdapterName as the key.
type IAdapterMap map[AdapterName]IAdapter

// Get retrieves an IAdapter by its AdapterName.
func (m IAdapterMap) Get(name AdapterName) (IAdapter, bool) {
	adapter, exists := m[name]
	return adapter, exists
}

// Set adds or updates an IAdapter in the map.
func (m IAdapterMap) Set(name AdapterName, adapter IAdapter) {
	if adapter != nil {
		m[name] = adapter
	}
}

// Remove deletes an IAdapter from the map by its AdapterName.
func (m IAdapterMap) Remove(name AdapterName) {
	delete(m, name)
}

// ToArray converts the IAdapterMap to an array of IAdapters.
func (m IAdapterMap) ToArray() IAdapters {
	adapters := IAdapters{}
	for _, adapter := range m {
		if adapter != nil {
			adapters = append(adapters, adapter)
		}
	}
	return adapters
}
