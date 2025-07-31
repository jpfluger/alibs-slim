package ageo

import "github.com/jpfluger/alibs-slim/autils"

// GeoFilter defines a geospatial match rule that restricts access based on location data.
// If IsDeny is true, the match will *reject* access instead of allowing it.
type GeoFilter struct {
	IsDeny bool `json:"is_deny,omitempty"`

	// List-based geographic constraints
	Countries []string `json:"countries,omitempty"` // ISO 3166-1 alpha-2 country codes (e.g., "US", "DE")
	Regions   []string `json:"regions,omitempty"`   // Region/state/province names
	Cities    []string `json:"cities,omitempty"`    // City names (case-sensitive)

	// GISPolygon defines a geographic boundary.
	// If specified, the user’s location must fall inside or near the polygon.
	GISPolygon GISPoints `json:"gis_points,omitempty"`
}

// Clone returns a deep copy of the GeoFilter.
func (gf *GeoFilter) Clone() *GeoFilter {
	if gf == nil {
		return nil
	}
	clone := &GeoFilter{
		IsDeny:     gf.IsDeny,
		Countries:  append([]string{}, gf.Countries...),
		Regions:    append([]string{}, gf.Regions...),
		Cities:     append([]string{}, gf.Cities...),
		GISPolygon: gf.GISPolygon.Clone(),
	}
	return clone
}

// requiresCheck returns true if the GeoFilter struct contains any filtering criteria.
// This includes countries, regions, cities, or GIS polygon geometry.
// Used to avoid unnecessary checks when no geo constraints are defined.
func (gf *GeoFilter) requiresCheck() bool {
	return len(gf.Countries) > 0 || len(gf.Regions) > 0 || len(gf.Cities) > 0 || gf.GISPolygon.IsMultiPoint()
}

// GeoCheck evaluates whether a given GeoInfo matches the filtering criteria.
// If a match is found:
//   - returns false if IsDeny is true (explicit deny)
//   - returns true if IsDeny is false (explicit allow)
//
// If no match is found:
//   - returns true if IsDeny is false (default allow)
//   - returns false if IsDeny is true (default deny)
func (gf *GeoFilter) GeoCheck(geo GeoInfo) bool {
	// If a filter requires checking but doesn't match, the check fails.
	if gf.requiresCheck() {
		matched := gf.Matches(geo)
		if matched {
			return !gf.IsDeny // Match found: respect allow/deny
		}
		return gf.IsDeny // Didn't match: deny if it's a deny rule
	}
	// No filtering criteria → default to allow
	return true
}

// Matches determines whether the given GeoInfo matches all specified filtering criteria:
// country, region, city, and optional GIS polygon proximity.
// Note: This method does not apply IsDeny logic — use GeoCheck for enforcement.
func (gf *GeoFilter) Matches(geo GeoInfo) bool {
	if !geo.IsValid() {
		return false
	}

	// Countries
	if len(gf.Countries) > 0 {
		found := false
		gCountry := autils.ToStringTrimLower(geo.CountryCode)
		for _, c := range gf.Countries {
			if gCountry == autils.ToStringTrimLower(c) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Regions
	if len(gf.Regions) > 0 {
		found := false
		gRegion := autils.ToStringTrimLower(geo.Region)
		for _, r := range gf.Regions {
			if gRegion == autils.ToStringTrimLower(r) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Cities
	if len(gf.Cities) > 0 {
		found := false
		gCity := autils.ToStringTrimLower(geo.City)
		for _, city := range gf.Cities {
			if gCity == autils.ToStringTrimLower(city) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// GIS Polygon
	if len(gf.GISPolygon) > 0 {
		if !gf.GISPolygon.IntersectsRadius(&geo.GISPoint, DefaultGeoFencingMeters) {
			return false
		}
	}

	return true
}

// GeoFilters is a slice of GeoFilter rules.
type GeoFilters []*GeoFilter

// Evaluate applies a sequence of GeoFilter rules to the given GeoInfo.
//
// Evaluation logic prioritizes denial:
//   - If any GeoFilter with IsDeny=true matches the input, access is denied immediately.
//   - Otherwise, if at least one GeoFilter with IsDeny=false matches, access is allowed.
//   - If no filters match:
//   - If no filters are defined (empty GeoFilters), access is allowed (default allow).
//   - If one or more filters exist but none match, access is denied (default deny).
//
// This function allows composite geo rulesets, such as:
//   - Allow if in US but deny if in specific cities.
//   - Allow access near a polygon, but deny certain countries or regions.
//
// Example behavior:
//
//	GeoFilters{
//	    &GeoFilter{Countries: []string{"US"}},                // allow if in US
//	    &GeoFilter{IsDeny: true, Cities: []string{"NYC"}},    // deny if in NYC
//	}
//
//	• geo = US, city = NYC → DENIED
//	• geo = US, city = LA  → ALLOWED
//	• geo = CA             → DENIED
func (gfs GeoFilters) Evaluate(geo GeoInfo) bool {
	allowMatched := false

	for _, gf := range gfs {
		if gf == nil || !gf.requiresCheck() {
			continue
		}
		if gf.Matches(geo) {
			if gf.IsDeny {
				return false // deny matched
			}
			allowMatched = true // allow matched, but keep checking for denies
		}
	}

	if allowMatched {
		return true // at least one non-deny matched, and no deny blocked it
	}

	// No filters matched: deny by default if any filters exist
	return len(gfs) == 0
}

// Clone returns a deep copy of the GeoFilters slice.
func (gfs GeoFilters) Clone() GeoFilters {
	if gfs == nil || len(gfs) == 0 {
		return GeoFilters{}
	}
	cloned := make(GeoFilters, 0, len(gfs))
	for _, gf := range gfs {
		if gf == nil {
			cloned = append(cloned, nil)
		} else {
			cloned = append(cloned, gf.Clone())
		}
	}
	return cloned
}
