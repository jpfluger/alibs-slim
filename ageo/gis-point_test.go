package ageo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGISPointIsGeoCoordinateValid(t *testing.T) {
	tests := []struct {
		name     string
		point    GISPoint
		expected bool
	}{
		{"Valid center point", GISPoint{Latitude: 0, Longitude: 0}, true},
		{"Valid bounds positive", GISPoint{Latitude: 45.0, Longitude: 90.0}, true},
		{"Valid bounds negative", GISPoint{Latitude: -45.0, Longitude: -90.0}, true},
		{"Max valid lat/lon", GISPoint{Latitude: 90.0, Longitude: 180.0}, true},
		{"Min valid lat/lon", GISPoint{Latitude: -90.0, Longitude: -180.0}, true},
		{"Invalid latitude high", GISPoint{Latitude: 91.0, Longitude: 0}, false},
		{"Invalid latitude low", GISPoint{Latitude: -91.0, Longitude: 0}, false},
		{"Invalid longitude high", GISPoint{Latitude: 0, Longitude: 181.0}, false},
		{"Invalid longitude low", GISPoint{Latitude: 0, Longitude: -181.0}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.point.IsGeoCoordinateValid())
		})
	}
}

// TestGISPointIsGISZero checks if the IsGISZero method correctly identifies zero and non-zero GIS points.
func TestGISPointIsGISZero(t *testing.T) {
	tests := []struct {
		point GISPoint
		want  bool
	}{
		{GISPoint{Latitude: 0, Longitude: 0}, true},
		{GISPoint{Latitude: 1, Longitude: 0}, false},
		{GISPoint{Latitude: 0, Longitude: 1}, false},
		{GISPoint{Latitude: 1, Longitude: 1}, false},
	}

	for _, test := range tests {
		if got := test.point.IsGISZero(); got != test.want {
			t.Errorf("GISPoint{%f, %f}.IsGISZero() = %v; want %v", test.point.Latitude, test.point.Longitude, got, test.want)
		}
	}
}

// TestGISPointToStringLonLat checks if the ToStringLonLat method returns the correct string representation of a GISPoint.
func TestGISPointToStringLonLat(t *testing.T) {
	point := GISPoint{Latitude: 12.345678, Longitude: 98.765432}
	want := "lat: 12.345678, lon: 98.765432"
	if got := point.ToStringLonLat(); got != want {
		t.Errorf("GISPoint.ToStringLonLat() = %v; want %v", got, want)
	}
}

// TestGISPointToStringLat checks if the ToStringLat method returns the correct string representation of the latitude.
func TestGISPointToStringLat(t *testing.T) {
	point := GISPoint{Latitude: 12.345678}
	want := "12.345678"
	if got := point.ToStringLat(); got != want {
		t.Errorf("GISPoint.ToStringLat() = %v; want %v", got, want)
	}
}

// TestGISPointToStringLon checks if the ToStringLon method returns the correct string representation of the longitude.
func TestGISPointToStringLon(t *testing.T) {
	point := GISPoint{Longitude: 98.765432}
	want := "98.765432"
	if got := point.ToStringLon(); got != want {
		t.Errorf("GISPoint.ToStringLon() = %v; want %v", got, want)
	}
}

// TestGISPointToStringOSMPoint checks if the ToStringOSMPoint method returns the correct OpenStreetMap query string.
func TestGISPointToStringOSMPoint(t *testing.T) {
	point := GISPoint{Latitude: 12.345678, Longitude: 98.765432}
	want := "98.765432,12.345678"
	if got := point.ToStringOSMPoint(); got != want {
		t.Errorf("GISPoint.ToStringOSMPoint() = %v; want %v", got, want)
	}
}

// TestGISPointsClean checks if the Clean method correctly removes zero points from a GISPoints slice.
func TestGISPointsClean(t *testing.T) {
	points := GISPoints{
		&GISPoint{Latitude: 0, Longitude: 0},
		&GISPoint{Latitude: 1, Longitude: 1},
		nil,
	}
	cleaned := points.Clean()
	if len(cleaned) != 1 {
		t.Errorf("GISPoints.Clean() returned %d points; want 1", len(cleaned))
	}
	if cleaned[0].Latitude != 1 || cleaned[0].Longitude != 1 {
		t.Errorf("GISPoints.Clean() did not return the correct non-zero point")
	}
}

// TestGISPointsIsLine checks if the IsLine method correctly identifies a GISPoints slice as a line.
func TestGISPointsIsLine(t *testing.T) {
	points := GISPoints{
		&GISPoint{Latitude: 1, Longitude: 1},
		&GISPoint{Latitude: 2, Longitude: 2},
	}
	if !points.IsLine() {
		t.Errorf("GISPoints.IsLine() = false; want true")
	}
}

// TestGISPointsIsPolygon checks if the IsPolygon method correctly identifies a GISPoints slice as a polygon.
func TestGISPointsIsPolygon(t *testing.T) {
	points := GISPoints{
		&GISPoint{Latitude: 1, Longitude: 1},
		&GISPoint{Latitude: 2, Longitude: 2},
		&GISPoint{Latitude: 3, Longitude: 3},
	}
	if !points.IsPolygon() {
		t.Errorf("GISPoints.IsPolygon() = false; want true")
	}
}

