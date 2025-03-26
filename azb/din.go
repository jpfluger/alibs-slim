package azb

// DIN represents a Data Input Node that includes an action and can be extended with additional data.
type DIN struct {
	ZAction ZAction `json:"zaction"` // The action associated with this data input node.
	// Extend DIN with specific or generic data structures as needed.
	// Data Contact `json:"data"` // Uncomment to use a specific struct.
	// Data interface{} `json:"data"` // Uncomment to use a generic struct.
}

// IDINPaginate is an interface for data input nodes that support pagination.
type IDINPaginate interface {
	Validate() error                                    // Validates the data input node.
	NewPaginate(totalItems, itemsPerPage int) *Paginate // Creates a new Paginate instance with total items and items per page.
}

// Validate synchronizes the page limit with the ZAction's page limit and resets the current page.
func (dinp *DIN) Validate() error {
	if dinp.ZAction.PageLimit != dinp.ZAction.PageLimit {
		dinp.ZAction.PageLimit = dinp.ZAction.PageLimit // Update the ZAction's page limit if different.
		dinp.ZAction.PageOn = 1                         // Reset to the first page.
	}
	return nil // No validation errors.
}

// NewPaginate creates a new Paginate instance with the provided total items and items per page.
func (dinp *DIN) NewPaginate(totalItems, itemsPerPage int) *Paginate {
	return NewPaginate(dinp.ZAction.PageOn, totalItems, itemsPerPage)
}
