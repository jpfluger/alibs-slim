package azb

import "errors"

// IDINPaginate is an interface for data input nodes that support pagination.
type IDINPaginate interface {
	Validate() error
	NewPaginate(totalItems int, cursor string) (*Paginate, error)
}

// DIN integrates ZAction for data input, with pagination.
// It implements IDINPaginate.
type DIN struct {
	ZAction ZAction `json:"zaction"`

	// Embed Paginate for direct access after validation.
	Paginate *Paginate `json:"paginate,omitempty"`
}

// Validate handles ZAction params, aligning with client JS (parseIntOrZero).
func (din *DIN) Validate() error {
	if din.ZAction.PageLimit <= 0 {
		din.ZAction.PageLimit = 25 // Default, or use GetPerPageLimitElseDefault.
	}
	if din.ZAction.PageOn <= 0 {
		din.ZAction.PageOn = 1
	}
	// Additional validation (e.g., for cursor if UUID query).
	return nil
}

// NewPaginate builds from ZAction, with totalItems from query results.
func (din *DIN) NewPaginate(totalItems int, cursor string) (*Paginate, error) {
	if err := din.Validate(); err != nil {
		return nil, err
	}
	if totalItems < 0 {
		return nil, errors.New("totalItems must be non-negative")
	}
	p := NewPaginate(din.ZAction.PageOn, totalItems, din.ZAction.PageLimit, cursor)
	din.Paginate = p
	return p, nil
}

// NewPaginateNoAll builds from ZAction, with totalItems from query results.
func (din *DIN) NewPaginateNoAll(totalItems int, cursor string) (*Paginate, error) {
	if err := din.Validate(); err != nil {
		return nil, err
	}
	if totalItems < 0 {
		return nil, errors.New("totalItems must be non-negative")
	}
	p := NewPaginateNoAll(din.ZAction.PageOn, totalItems, din.ZAction.PageLimit, cursor)
	din.Paginate = p
	return p, nil
}
