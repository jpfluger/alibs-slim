package ahttp

import (
	"github.com/jpfluger/alibs-slim/azb"
)

/*
CONSIDERING TO DEPRECATE THIS ITEM.
*/

// Constants for pagination navigation control types.
const (
	PAGINATENAV_CONTROL_LABEL_RESULTS = azb.ZBType("label-results") // Label for results.
	PAGINATENAV_CONTROL_LABEL_COUNT   = azb.ZBType("label-count")   // Label for count.
	PAGINATENAV_CONTROL_SHOW_PER_PAGE = azb.ZBType("show-per-page") // Control to show per page options.
	PAGINATENAV_CONTROL_NAVLINKS      = azb.ZBType("nav-links")     // Navigation links.
)

// PaginateNav struct defines the structure for pagination navigation.
type PaginateNav struct {
	azb.Paginate // Embedding azb.Paginate for pagination functionality.

	AddZClick bool   `json:"addZClick,omitempty"` // Flag to add a 'zclick' class.
	ZUrl      string `json:"zurl,omitempty"`      // URL for the Z click action.

	Label string `json:"label,omitempty"` // Label for the pagination navigation.

	ShowPerPage   int            `json:"showPerPage,omitempty"`   // Number of items to show per page.
	LimitsPerPage PaginateLimits `json:"limitsPerPage,omitempty"` // Limits for items per page.

	// Controls specifies the order of display for pagination controls.
	Controls azb.ZBTypes `json:"controls,omitempty"`

	// LinkRender defines which navigation links to render.
	LinkRender struct {
		NoStart bool `json:"noStart,omitempty"` // Do not render the start link.
		NoEnd   bool `json:"noEnd,omitempty"`   // Do not render the end link.
		NoPrev  bool `json:"noPrev,omitempty"`  // Do not render the previous link.
		NoNext  bool `json:"noNext,omitempty"`  // Do not render the next link (corrected json tag from noText to noNext).
	} `json:"linkRender,omitempty"`

	// see azb.Paginate.TotalItems
	//Total int `json:"total,omitempty"` // Total number of items.
}

// GetZClick returns the 'zclick' class if AddZClick is true.
func (pn *PaginateNav) GetZClick() string {
	if pn.AddZClick {
		return " zclick"
	}
	return ""
}
