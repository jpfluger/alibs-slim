package aconns

import (
	"fmt"
	"sort"
)

// AuthAdapterEntry represents a specific authentication adapter instance used
// for a given AuthMethod in a multi-tenant authentication pipeline.
//
// It links a connection (via ConnId) to its underlying adapter instance
// and includes an optional priority to determine execution order.
//
// This structure is typically used within an AuthPipeline to define how
// each authentication method (e.g., primary, MFA, SSPR) should be evaluated
// across different adapters.
type AuthAdapterEntry struct {
	// ConnId refers to the ID of the IConn that this adapter belongs to.
	// It allows the entry to be traced back to its originating connection.
	ConnId ConnId `json:"connId,omitempty"`

	// Adapter is the in-memory adapter instance used to execute the authentication logic.
	// It is excluded from JSON serialization for safety and clarity.
	Adapter IAdapter `json:"-"`

	// Priority is used to determine the evaluation order within a method slice.
	// Lower values are evaluated earlier. This allows fine-grained control over fallback order.
	Priority int `json:"priority,omitempty"`
}

// HasConnId checks if the entry matches the given ConnId.
func (entry *AuthAdapterEntry) HasConnId(id ConnId) bool {
	return entry != nil && entry.ConnId == id
}

// GetAdapter safely returns the adapter if available.
func (entry *AuthAdapterEntry) GetAdapter() IAdapter {
	if entry == nil {
		return nil
	}
	return entry.Adapter
}

// AuthAdapterEntries is a slice of AuthAdapterEntry pointers.
type AuthAdapterEntries []*AuthAdapterEntry

// GetByConnId returns the entry with the given ConnId, if it exists.
func (entries AuthAdapterEntries) GetByConnId(id ConnId) (*AuthAdapterEntry, bool) {
	for _, entry := range entries {
		if entry != nil && entry.ConnId == id {
			return entry, true
		}
	}
	return nil, false
}

// GetAdapters returns a slice of IAdapters in order.
func (entries AuthAdapterEntries) GetAdapters() IAdapters {
	var result IAdapters
	for _, entry := range entries {
		if entry != nil && entry.Adapter != nil {
			result = append(result, entry.Adapter)
		}
	}
	return result
}

// GetConnIds returns all ConnIds in order.
func (entries AuthAdapterEntries) GetConnIds() []ConnId {
	var ids []ConnId
	for _, entry := range entries {
		if entry != nil {
			ids = append(ids, entry.ConnId)
		}
	}
	return ids
}

// SortByPriority modifies the slice to be sorted by ascending Priority.
func (entries AuthAdapterEntries) SortByPriority() {
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i] == nil || entries[j] == nil {
			return false
		}
		return entries[i].Priority < entries[j].Priority
	})
}

// Validate ensures all entries have a valid adapter and non-nil ConnId.
func (entries AuthAdapterEntries) Validate() error {
	for i, entry := range entries {
		if entry == nil {
			return fmt.Errorf("auth adapter entry at index %d is nil", i)
		}
		if entry.ConnId.IsNil() {
			return fmt.Errorf("auth adapter entry at index %d has empty ConnId", i)
		}
		if entry.Adapter == nil {
			return fmt.Errorf("auth adapter entry at index %d has nil adapter", i)
		}
	}
	return nil
}
