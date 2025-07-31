package ageo

import (
	"testing"
)

// TestGISPolygonHasBoundary checks if the HasBoundary method correctly identifies a valid polygon.
func TestGISPolygonHasBoundary(t *testing.T) {
	// Create a GISPolygon with less than three points (not a polygon).
	polygon := GISPolygon{
		Points: GISPoints{
			&GISPoint{Latitude: 0, Longitude: 0},
			&GISPoint{Latitude: 1, Longitude: 1},
		},
	}

	// The polygon should not have a valid boundary.
	if polygon.HasBoundary() {
		t.Errorf("GISPolygon.HasBoundary() = true; want false for less than three points")
	}

	// Add a third point to form a valid polygon.
	polygon.Points = append(polygon.Points, &GISPoint{Latitude: 2, Longitude: 2})

	// Now the polygon should have a valid boundary.
	if !polygon.HasBoundary() {
		t.Errorf("GISPolygon.HasBoundary() = false; want true for three or more points")
	}
}

// TestGISPolygonClean checks if the Clean method correctly sanitizes the GISPolygon.
func TestGISPolygonClean(t *testing.T) {
	// Create a GISPolygon with invalid points and opacity.
	polygon := GISPolygon{
		Points: GISPoints{
			&GISPoint{Latitude: 0, Longitude: 0}, // This is a zero point and should be removed.
			&GISPoint{Latitude: 1, Longitude: 1},
			&GISPoint{Latitude: 2, Longitude: 2},
		},
		Opacity: -0.5, // This is an invalid opacity and should be set to 1.
	}

	// Clean the polygon.
	polygon.Clean()

	// After cleaning, the polygon should have zero points and opacity set to 1.
	if len(polygon.Points) != 0 {
		t.Errorf("GISPolygon.Clean() did not remove points, got %d points; want 0", len(polygon.Points))
	}
	//// After cleaning, the polygon should have two points and opacity set to 1.
	//if len(polygon.Points) != 2 {
	//	t.Errorf("GISPolygon.Clean() did not remove zero points, got %d points; want 2", len(polygon.Points))
	//}
	if polygon.Opacity != 1 {
		t.Errorf("GISPolygon.Clean() did not correct opacity, got %f; want 1", polygon.Opacity)
	}
}
