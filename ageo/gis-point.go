package ageo

import (
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
	"math"
	"strconv"
)

// GISPoint represents a geographical point with latitude and longitude.
// Latitude corresponds to the Y-axis, and Longitude corresponds to the X-axis.
// This convention is commonly used in geographic information systems (GIS).
type GISPoint struct {
	Latitude  float64 `json:"y,omitempty"` // Latitude of the point.
	Longitude float64 `json:"x,omitempty"` // Longitude of the point.
}

// ToOrb converts to GISPoint to its Orb equivalent.
func (pt *GISPoint) ToOrb() orb.Point {
	return orb.Point{pt.Longitude, pt.Latitude}
}

// Clone returns a deep copy of the GISPoint.
func (pt *GISPoint) Clone() *GISPoint {
	if pt == nil {
		return nil
	}
	return &GISPoint{
		Latitude:  pt.Latitude,
		Longitude: pt.Longitude,
	}
}

// IsGeoCoordinateValid returns true if the GISPoint has valid latitude and longitude
// values according to geographic coordinate limits.
// This includes the (0,0) point, which some systems (e.g., fallback responses) may return.
// Use this for validating raw coordinate bounds without rejecting placeholder locations.
//
// Latitude must be between -90 and 90.
// Longitude must be between -180 and 180.
func (p GISPoint) IsGeoCoordinateValid() bool {
	return p.Latitude >= -90 && p.Latitude <= 90 &&
		p.Longitude >= -180 && p.Longitude <= 180
}

// IsPracticallyValid returns true if the GISPoint has valid, non-zero latitude and longitude values.
// This stricter check is useful for real-world geolocation logic and filtering out placeholder values like (0,0).
//
// Zero coordinates (0,0) are treated as invalid for practical geo use. Why?
// * Many systems, including Firezone and security platforms, treat (0,0) as a fallback/null-equivalent.
// * (0°N, 0°E) is the intersection of the Equator and Prime Meridian.
// * It's located in the Gulf of Guinea, off the coast of West Africa.
// * There is no populated location or real-world data associated with it in most geo IP databases.
func (p GISPoint) IsPracticallyValid() bool {
	return p.IsGeoCoordinateValid() &&
		!(p.Latitude == 0 && p.Longitude == 0)
}

// IsGISZero checks if the GISPoint is at the zero coordinates, which is unlikely to be used in real-life scenarios.
func (pt *GISPoint) IsGISZero() bool {
	// A nil GISPoint or one with both latitude and longitude set to 0 is considered "zero".
	return pt == nil || (pt.Latitude == 0 && pt.Longitude == 0)
}

// ToStringLonLat formats the GISPoint as a string with latitude and longitude.
func (pt *GISPoint) ToStringLonLat() string {
	if pt == nil {
		return ""
	}
	// Format the point using latitude and longitude values.
	return fmt.Sprintf("lat: %s, lon: %s", strconv.FormatFloat(pt.Latitude, 'f', -1, 64), strconv.FormatFloat(pt.Longitude, 'f', -1, 64))
}

// ToStringLat formats the latitude of the GISPoint as a string.
func (pt *GISPoint) ToStringLat() string {
	if pt == nil {
		return ""
	}
	// Convert the latitude to a string.
	return strconv.FormatFloat(pt.Latitude, 'f', -1, 64)
}

// ToStringLon formats the longitude of the GISPoint as a string.
func (pt *GISPoint) ToStringLon() string {
	if pt == nil {
		return ""
	}
	// Convert the longitude to a string.
	return strconv.FormatFloat(pt.Longitude, 'f', -1, 64)
}

// ToStringOSMPoint formats the GISPoint as a string suitable for OpenStreetMap queries.
func (pt *GISPoint) ToStringOSMPoint() string {
	if pt == nil {
		return ""
	}
	// Format the point with longitude first, followed by latitude, separated by a comma.
	return fmt.Sprintf("%s,%s", pt.ToStringLon(), pt.ToStringLat())
}

// GISPoints is a slice of pointers to GISPoint, representing multiple geographical points.
type GISPoints []*GISPoint

// Clean returns a new slice of GISPoints with zero points removed.
func (pts GISPoints) Clean() GISPoints {
	newArr := GISPoints{}
	if pts == nil || len(pts) == 0 {
		return newArr
	}
	for _, pt := range pts {
		if !pt.IsGISZero() {
			newArr = append(newArr, pt)
		}
	}
	return newArr
}

