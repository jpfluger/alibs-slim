package autils

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
