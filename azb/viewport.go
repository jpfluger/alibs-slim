package azb

// ViewPortSize represents a simple width and height measurement.
type ViewPortSize struct {
	W int `json:"w"` // Width in pixels
	H int `json:"h"` // Height in pixels
}

// ViewPort holds the dimensions of the window and main content area.
type ViewPort struct {
	Window  ViewPortSize `json:"window"`  // Full browser window size
	Content ViewPortSize `json:"content"` // Inner content area/container size
}

// Width returns the width of the browser window.
func (vp ViewPort) Width() int {
	return vp.Window.W
}

// Height returns the height of the browser window.
func (vp ViewPort) Height() int {
	return vp.Window.H
}

// ContentWidth returns the width of the content container.
func (vp ViewPort) ContentWidth() int {
	return vp.Content.W
}

// ContentHeight returns the height of the content container.
func (vp ViewPort) ContentHeight() int {
	return vp.Content.H
}

// DeviceType returns a simplified device classification based on window width.
//
// This method is not based on a browser-native feature, but instead follows
// widely accepted responsive design conventions. The breakpoints here are
// loosely inspired by Bootstrap's grid system and common UX standards.
//
// - Width < 768px      → "phone"   (small mobile devices)
// - 768px to < 1024px  → "tablet"  (typical tablets and small screens)
// - ≥ 1024px           → "desktop" (laptops, desktops, and large tablets)
//
// These values provide a useful abstraction for server-side rendering logic,
// allowing templates to adjust layout (e.g., cards vs. tables) depending on
// the screen size as reported by the client.
//
// For more granular layout control aligned with Bootstrap’s full breakpoint
// system, see the BootstrapType() method.
func (vp ViewPort) DeviceType() string {
	w := vp.Window.W
	switch {
	case w == 0:
		return "unknown"
	case w < 768:
		return "phone"
	case w < 1024:
		return "tablet"
	default:
		return "desktop"
	}
}

// IsPhone returns true if the window width matches a phone-sized screen.
func (vp ViewPort) IsPhone() bool {
	return vp.DeviceType() == "phone"
}

// IsTablet returns true if the window width matches a tablet-sized screen.
func (vp ViewPort) IsTablet() bool {
	return vp.DeviceType() == "tablet"
}

// IsDesktop returns true if the window width matches a desktop-sized screen.
func (vp ViewPort) IsDesktop() bool {
	return vp.DeviceType() == "desktop"
}

// BootstrapType returns the Bootstrap 5 screen category (e.g., "sm", "md") based on window width.
func (vp ViewPort) BootstrapType() string {
	w := vp.Window.W
	switch {
	case w < 576:
		return "xs"
	case w < 768:
		return "sm"
	case w < 992:
		return "md"
	case w < 1200:
		return "lg"
	case w < 1400:
		return "xl"
	default:
		return "xxl"
	}
}

// IsXs returns true if the viewport matches Bootstrap's "xs" breakpoint.
func (vp ViewPort) IsXs() bool { return vp.BootstrapType() == "xs" }

// IsSm returns true if the viewport matches Bootstrap's "sm" breakpoint.
func (vp ViewPort) IsSm() bool { return vp.BootstrapType() == "sm" }

// IsMd returns true if the viewport matches Bootstrap's "md" breakpoint.
func (vp ViewPort) IsMd() bool { return vp.BootstrapType() == "md" }

// IsLg returns true if the viewport matches Bootstrap's "lg" breakpoint.
func (vp ViewPort) IsLg() bool { return vp.BootstrapType() == "lg" }

// IsXl returns true if the viewport matches Bootstrap's "xl" breakpoint.
func (vp ViewPort) IsXl() bool { return vp.BootstrapType() == "xl" }

// IsXxl returns true if the viewport matches Bootstrap's "xxl" breakpoint.
func (vp ViewPort) IsXxl() bool { return vp.BootstrapType() == "xxl" }

// BootstrapBreakpoints maps Bootstrap 5 breakpoint names to their min-widths in pixels.
var BootstrapBreakpoints = map[string]int{
	"xs":  0,
	"sm":  576,
	"md":  768,
	"lg":  992,
	"xl":  1200,
	"xxl": 1400,
}

// IsLtBootstrap returns true if the viewport width is less than the given Bootstrap breakpoint.
func (vp ViewPort) IsLtBootstrap(bp string) bool {
	w := vp.Window.W
	if min, ok := BootstrapBreakpoints[bp]; ok {
		return w < min
	}
	return false // unknown breakpoint
}

// IsLteBootstrap returns true if the viewport width is less than or equal to the given Bootstrap breakpoint.
func (vp ViewPort) IsLteBootstrap(bp string) bool {
	w := vp.Window.W
	if min, ok := BootstrapBreakpoints[bp]; ok {
		return w <= min
	}
	return false
}

// IsGteBootstrap returns true if the viewport width is greater than or equal to the given Bootstrap breakpoint.
func (vp ViewPort) IsGteBootstrap(bp string) bool {
	w := vp.Window.W
	if min, ok := BootstrapBreakpoints[bp]; ok {
		return w >= min
	}
	return false
}

// IsGtBootstrap returns true if the viewport width is greater than the given Bootstrap breakpoint.
func (vp ViewPort) IsGtBootstrap(bp string) bool {
	w := vp.Window.W
	if min, ok := BootstrapBreakpoints[bp]; ok {
		return w > min
	}
	return false
}
