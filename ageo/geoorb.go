package ageo

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/orb/project"
	"math"
)

const (
	// EarthRadiusKM defines the Earth's radius in kilometers.
	// Used for haversine calculations when km output is preferred.
	EarthRadiusKM = 6371.0

	// EarthRadiusM defines the Earth's radius in meters.
	// Used for accurate great-circle distance calculations in meters.
	EarthRadiusM = 6371000.0

	// DistanceEpsilon is a generic precision threshold (in meters)
	// used in approximate comparisons or equality checks between distances.
	DistanceEpsilon = 1.0

	// LineMatchEpsilon defines the tolerance (in meters) for considering
	// a point to lie "approximately" on a line segment. Higher = more forgiving.
	// Six (6) meters is appropriate for GPS-level accuracy. It allows for minor
	// drift, but still ensures a point is geometrically aligned with the segment.
	LineMatchEpsilon = 6.0

	// DefaultToleranceMeters is a fallback threshold for radius-based
	// proximity checks when none is provided explicitly. It's designed
	// to fail fast unless the user is truly nearby (e.g., standing next
	// to a device or physical marker).
	DefaultToleranceMeters = 5.0

	// DefaultGeoFencingMeters defines the default proximity threshold (in meters)
	// used for geofencing operations — such as determining if a point lies near the
	// edge of a polygon or within a surrounding radius.
	//
	// Why 6000 meters?
	// - It accounts for real-world imprecision in GPS coordinates and polygon outlines.
	// - It captures users who are **near**, but not strictly inside, a defined area.
	// - It's suitable for applications like:
	//     * Campus or facility perimeter detection
	//     * Regional geofencing (e.g., city borders, delivery zones)
	//     * Maritime or airspace boundary monitoring
	//
	// Unlike tighter tolerances (like LineMatchEpsilon or DefaultToleranceMeters),
	// this threshold supports fuzzy spatial logic rather than precise containment.
	DefaultGeoFencingMeters = 6000.0
)

// HaversineDistanceBetweenPointsKM calculates the great-circle distance in kilometers
// between two geographic points using the Haversine formula.
func HaversineDistanceBetweenPointsKM(a, b *GISPoint) float64 {
	return geo.DistanceHaversine(a.ToOrb(), b.ToOrb()) / 1000.0
}

// HaversineDistanceKM computes the Haversine distance in kilometers
// from raw latitude and longitude pairs.
func HaversineDistanceKM(lat1, lon1, lat2, lon2 float64) float64 {
	return HaversineDistanceM(lat1, lon1, lat2, lon2) / 1000.0
}

// HaversineDistanceBetweenPointsM calculates the great-circle distance in meters
// between two GISPoint coordinates.
func HaversineDistanceBetweenPointsM(a, b *GISPoint) float64 {
	return geo.DistanceHaversine(a.ToOrb(), b.ToOrb())
}

// HaversineDistanceM computes the Haversine distance in meters
// between two latitude/longitude pairs.
func HaversineDistanceM(lat1, lon1, lat2, lon2 float64) float64 {
	p1 := orb.Point{lon1, lat1}
	p2 := orb.Point{lon2, lat2}
	return geo.Distance(p1, p2)
}

// IsPointInsidePolygon returns true if the point lies inside a polygon defined by GISPoints.
// Uses orb's planar.RingContains algorithm which assumes a closed ring.
func IsPointInsidePolygon(point *GISPoint, polygon GISPoints) bool {
	ring := make(orb.Ring, 0, len(polygon))
	for _, pt := range polygon {
		ring = append(ring, pt.ToOrb())
	}
	return planar.RingContains(ring, point.ToOrb())
}

// IsPointApproximatelyOnLine returns true if point p is within LineMatchEpsilon
// meters of the line segment from point a to b.
func IsPointApproximatelyOnLine(a, b, p *GISPoint) bool {
	return DistanceFromPointToLineSegmentM(p.ToOrb(), a.ToOrb(), b.ToOrb()) <= LineMatchEpsilon
}

// DistanceFromPointToLineSegment returns the orthogonal distance in meters
// from point p to the segment defined by points a–b. Nil-safe fallback returns max float.
func DistanceFromPointToLineSegment(a, b, p *GISPoint) float64 {
	if a == nil || b == nil || p == nil {
		return math.MaxFloat64
	}

	// Convert to orb.Point and project to Mercator (meters)
	pa := project.WGS84.ToMercator(a.ToOrb())
	pb := project.WGS84.ToMercator(b.ToOrb())
	pp := project.WGS84.ToMercator(p.ToOrb())

	return planar.DistanceFromSegment(pa, pb, pp)
}

// DistanceFromPointToLineSegmentM computes distance in meters from point `p`
// to the segment [a–b], all given as orb.Points. Uses Mercator projection.
func DistanceFromPointToLineSegmentM(p, a, b orb.Point) float64 {
	pM := project.WGS84.ToMercator(p)
	aM := project.WGS84.ToMercator(a)
	bM := project.WGS84.ToMercator(b)
	return planar.DistanceFromSegment(aM, bM, pM)
}

// DistanceFromSegmentMeters estimates the orthogonal distance in meters from point p
// to a geographic line segment a–b using an equirectangular approximation.
func DistanceFromSegmentMeters(p, a, b orb.Point) float64 {
	const R = EarthRadiusM // Earth radius in meters

	meanLat := (a[1] + b[1]) / 2.0
	cosMeanLat := math.Cos(meanLat * math.Pi / 180)

	ax := a[0] * cosMeanLat
	ay := a[1]
	bx := b[0] * cosMeanLat
	by := b[1]
	px := p[0] * cosMeanLat
	py := p[1]

	dx := bx - ax
	dy := by - ay

	if dx == 0 && dy == 0 {
		return geo.Distance(p, a) // degenerate segment
	}

	// Project p onto line a–b
	t := ((px-ax)*dx + (py-ay)*dy) / (dx*dx + dy*dy)
	t = math.Max(0, math.Min(1, t)) // clamp to segment

	closestX := ax + t*dx
	closestY := ay + t*dy
	closest := orb.Point{closestX / cosMeanLat, closestY}

	return geo.Distance(p, closest)
}
