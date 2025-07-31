package ageo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeoFilter_requiresCheck(t *testing.T) {
	tests := []struct {
		name string
		gf   GeoFilter
		want bool
	}{
		{"Empty", GeoFilter{}, false},
		{"Countries only", GeoFilter{Countries: []string{"US"}}, true},
		{"Regions only", GeoFilter{Regions: []string{"CA"}}, true},
		{"Cities only", GeoFilter{Cities: []string{"Paris"}}, true},
		{"Polygon only", GeoFilter{GISPolygon: polygonSquare()}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.gf.requiresCheck())
		})
	}
}

func TestGeoFilter_Matches(t *testing.T) {
	pointInside := GeoInfo{
		CountryCode: " us ",
		Region:      "  ny ",
		City:        " new york ",
		GISPoint:    GISPoint{Latitude: 5.0, Longitude: 5.0},
	}

	tests := []struct {
		name     string
		filter   GeoFilter
		geo      GeoInfo
		expected bool
	}{
		{"Exact country match", GeoFilter{Countries: []string{"us"}}, pointInside, true},
		{"Case/trim-insensitive region", GeoFilter{Regions: []string{"NY"}}, pointInside, true},
		{"City match insensitive", GeoFilter{Cities: []string{"NEW YORK"}}, pointInside, true},
		{"Polygon match", GeoFilter{GISPolygon: polygonSquare()}, pointInside, true},

		{"Country mismatch", GeoFilter{Countries: []string{"de"}}, pointInside, false},
		{"Region mismatch", GeoFilter{Regions: []string{"CA"}}, pointInside, false},
		{"City mismatch", GeoFilter{Cities: []string{"London"}}, pointInside, false},
		{"Outside polygon", GeoFilter{GISPolygon: polygonSquare()}, GeoInfo{
			CountryCode: "us", Region: "ny", City: "new york",
			GISPoint: GISPoint{Latitude: 20.0, Longitude: 20.0},
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.Matches(tt.geo)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGeoFilter_GeoCheck(t *testing.T) {
	point := GeoInfo{
		CountryCode: "us",
		Region:      "ny",
		City:        "new york",
		GISPoint:    GISPoint{Latitude: 5.0, Longitude: 5.0},
	}

	tests := []struct {
		name     string
		filter   GeoFilter
		geo      GeoInfo
		expected bool
	}{
		{"Allow country match", GeoFilter{Countries: []string{"us"}}, point, true},
		{"Deny country match", GeoFilter{IsDeny: true, Countries: []string{"us"}}, point, false},
		{"Deny not matched → allow", GeoFilter{IsDeny: true, Countries: []string{"ca"}}, point, true},
		{"Allow not matched → deny", GeoFilter{IsDeny: false, Countries: []string{"ca"}}, point, false},
		{"Empty filter → allow", GeoFilter{}, point, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.filter.GeoCheck(tt.geo)
			assert.Equal(t, tt.expected, got)
		})
	}
}

// polygonSquare returns a 10x10 square polygon covering point (5,5)
func polygonSquare() GISPoints {
	return GISPoints{
		{Latitude: 0.0, Longitude: 0.0},
		{Latitude: 0.0, Longitude: 10.0},
		{Latitude: 10.0, Longitude: 10.0},
		{Latitude: 10.0, Longitude: 0.0},
	}
}

func TestGeoFilters_Evaluate(t *testing.T) {
	point := GeoInfo{
		CountryCode: "us", Region: "ny", City: "new york",
		GISPoint: GISPoint{Latitude: 5.0, Longitude: 5.0},
	}

	allow := &GeoFilter{Countries: []string{"us"}}
	deny := &GeoFilter{IsDeny: true, Cities: []string{"new york"}}
	denyNotMatched := &GeoFilter{IsDeny: true, Countries: []string{"de"}}

	tests := []struct {
		name     string
		rules    GeoFilters
		expected bool
	}{
		{"Allow wins", GeoFilters{allow}, true},
		{"Deny matched", GeoFilters{allow, deny}, false},
		{"Deny not matched", GeoFilters{allow, denyNotMatched}, true},
		{"Deny exists but unmatched", GeoFilters{denyNotMatched}, false},
		{"No rules", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rules.Evaluate(point)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGeoFilters_Evaluate_ExampleBehavior(t *testing.T) {
	geoUS_NYC := GeoInfo{
		CountryCode: "US",
		City:        "New York",
		GISPoint:    GISPoint{Latitude: 40.7128, Longitude: -74.0060},
	}

	geoUS_LA := GeoInfo{
		CountryCode: "US",
		City:        "Los Angeles",
		GISPoint:    GISPoint{Latitude: 34.0522, Longitude: -118.2437},
	}

	geoCanada := GeoInfo{
		CountryCode: "CA",
		City:        "Toronto",
		GISPoint:    GISPoint{Latitude: 43.6510, Longitude: -79.3470},
	}

	filters := GeoFilters{
		&GeoFilter{Countries: []string{"US"}},                  // Allow US
		&GeoFilter{IsDeny: true, Cities: []string{"New York"}}, // Deny NYC
	}

	tests := []struct {
		name     string
		geo      GeoInfo
		expected bool
	}{
		{"US and NYC → DENIED", geoUS_NYC, false},
		{"US and LA → ALLOWED", geoUS_LA, true},
		{"Canada → DENIED", geoCanada, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filters.Evaluate(tt.geo)
			if got != tt.expected {
				t.Errorf("Evaluate() = %v, want %v", got, tt.expected)
			}
		})
	}
}
