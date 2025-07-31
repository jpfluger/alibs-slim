package ageo

import (
	"testing"
)

// TestGISBoundsHasBounds checks if the HasBounds method correctly identifies valid and invalid bounds.
func TestGISBoundsHasBounds(t *testing.T) {
	tests := []struct {
		name     string
		bounds   GISBounds
		expected bool
	}{
		{
			name: "Valid bounds",
			bounds: GISBounds{
				SW: GISPoint{Latitude: -90, Longitude: -180},
				NE: GISPoint{Latitude: 90, Longitude: 180},
			},
			expected: true,
		},
		{
			name: "Invalid bounds with zero southwest point",
			bounds: GISBounds{
				SW: GISPoint{Latitude: 0, Longitude: 0},
				NE: GISPoint{Latitude: 90, Longitude: 180},
			},
			expected: false,
		},
		{
			name: "Invalid bounds with zero northeast point",
			bounds: GISBounds{
				SW: GISPoint{Latitude: -90, Longitude: -180},
				NE: GISPoint{Latitude: 0, Longitude: 0},
			},
			expected: false,
		},
		{
			name:     "Nil bounds",
			bounds:   GISBounds{},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.bounds.HasBounds(); got != test.expected {
				t.Errorf("GISBounds.HasBounds() = %v, want %v", got, test.expected)
			}
		})
	}
}

// TestGISBounds_ContainsPoint tests the ContainsPoint method for GISBounds.
func TestGISBounds_ContainsPoint(t *testing.T) {
	bounds := &GISBounds{
		SW: GISPoint{Latitude: 10.0, Longitude: 10.0},
		NE: GISPoint{Latitude: 20.0, Longitude: 20.0},
	}

	tests := []struct {
		name     string
		point    *GISPoint
		expected bool
	}{
		{"Inside bounds", &GISPoint{Latitude: 15.0, Longitude: 15.0}, true},
		{"On SW corner", &GISPoint{Latitude: 10.0, Longitude: 10.0}, true},
		{"On NE corner", &GISPoint{Latitude: 20.0, Longitude: 20.0}, true},
		{"Outside north", &GISPoint{Latitude: 21.0, Longitude: 15.0}, false},
		{"Outside west", &GISPoint{Latitude: 15.0, Longitude: 9.0}, false},
		{"Nil point", nil, false},
		{"Invalid point", &GISPoint{Latitude: 0, Longitude: 0}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := bounds.ContainsPoint(tc.point)
			if got != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}

// TestGISBounds_IntersectsRadius tests the IntersectsRadius method for GISBounds.
func TestGISBounds_IntersectsRadius(t *testing.T) {
	bounds := &GISBounds{
		SW: GISPoint{Latitude: 10.0, Longitude: 10.0},
		NE: GISPoint{Latitude: 20.0, Longitude: 20.0},
	}

	tests := []struct {
		name         string
		point        *GISPoint
		radiusMeters float64
		expected     bool
	}{
		{"Inside bounds", &GISPoint{Latitude: 15.0, Longitude: 15.0}, 500, true},
		{"Near edge", &GISPoint{Latitude: 9.95, Longitude: 15.0}, 6000, true},    // ~5km from southern edge
		{"Near corner", &GISPoint{Latitude: 9.95, Longitude: 9.95}, 10000, true}, // ~7km from SW corner
		{"Far outside", &GISPoint{Latitude: 0, Longitude: 0}, 10000, false},
		{"Zero radius", &GISPoint{Latitude: 15.0, Longitude: 15.0}, 0, false},
		{"Invalid point", &GISPoint{Latitude: 0, Longitude: 0}, 10000, false},
		{"Nil point", nil, 10000, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := bounds.IntersectsRadius(tc.point, tc.radiusMeters)
			if got != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}
