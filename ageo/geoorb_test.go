package ageo

import (
	"github.com/paulmach/orb"
	"math"
	"testing"
)

// almostEqualWithEpsilon checks floating-point equality within a small tolerance.
func almostEqualWithEpsilon(a, b, epsilon float64) bool {
	return math.Abs(a-b) <= epsilon
}

func TestHaversineDistance(t *testing.T) {
	a := &GISPoint{Latitude: 40.748817, Longitude: -73.985428} // Empire State Building
	b := &GISPoint{Latitude: 40.689247, Longitude: -74.044502} // Statue of Liberty

	wantKM := 8.3
	wantM := wantKM * 1000

	gotKM := HaversineDistanceBetweenPointsKM(a, b)
	gotM := HaversineDistanceBetweenPointsM(a, b)

	if !almostEqualWithEpsilon(gotKM, wantKM, 0.1) {
		t.Errorf("HaversineDistanceBetweenPointsKM = %v; want approx %v", gotKM, wantKM)
	}
	if !almostEqualWithEpsilon(gotM, wantM, 100) {
		t.Errorf("HaversineDistanceBetweenPointsM = %v; want approx %v", gotM, wantM)
	}
}

func TestDistanceFromPointToLineSegment(t *testing.T) {
	a := orb.Point{0, 0}
	b := orb.Point{2, 0}
	p := orb.Point{1, 1} // ~111 km north of the line

	distance := DistanceFromSegmentMeters(p, a, b)
	expected := 111000.0 // ~111 km

	if !almostEqualWithEpsilon(distance, expected, 2000) {
		t.Errorf("Expected distance ~%v meters, got %v", expected, distance)
	}
}

func TestIsPointInsidePolygon(t *testing.T) {
	polygon := GISPoints{
		{Latitude: 0, Longitude: 0},
		{Latitude: 0, Longitude: 2},
		{Latitude: 2, Longitude: 2},
		{Latitude: 2, Longitude: 0},
	}

	tests := []struct {
		name     string
		point    *GISPoint
		expected bool
	}{
		{"Inside", &GISPoint{Latitude: 1, Longitude: 1}, true},
		{"On edge", &GISPoint{Latitude: 0, Longitude: 1}, true},
		{"On vertex", &GISPoint{Latitude: 0, Longitude: 0}, true},
		{"Outside", &GISPoint{Latitude: 3, Longitude: 3}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsPointInsidePolygon(tc.point, polygon)
			if got != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestIsPointApproximatelyOnLine(t *testing.T) {
	a := &GISPoint{Latitude: 0, Longitude: 0}
	b := &GISPoint{Latitude: 0, Longitude: 2}

	tests := []struct {
		name   string
		p      *GISPoint
		expect bool
	}{
		{"On line", &GISPoint{Latitude: 0, Longitude: 1}, true},
		{"Near line", &GISPoint{Latitude: 0.00005, Longitude: 1}, true}, // ~5.5m north
		{"Off line", &GISPoint{Latitude: 1, Longitude: 1}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsPointApproximatelyOnLine(a, b, tc.p)
			if got != tc.expect {
				t.Errorf("Expected %v, got %v", tc.expect, got)
			}
		})
	}
}