// Clone returns a deep copy of the GISPoints slice.
func (pts GISPoints) Clone() GISPoints {
	if pts == nil || len(pts) == 0 {
		return GISPoints{}
	}
	cloned := make(GISPoints, 0, len(pts))
	for _, pt := range pts {
		if pt == nil {
			cloned = append(cloned, nil)
		} else {
			cloned = append(cloned, pt.Clone())
		}
	}
	return cloned
}

// IsLine checks if the GISPoints represent a line, which requires exactly two points.
func (pts GISPoints) IsLine() bool {
	return pts != nil && len(pts) == 2
}

// IsPolygon checks if the GISPoints represent a polygon, which requires at least three points.
func (pts GISPoints) IsPolygon() bool {
	return pts != nil && len(pts) >= 3
}

// IsMultiPoint checks if the GISPoints has two or more points.
func (pts GISPoints) IsMultiPoint() bool {
	return pts != nil && len(pts) >= 2
}

// ContainsPoint determines whether the target point:
// - exactly matches any point in the set,
// - lies within the polygon (if applicable),
// - lies on a line segment (within a tolerance).
func (pts GISPoints) ContainsPoint(target *GISPoint) bool {
	return pts.ContainsPointWithinTolerance(target, DefaultToleranceMeters)
}

// ContainsPointWithinTolerance checks if the target point:
// - exactly matches a point,
// - lies on a line segment (within tolerance),
// - or lies inside the polygon.
func (pts GISPoints) ContainsPointWithinTolerance(target *GISPoint, toleranceMeters float64) bool {
	if target == nil || !target.IsGeoCoordinateValid() {
		return false
	}

	if toleranceMeters <= 0 {
		toleranceMeters = DefaultToleranceMeters
	}

	cleaned := pts.Clean()
	if len(cleaned) == 0 {
		return false
	}

	tp := orb.Point{target.Longitude, target.Latitude}

	// Exact match
	for _, pt := range cleaned {
		if almostEqual(pt.Latitude, target.Latitude) && almostEqual(pt.Longitude, target.Longitude) {
			return true
		}
	}

	if cleaned.IsLine() {
		p1 := orb.Point{cleaned[0].Longitude, cleaned[0].Latitude}
		p2 := orb.Point{cleaned[1].Longitude, cleaned[1].Latitude}
		distance := planar.DistanceFromSegment(tp, p1, p2)
		return distance <= toleranceMeters
	}

	if cleaned.IsPolygon() {
		ring := make(orb.Ring, len(cleaned))
		for i, pt := range cleaned {
			ring[i] = orb.Point{pt.Longitude, pt.Latitude}
		}
		return planar.RingContains(ring, tp)
	}

	return false
}

// IntersectsRadius determines whether a circular area around the target point
// (defined by a radius in meters) intersects:
// - any point in the set,
// - the line segment if the set forms a line,
// - any edge of the polygon if the set forms a polygon.
func (pts GISPoints) IntersectsRadius(target *GISPoint, radiusMeters float64) bool {
	if target == nil || !target.IsGeoCoordinateValid() || radiusMeters <= 0 {
		return false
	}

	cleaned := pts.Clean()
	if len(cleaned) == 0 {
		return false
	}

	targetPt := target.ToOrb()

	// Check proximity to each point
	for _, pt := range cleaned {
		if geo.Distance(targetPt, pt.ToOrb()) <= radiusMeters {
			return true
		}
	}

	// Check proximity to line segment
	if cleaned.IsLine() {
		a := cleaned[0].ToOrb()
		b := cleaned[1].ToOrb()
		if DistanceFromPointToLineSegmentM(targetPt, a, b) <= radiusMeters {
			return true
		}
	}

	// Check polygon edge proximity
	if cleaned.IsPolygon() {
		for i := 0; i < len(cleaned); i++ {
			j := (i + 1) % len(cleaned)
			a := cleaned[i].ToOrb()
			b := cleaned[j].ToOrb()
			if DistanceFromPointToLineSegmentM(targetPt, a, b) <= radiusMeters {
				return true
			}
		}
	}

	return false
}

// almostEqual compares float64 values with a tolerance suitable for geo coordinates.
func almostEqual(a, b float64) bool {
	const epsilon = 1e-6
	return math.Abs(a-b) <= epsilon
}
