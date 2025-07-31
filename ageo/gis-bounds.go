package ageo

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/orb/project"
)

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

// ContainsPoint returns true if the target point is within the bounding box defined by SW and NE corners.
func (b *GISBounds) ContainsPoint(target *GISPoint) bool {
	if b == nil || !b.HasBounds() || target == nil || !target.IsGeoCoordinateValid() {
		return false
	}

	bound := orb.Bound{
		Min: orb.Point{b.SW.Longitude, b.SW.Latitude},
		Max: orb.Point{b.NE.Longitude, b.NE.Latitude},
	}

	return bound.Contains(orb.Point{target.Longitude, target.Latitude})
}

// IntersectsRadius returns true if the radius around the target point intersects the bounds.
// Uses bounding box expansion and distance approximation for edge contact.
func (b *GISBounds) IntersectsRadius(target *GISPoint, radiusMeters float64) bool {
	if b == nil || !b.HasBounds() || target == nil || !target.IsGeoCoordinateValid() || radiusMeters <= 0 {
		return false
	}

	if b.ContainsPoint(target) {
		return true
	}

	// Convert target point to projected meters
	tp := project.WGS84.ToMercator(target.ToOrb())

	// Define bounds edges and project to meters
	edges := [][2]orb.Point{
		// South edge
		{
			project.WGS84.ToMercator(orb.Point{b.SW.Longitude, b.SW.Latitude}),
			project.WGS84.ToMercator(orb.Point{b.NE.Longitude, b.SW.Latitude}),
		},
		// North edge
		{
			project.WGS84.ToMercator(orb.Point{b.SW.Longitude, b.NE.Latitude}),
			project.WGS84.ToMercator(orb.Point{b.NE.Longitude, b.NE.Latitude}),
		},
		// West edge
		{
			project.WGS84.ToMercator(orb.Point{b.SW.Longitude, b.SW.Latitude}),
			project.WGS84.ToMercator(orb.Point{b.SW.Longitude, b.NE.Latitude}),
		},
		// East edge
		{
			project.WGS84.ToMercator(orb.Point{b.NE.Longitude, b.SW.Latitude}),
			project.WGS84.ToMercator(orb.Point{b.NE.Longitude, b.NE.Latitude}),
		},
	}

	// Check if the projected target is within threshold distance from any edge
	for _, edge := range edges {
		dist := planar.DistanceFromSegment(edge[0], edge[1], tp)
		if dist <= radiusMeters {
			return true
		}
	}

	return false
}
