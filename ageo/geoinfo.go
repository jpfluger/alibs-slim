package ageo

type GeoInfo struct {
	CountryCode string `json:"country_code"`
	IsEU        bool   `json:"isEU,omitempty"`
	Region      string `json:"region"`
	City        string `json:"city"`
	GISPoint
	IPv4 string `json:"ipv4"`
}

// IsValid returns true if the GeoInfo contains usable coordinates.
// City/Region may be empty depending on IP resolution.
func (g GeoInfo) IsValid() bool {
	return g.GISPoint.IsPracticallyValid()
}

// DistanceTo returns the Haversine distance in meters between this GeoInfo's location
// and another GeoInfo's location. Returns -1 if either point is invalid.
func (g GeoInfo) DistanceTo(other GeoInfo) float64 {
	if !g.IsValid() || !other.IsValid() {
		return -1
	}
	return HaversineDistanceBetweenPointsM(&g.GISPoint, &other.GISPoint)
}

// DistanceToKM returns the Haversine distance in kilometers between this GeoInfo's location
// and another GeoInfo's location. Returns -1 if either point is invalid.
func (g GeoInfo) DistanceToKM(other GeoInfo) float64 {
	if !g.IsValid() || !other.IsValid() {
		return -1
	}
	return HaversineDistanceBetweenPointsKM(&g.GISPoint, &other.GISPoint)
}
