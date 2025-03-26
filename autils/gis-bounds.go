package autils

// GISBounds represents geographical bounds with southwest (SW) and northeast (NE) points.
// The structure follows the convention used by MapLibre for defining longitude and latitude bounds.
// Reference: https://maplibre.org/maplibre-gl-js-docs/api/geography/#lnglatboundslike
type GISBounds struct {
	SW GISPoint `json:"sw,omitempty"` // SW represents the southwest corner of the bounds.
	NE GISPoint `json:"ne,omitempty"` // NE represents the northeast corner of the bounds.
}

// HasBounds checks if both the southwest and northeast points of the bounds are defined and non-zero.
func (b *GISBounds) HasBounds() bool {
	// Ensure the GISBounds pointer is not nil and both points are non-zero.
	return b != nil && !b.SW.IsGISZero() && !b.NE.IsGISZero()
}