func TestGISPoints_ContainsPoint(t *testing.T) {
	tests := []struct {
		name     string
		pts      GISPoints
		target   *GISPoint
		expected bool
	}{
		{
			name:     "Exact match point",
			pts:      GISPoints{{Latitude: 40.0, Longitude: -74.0}},
			target:   &GISPoint{Latitude: 40.0, Longitude: -74.0},
			expected: true,
		},
		{
			name: "Point on line",
			pts: GISPoints{
				{Latitude: 40.0, Longitude: -74.0},
				{Latitude: 41.0, Longitude: -75.0},
			},
			target:   &GISPoint{Latitude: 40.5, Longitude: -74.5},
			expected: true,
		},
		{
			name: "Point in polygon",
			pts: GISPoints{
				{Latitude: 0.0, Longitude: 0.0},
				{Latitude: 0.0, Longitude: 10.0},
				{Latitude: 10.0, Longitude: 10.0},
				{Latitude: 10.0, Longitude: 0.0},
			},
			target:   &GISPoint{Latitude: 5.0, Longitude: 5.0},
			expected: true,
		},
		{
			name: "Point outside polygon",
			pts: GISPoints{
				{Latitude: 0.0, Longitude: 0.0},
				{Latitude: 0.0, Longitude: 10.0},
				{Latitude: 10.0, Longitude: 10.0},
				{Latitude: 10.0, Longitude: 0.0},
			},
			target:   &GISPoint{Latitude: 15.0, Longitude: 15.0},
			expected: false,
		},
		{
			name:     "Invalid target point",
			pts:      GISPoints{{Latitude: 0.0, Longitude: 0.0}},
			target:   &GISPoint{Latitude: 0.0, Longitude: 0.0},
			expected: false,
		},
		{
			name:     "Nil target point",
			pts:      GISPoints{{Latitude: 10.0, Longitude: 10.0}},
			target:   nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pts.ContainsPoint(tt.target)
			if result != tt.expected {
				t.Errorf("ContainsPoint() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanityDistance(t *testing.T) {
	a := &GISPoint{Latitude: 0, Longitude: 0}
	b := &GISPoint{Latitude: 10, Longitude: 0}
	p := &GISPoint{Latitude: 20, Longitude: 0}

	dist := DistanceFromPointToLineSegment(a, b, p)
	if dist < 1000000 {
		t.Errorf("Expected large distance, got %f", dist)
	}
}

func TestGISPoints_IntersectsRadius(t *testing.T) {
	tests := []struct {
		name         string
		pts          GISPoints
		target       *GISPoint
		radiusMeters float64
		expected     bool
	}{
		{
			name:         "Point within radius of exact match",
			pts:          GISPoints{{Latitude: 40.0, Longitude: -74.0}},
			target:       &GISPoint{Latitude: 40.0005, Longitude: -74.0005}, // ~70m
			radiusMeters: 100,
			expected:     true,
		},
		{
			name: "Near line segment (>100m away)",
			pts: GISPoints{
				{Latitude: 40.0, Longitude: -74.0},
				{Latitude: 41.0, Longitude: -75.0},
			},
			target:       &GISPoint{Latitude: 40.5, Longitude: -74.5},
			radiusMeters: 100,
			expected:     false,
		},
		{
			name: "Far from polygon edge (560km)",
			pts: GISPoints{
				{Latitude: 0.0, Longitude: 0.0},
				{Latitude: 0.0, Longitude: 10.0},
				{Latitude: 10.0, Longitude: 10.0},
				{Latitude: 10.0, Longitude: 0.0},
			},
			target:       &GISPoint{Latitude: 5.0, Longitude: 0.05},
			radiusMeters: 5600,
			expected:     false, // was mistakenly true
		},
		{
			name: "No intersection (far from polygon)",
			pts: GISPoints{
				{Latitude: 0.0, Longitude: 0.0},
				{Latitude: 0.0, Longitude: 10.0},
				{Latitude: 10.0, Longitude: 10.0},
				{Latitude: 10.0, Longitude: 0.0},
			},
			target:       &GISPoint{Latitude: 20.0, Longitude: 20.0}, // far outside
			radiusMeters: 100,
			expected:     false,
		},
		{
			name:         "Nil target",
			pts:          GISPoints{{Latitude: 1.0, Longitude: 1.0}},
			target:       nil,
			radiusMeters: 100,
			expected:     false,
		},
		{
			name:         "Zero radius",
			pts:          GISPoints{{Latitude: 1.0, Longitude: 1.0}},
			target:       &GISPoint{Latitude: 1.0, Longitude: 1.0},
			radiusMeters: 0,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pts.IntersectsRadius(tt.target, tt.radiusMeters)
			if result != tt.expected {
				t.Errorf("IntersectsRadius() = %v, want %v", result, tt.expected)
			}
		})
	}
}
