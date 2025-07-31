package ageo

import (
	"math"
	"testing"
)

func TestGeoInfo_DistanceTo(t *testing.T) {
	tests := []struct {
		name     string
		g1       GeoInfo
		g2       GeoInfo
		expected float64 // in meters
		valid    bool
	}{
		{
			name:     "Same location",
			g1:       GeoInfo{City: "Same", GISPoint: GISPoint{Latitude: 40.0, Longitude: -74.0}},
			g2:       GeoInfo{City: "Same", GISPoint: GISPoint{Latitude: 40.0, Longitude: -74.0}},
			expected: 0,
			valid:    true,
		},
		{
			name:     "New York to Los Angeles",
			g1:       GeoInfo{City: "NYC", GISPoint: GISPoint{Latitude: 40.7128, Longitude: -74.0060}},
			g2:       GeoInfo{City: "LA", GISPoint: GISPoint{Latitude: 34.0522, Longitude: -118.2437}},
			expected: 3936000, // ~3936 km
			valid:    true,
		},
		{
			name:     "Invalid source point (0,0)",
			g1:       GeoInfo{City: "Invalid", GISPoint: GISPoint{Latitude: 0, Longitude: 0}},
			g2:       GeoInfo{City: "Valid", GISPoint: GISPoint{Latitude: 40.0, Longitude: -74.0}},
			expected: -1,
			valid:    false,
		},
		{
			name:     "Invalid destination point (0,0)",
			g1:       GeoInfo{City: "Valid", GISPoint: GISPoint{Latitude: 40.0, Longitude: -74.0}},
			g2:       GeoInfo{City: "Invalid", GISPoint: GISPoint{Latitude: 0, Longitude: 0}},
			expected: -1,
			valid:    false,
		},
	}

	const tolerance = 25000.0 // 25 km

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.g1.DistanceTo(tt.g2)
			if tt.valid {
				if math.Abs(got-tt.expected) > tolerance {
					t.Errorf("DistanceTo() = %.2f, want %.2f ±%.2f", got, tt.expected, tolerance)
				}
			} else {
				if got != -1 {
					t.Errorf("DistanceTo() = %.2f, want -1 for invalid input", got)
				}
			}
		})
	}
}

func TestGeoInfo_DistanceToKM(t *testing.T) {
	tests := []struct {
		name     string
		g1       GeoInfo
		g2       GeoInfo
		expected float64 // in km
		valid    bool
	}{
		{
			name:     "London to Paris",
			g1:       GeoInfo{City: "London", GISPoint: GISPoint{Latitude: 51.5074, Longitude: -0.1278}},
			g2:       GeoInfo{City: "Paris", GISPoint: GISPoint{Latitude: 48.8566, Longitude: 2.3522}},
			expected: 344,
			valid:    true,
		},
		{
			name:     "Invalid both points",
			g1:       GeoInfo{City: "", GISPoint: GISPoint{Latitude: 0, Longitude: 0}},
			g2:       GeoInfo{City: "", GISPoint: GISPoint{Latitude: 0, Longitude: 0}},
			expected: -1,
			valid:    false,
		},
	}

	const tolerance = 10.0 // 10 km

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.g1.DistanceToKM(tt.g2)
			if tt.valid {
				if math.Abs(got-tt.expected) > tolerance {
					t.Errorf("DistanceToKM() = %.2f, want %.2f ±%.2f", got, tt.expected, tolerance)
				}
			} else {
				if got != -1 {
					t.Errorf("DistanceToKM() = %.2f, want -1 for invalid input", got)
				}
			}
		})
	}
}
