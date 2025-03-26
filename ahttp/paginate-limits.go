package ahttp

// PaginateLimit defines the structure for a single pagination limit.
type PaginateLimit struct {
	Label     string `json:"label,omitempty"`     // Display label for the limit.
	Limit     int    `json:"limit,omitempty"`     // Numeric value of the limit.
	IsDefault bool   `json:"isDefault,omitempty"` // Indicates if this limit is the default selection.
}

// PaginateLimits is a slice of pointers to PaginateLimit.
type PaginateLimits []*PaginateLimit

// defaultPageLimits holds the default set of page limits.
var defaultPageLimits = PaginateLimits{
	&PaginateLimit{"All", 0, false},   // Option for showing all items (no limit).
	&PaginateLimit{"25", 25, true},    // Option for 25 items per page.
	&PaginateLimit{"50", 50, false},   // Option for 50 items per page.
	&PaginateLimit{"100", 100, false}, // Option for 100 items per page.
}

// DefaultPaginateLimits returns the default set of pagination limits.
func DefaultPaginateLimits() PaginateLimits {
	return defaultPageLimits
}

// HasLimit checks if the specified target limit exists within the limits.
func (pls PaginateLimits) HasLimit(target int) bool {
	for _, pl := range pls {
		if pl.Limit == target {
			return true // Found the target limit.
		}
	}
	return false // Target limit not found.
}

// GetLimitElseDefault returns the specified target limit if it exists,
// otherwise returns the default limit.
func (pls PaginateLimits) GetLimitElseDefault(target int) int {
	var defaultLimit int // Variable to hold the default limit.
	for _, pl := range pls {
		if pl.Limit == target {
			return target // Return the target limit as it exists.
		}
		if pl.IsDefault {
			defaultLimit = pl.Limit // Store the default limit.
		}
	}
	return defaultLimit // Return the default limit if target limit is not found.
}
