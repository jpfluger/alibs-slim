package ageo

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var geoDBPath = os.Getenv("MAXMIND_DB_PATH") // e.g. set in test env

func TestMain(m *testing.M) {
	if geoDBPath == "" {
		println("MAXMIND_DB_PATH env var is required for GeoIP tests")
		os.Exit(1)
	}

	err := InitGeoInfo(geoDBPath, false)
	if err != nil {
		println("Failed to initialize GeoIP database:", err.Error())
		os.Exit(1)
	}
	defer CloseGeoInfo()

	code := m.Run()
	os.Exit(code)
}

func TestLookupGeoInfoForValidPublicIP(t *testing.T) {
	ip := "8.8.8.8" // Google's DNS — should always have a match

	geo := LookupGeoInfoForIP(ip)
	assert.NotNil(t, geo, "GeoInfo should not be nil for a valid IP")
	assert.True(t, geo.IsValid(), "GeoInfo should be valid")
	assert.Equal(t, ip, geo.IPv4)
	assert.NotZero(t, geo.Latitude)
	assert.NotZero(t, geo.Longitude)
}

func TestLookupGeoInfoForInvalidIP(t *testing.T) {
	ip := "999.999.999.999"
	geo := LookupGeoInfoForIP(ip)
	assert.Nil(t, geo, "GeoInfo should be nil for invalid IP")
}

func TestLookupGeoInfoForEmptyIP(t *testing.T) {
	geo := LookupGeoInfoForIP("")
	assert.Nil(t, geo, "GeoInfo should be nil for empty IP")
}

func TestMustLookupGeoInfoForInvalidIP(t *testing.T) {
	geo := MustLookupGeoInfoForIP("invalid-ip")
	assert.NotNil(t, geo, "MustLookup should not return nil")
	assert.False(t, geo.IsValid(), "GeoInfo should not be valid for invalid IP")
}

func TestReloadDB(t *testing.T) {
	// Directly trigger reload for coverage
	CloseGeoInfo()
	err := InitGeoInfo(geoDBPath, false)
	assert.NoError(t, err, "should reload without error")

	geo := LookupGeoInfoForIP("1.1.1.1")
	assert.NotNil(t, geo, "GeoInfo should be valid after reload")
}

func TestGeoInfo_IsValid(t *testing.T) {
	valid := GeoInfo{
		City:   "San Francisco",
		Region: "CA",
		GISPoint: GISPoint{
			Latitude:  37.7749,
			Longitude: -122.4194,
		},
		IPv4: "8.8.8.8",
	}
	assert.True(t, valid.IsValid())

	invalid := GeoInfo{
		City:   "",
		Region: "",
		GISPoint: GISPoint{
			Latitude:  0,
			Longitude: 0,
		},
	}
	assert.False(t, invalid.IsValid())
}

func TestHaversineKM(t *testing.T) {
	const epsilonKM = 2.0 // acceptable error in kilometers
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64 // expected distance in km (approximate)
	}{
		{
			name:     "Same point",
			lat1:     0,
			lon1:     0,
			lat2:     0,
			lon2:     0,
			expected: 0,
		},
		{
			name:     "NYC to LA",
			lat1:     40.7128,
			lon1:     -74.0060,
			lat2:     34.0522,
			lon2:     -118.2437,
			expected: 3983, // approximate distance in km
		},
		{
			name:     "London to Paris",
			lat1:     51.5074,
			lon1:     -0.1278,
			lat2:     48.8566,
			lon2:     2.3522,
			expected: 344,
		},
		{
			name:     "Tokyo to Sydney",
			lat1:     35.6895,
			lon1:     139.6917,
			lat2:     -33.8688,
			lon2:     151.2093,
			expected: 7848,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dist := HaversineDistanceKM(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if !almostEqualWithEpsilon(dist, tt.expected, epsilonKM) {
				t.Errorf("%s: got %.2f km, expected %.2f ± %.2f", tt.name, dist, tt.expected, epsilonKM)
			}
		})
	}
}
