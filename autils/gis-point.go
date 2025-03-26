package autils

import (
	"fmt"
	"strconv"
)

// GISPoint represents a geographical point with latitude and longitude.
// Latitude corresponds to the Y-axis, and Longitude corresponds to the X-axis.
// This convention is commonly used in geographic information systems (GIS).
type GISPoint struct {
	Latitude  float64 `json:"y,omitempty"` // Latitude of the point.
	Longitude float64 `json:"x,omitempty"` // Longitude of the point.
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

// IsLine checks if the GISPoints represent a line, which requires exactly two points.
func (pts GISPoints) IsLine() bool {
	return pts != nil && len(pts) == 2
}

// IsPolygon checks if the GISPoints represent a polygon, which requires at least three points.
func (pts GISPoints) IsPolygon() bool {
	return pts != nil && len(pts) >= 3
}
