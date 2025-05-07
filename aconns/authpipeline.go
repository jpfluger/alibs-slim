package aconns

import "fmt"

// AuthPipeline maps AuthMethod → AuthAdapterEntries (ordered slice).
type AuthPipeline map[AuthMethod]AuthAdapterEntries

// GetEntries returns all entries for a specific AuthMethod.
func (ap AuthPipeline) GetEntries(method AuthMethod) AuthAdapterEntries {
	return ap[method]
}

// GetAdapters returns a flat list of adapters for a given method.
func (ap AuthPipeline) GetAdapters(method AuthMethod) IAdapters {
	return ap[method].GetAdapters()
}

// GetConnIds returns ConnIds in order for the given method.
func (ap AuthPipeline) GetConnIds(method AuthMethod) []ConnId {
	return ap[method].GetConnIds()
}

// Validate runs validation across all method entries in the pipeline.
func (ap AuthPipeline) Validate() error {
	for method, entries := range ap {
		if err := entries.Validate(); err != nil {
			return fmt.Errorf("auth pipeline validation failed for method %s: %w", method, err)
		}
	}
	return nil
}

// Methods returns all AuthMethods defined in the pipeline.
func (ap AuthPipeline) Methods() []AuthMethod {
	var methods []AuthMethod
	for method := range ap {
		methods = append(methods, method)
	}
	return methods
}
