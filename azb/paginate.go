package azb

// IPaginate defines the interface for pagination operations.
type IPaginate interface {
	SetCurrentPage(page int)
	Offset() int
	Limit() int
	NavNext()
	NavPrev()
	NavFirst()
	NavLast()
	PeekNext() int
	PeekPrev() int
}

// Paginate holds pagination state and UI config (merged from PaginateNav).
// It implements IPaginate.
type Paginate struct {
	CurrentPage  int `json:"cp" query:"cp"` // Current page number, starting at 1.
	TotalItems   int `json:"ti" query:"ti"` // Total number of items.
	ItemsPerPage int `json:"ip" query:"ip"` // Items per page.
	TotalPages   int `json:"tp" query:"tp"` // Calculated total pages.

	// Optional for UUIDv7/cursor-based queries.
	Cursor string `json:"cursor,omitempty" query:"cursor"`

	// UI fields (from PaginateNav).
	AddZClick bool   `json:"addZClick,omitempty"`
	ZUrl      string `json:"zurl,omitempty"`
	Label     string `json:"label,omitempty"`

	// Per-page options (inlined from PaginateLimits).
	PerPageOptions []PaginateLimit `json:"perPageOptions,omitempty"`

	// Controls order (e.g., label-results, nav-links).
	Controls ZBTypes `json:"controls,omitempty"`

	// Link rendering options.
	LinkRender struct {
		NoStart bool `json:"noStart,omitempty"`
		NoEnd   bool `json:"noEnd,omitempty"`
		NoPrev  bool `json:"noPrev,omitempty"`
		NoNext  bool `json:"noNext,omitempty"`
	} `json:"linkRender,omitempty"`
}

// PaginateLimit (simple struct for per-page options).
type PaginateLimit struct {
	Label     string `json:"label,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	IsDefault bool   `json:"isDefault,omitempty"`
}

// DefaultPerPageOptions is the global default list of per-page options.
// Can be overridden at init time for custom defaults.
var DefaultPerPageOptions = []PaginateLimit{
	{Label: "All", Limit: 0, IsDefault: false},
	{Label: "25", Limit: 25, IsDefault: true},
	{Label: "50", Limit: 50, IsDefault: false},
	{Label: "100", Limit: 100, IsDefault: false},
}

// DefaultPerPageOptionsNoAll is the global default list of per-page options without "All".
var DefaultPerPageOptionsNoAll = []PaginateLimit{
	{Label: "25", Limit: 25, IsDefault: true},
	{Label: "50", Limit: 50, IsDefault: false},
	{Label: "100", Limit: 100, IsDefault: false},
}

// NewPaginate initializes with validation using DefaultPerPageOptions.
func NewPaginate(currentPage, totalItems, itemsPerPage int, cursor string) *Paginate {
	return newPaginateWithOptions(currentPage, totalItems, itemsPerPage, cursor, DefaultPerPageOptions)
}

// NewPaginateNoAll initializes with validation using DefaultPerPageOptionsNoAll.
func NewPaginateNoAll(currentPage, totalItems, itemsPerPage int, cursor string) *Paginate {
	return newPaginateWithOptions(currentPage, totalItems, itemsPerPage, cursor, DefaultPerPageOptionsNoAll)
}

// newPaginateWithOptions is the internal initializer with custom per-page options.
func newPaginateWithOptions(currentPage, totalItems, itemsPerPage int, cursor string, options []PaginateLimit) *Paginate {
	if itemsPerPage < 1 {
		itemsPerPage = 25
	}
	if currentPage < 1 {
		currentPage = 1
	}
	if totalItems < 1 {
		totalItems = 0
	}
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	return &Paginate{
		CurrentPage:    currentPage,
		TotalItems:     totalItems,
		ItemsPerPage:   itemsPerPage,
		TotalPages:     totalPages,
		Cursor:         cursor,
		PerPageOptions: options,
	}
}

// SetPerPageOptions updates the per-page options for this instance.
func (p *Paginate) SetPerPageOptions(options []PaginateLimit) {
	if len(options) > 0 {
		p.PerPageOptions = options
	}
}

// SetCurrentPage bounds the page.
func (p *Paginate) SetCurrentPage(page int) {
	if page < 1 {
		page = 1
	} else if page > p.TotalPages {
		page = p.TotalPages
	}
	p.CurrentPage = page
}

// Offset for queries (0 if cursor mode).
func (p *Paginate) Offset() int {
	if p.Cursor != "" {
		return 0
	}
	return (p.CurrentPage - 1) * p.ItemsPerPage
}

// Limit returns items per page.
func (p *Paginate) Limit() int {
	return p.ItemsPerPage
}

// NavNext advances page or cursor.
func (p *Paginate) NavNext() {
	if p.Cursor != "" {
		// App-specific: Update cursor based on last item.
		return
	}
	p.SetCurrentPage(p.CurrentPage + 1)
}

// NavPrev decrements the current page.
func (p *Paginate) NavPrev() {
	if p.Cursor != "" {
		// App-specific logic.
		return
	}
	p.SetCurrentPage(p.CurrentPage - 1)
}

// NavFirst sets to first page.
func (p *Paginate) NavFirst() {
	if p.Cursor != "" {
		// Reset cursor.
		p.Cursor = ""
	}
	p.SetCurrentPage(1)
}

// NavLast sets to last page.
func (p *Paginate) NavLast() {
	if p.Cursor != "" {
		// Set cursor to last batch.
		return
	}
	p.SetCurrentPage(p.TotalPages)
}

// PeekNext returns next page without changing.
func (p *Paginate) PeekNext() int {
	if p.Cursor != "" {
		return p.CurrentPage + 1 // Or app-specific.
	}
	if p.CurrentPage < p.TotalPages {
		return p.CurrentPage + 1
	}
	return p.CurrentPage
}

// PeekPrev returns previous page without changing.
func (p *Paginate) PeekPrev() int {
	if p.Cursor != "" {
		return p.CurrentPage - 1 // Or app-specific.
	}
	if p.CurrentPage > 1 {
		return p.CurrentPage - 1
	}
	return p.CurrentPage
}

// GetZClick for UI class.
func (p *Paginate) GetZClick() string {
	if p.AddZClick {
		return " zclick"
	}
	return ""
}

// HasPerPageLimit checks for target limit.
func (p *Paginate) HasPerPageLimit(target int) bool {
	for _, opt := range p.PerPageOptions {
		if opt.Limit == target {
			return true
		}
	}
	return false
}

// GetPerPageLimitElseDefault returns target or default.
func (p *Paginate) GetPerPageLimitElseDefault(target int) int {
	var def int
	for _, opt := range p.PerPageOptions {
		if opt.Limit == target {
			return target
		}
		if opt.IsDefault {
			def = opt.Limit
		}
	}
	return def
}

// PageNumbers for display.
func (p *Paginate) PageNumbers() []int {
	pages := make([]int, p.TotalPages)
	for i := range pages {
		pages[i] = i + 1
	}
	return pages
}
