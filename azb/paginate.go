package azb

// IPaginate defines the interface for pagination operations.
type IPaginate interface {
	// SetCurrentPage sets the current page number within valid bounds.
	SetCurrentPage(page int)

	// Offset calculates the offset for database queries.
	Offset() int

	// Limit returns the limit for items per page.
	Limit() int

	// NavNext increments the current page to the next page.
	NavNext()

	// NavPrev decrements the current page to the previous page.
	NavPrev()

	// NavFirst sets the current page to the first page.
	NavFirst()

	// NavLast sets the current page to the last page.
	NavLast()

	// PeekNext returns the next page number without modifying the current page.
	PeekNext() int

	// PeekPrev returns the previous page number without modifying the current page.
	PeekPrev() int
}

// Paginate holds the pagination parameters and provides methods for pagination control.
type Paginate struct {
	CurrentPage  int `json:"cp" query:"cp"` // Current page number, starting at 1.
	TotalItems   int `json:"ti" query:"ti"` // Total number of items to be paginated.
	ItemsPerPage int `json:"ip" query:"ip"` // Number of items per page.
	TotalPages   int `json:"tp" query:"tp"` // Total number of pages, calculated based on TotalItems and ItemsPerPage.
}

// NewPaginate creates a new Paginate instance with the given current page, total items, and items per page.
func NewPaginate(currentPage, totalItems, itemsPerPage int) *Paginate {
	// If itemsPerPage is less than 1, set a default value of 25.
	if itemsPerPage < 1 {
		itemsPerPage = 25
	}
	// If currentPage is less than 1, set it to the first page.
	if currentPage < 1 {
		currentPage = 1
	}
	// If totalItems is less than 1, set it to zero.
	if totalItems < 1 {
		totalItems = 0
	}
	// Calculate the total number of pages based on totalItems and itemsPerPage.
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	// Return a new Paginate instance with validated parameters.
	return &Paginate{
		CurrentPage:  currentPage,
		TotalItems:   totalItems,
		ItemsPerPage: itemsPerPage,
		TotalPages:   totalPages,
	}
}

// SetCurrentPage sets the current page to the provided page number, ensuring it falls within the valid range.
func (p *Paginate) SetCurrentPage(page int) {
	if page < 1 {
		page = 1 // Default to the first page if an invalid page number is provided.
	} else if page > p.TotalPages {
		page = p.TotalPages // Set to the last page if the provided page number exceeds the total pages.
	}
	p.CurrentPage = page
}

// Offset calculates the database query offset based on the current page and items per page.
func (p *Paginate) Offset() int {
	return (p.CurrentPage - 1) * p.ItemsPerPage
}

// Limit returns the number of items per page, to be used as the limit in database queries.
func (p *Paginate) Limit() int {
	return p.ItemsPerPage
}

// NavNext increments the current page number, moving to the next page.
func (p *Paginate) NavNext() {
	p.SetCurrentPage(p.CurrentPage + 1)
}

// NavPrev decrements the current page number, moving to the previous page.
func (p *Paginate) NavPrev() {
	p.SetCurrentPage(p.CurrentPage - 1)
}

// NavFirst sets the current page number to the first page.
func (p *Paginate) NavFirst() {
	p.SetCurrentPage(1)
}

// NavLast sets the current page number to the last page.
func (p *Paginate) NavLast() {
	p.SetCurrentPage(p.TotalPages)
}

// PeekNext returns the next page number without changing the current page.
func (p *Paginate) PeekNext() int {
	if p.CurrentPage < p.TotalPages {
		return p.CurrentPage + 1
	}
	return p.CurrentPage // If on the last page, return the current page.
}

// PeekPrev returns the previous page number without changing the current page.
func (p *Paginate) PeekPrev() int {
	if p.CurrentPage > 1 {
		return p.CurrentPage - 1
	}
	return p.CurrentPage // If on the first page, return the current page.
}
