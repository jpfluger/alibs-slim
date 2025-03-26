package acontact

// ILocation defines an interface for location-related operations.
type ILocation interface {
	// FindCountryName takes a target string, which could be a country code or name,
	// and returns the full name of the country.
	FindCountryName(target string) string

	// FindCountryShort takes a target string, which could be a country code or name,
	// and returns the abbreviated form or short code of the country.
	FindCountryShort(target string) string

	// FindDivisionName takes a countryKey, which could be a country code or name,
	// a division string which could be a state or province, and a getShort boolean
	// indicating whether to return the short form (true) or full name (false) of the division.
	FindDivisionName(countryKey string, division string, getShort bool) string
}
