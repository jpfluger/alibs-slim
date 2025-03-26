package autils

import (
	"testing"
)

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
