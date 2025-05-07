package aconns

import (
	"errors"
	"fmt"
	"github.com/jpfluger/alibs-slim/auuids"
	"strings"
)

// ConnTenantInfo contains metadata about a tenant-specific connection.
// It is useful for routing, sharding, and access control in multi-tenant systems.
type ConnTenantInfo struct {
	// Region is a geographic or logical zone where the tenant's resources reside.
	// This is used for directing requests to regionally optimized endpoints.
	Region string `json:"region,omitempty"`

	// TenantId uniquely identifies the tenant. Required for most operations.
	TenantId auuids.UUID `json:"tenantId,omitempty"`

	// Priority defines the connection's preference when multiple options are available.
	// Lower values represent higher priority (e.g., Priority 1 is preferred over Priority 2).
	// This is helpful for failover or multi-region load balancing strategies.
	Priority int `json:"priority,omitempty"`

	// Label is an optional, human-readable identifier for the tenant context.
	// For example, this could be "US-West Primary" or "Client-A Staging".
	// Not required for system logic but useful for diagnostics and UI mapping.
	Label string `json:"label,omitempty"`
}

// Validate checks that the required fields in ConnTenantInfo are set correctly.
func (ti ConnTenantInfo) Validate() error {
	if ti.TenantId.IsNil() {
		return errors.New("tenantId is required and cannot be zero")
	}
	if strings.TrimSpace(ti.Region) == "" {
		return errors.New("region is required and cannot be empty")
	}
	if ti.Priority < 0 {
		return fmt.Errorf("priority must be non-negative, got %d", ti.Priority)
	}
	return nil
}

// ConnTenantInfos is a slice of ConnTenantInfo structs.
type ConnTenantInfos []ConnTenantInfo

// Validate checks all entries in the list for correctness.
func (infos ConnTenantInfos) Validate() error {
	for i, info := range infos {
		if err := info.Validate(); err != nil {
			return fmt.Errorf("ConnTenantInfo at index %d: %w", i, err)
		}
	}
	return nil
}

// HasTenant returns true if a tenant with the given ID exists in the list.
func (infos ConnTenantInfos) HasTenant(id auuids.UUID) bool {
	for _, info := range infos {
		if info.TenantId == id {
			return true
		}
	}
	return false
}

// GetByTenantId returns the ConnTenantInfo for the specified tenant, if present.
func (infos ConnTenantInfos) GetByTenantId(id auuids.UUID) (ConnTenantInfo, bool) {
	for _, info := range infos {
		if info.TenantId == id {
			return info, true
		}
	}
	return ConnTenantInfo{}, false
}

// FilterByRegion returns all tenant infos that match the specified region (case-insensitive).
func (infos ConnTenantInfos) FilterByRegion(region string) ConnTenantInfos {
	var result ConnTenantInfos
	region = strings.ToLower(strings.TrimSpace(region))
	for _, info := range infos {
		if strings.ToLower(strings.TrimSpace(info.Region)) == region {
			result = append(result, info)
		}
	}
	return result
}
