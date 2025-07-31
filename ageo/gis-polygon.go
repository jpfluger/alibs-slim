package ageo

// GISPolygon represents a polygon with a set of points, background color, and opacity.
type GISPolygon struct {
	Points          GISPoints `json:"points,omitempty"`  // Points defining the polygon.
	BackgroundColor string    `json:"bg,omitempty"`      // Background color of the polygon.
	Opacity         float64   `json:"opacity,omitempty"` // Opacity of the polygon, ranging from 0 (transparent) to 1 (opaque).
}

// HasBoundary checks if the GISPolygon has a valid boundary defined by at least three points.
func (b *GISPolygon) HasBoundary() bool {
	// Ensure the GISPolygon pointer is not nil and the points define a polygon.
	return b != nil && b.Points.IsPolygon()
}

// Clean sanitizes the GISPolygon by cleaning up the points and ensuring the opacity is within valid range.
func (b *GISPolygon) Clean() {
	// Clean the points to remove any zero points.
	b.Points = b.Points.Clean()

	// If the cleaned points do not form a polygon, reset the points to an empty slice.
	if !b.HasBoundary() {
		b.Points = GISPoints{}
	}

	// Ensure the opacity is within the range [0, 1]. If not, set it to fully opaque.
	if b.Opacity < 0 || b.Opacity > 1 {
		b.Opacity = 1
	}
}
