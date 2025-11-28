package acontact

import (
	"fmt"
	"strings"

	"github.com/bojanz/address"
)

// MailAddress represents a physical mailing address.
// It uses the bojanz/address library to handle address formatting and validation.
type MailAddress struct {
	address.Address // Embedded Address struct from the bojanz/address library

	Raw        string `json:"raw,omitempty"`        // Raw address data, used if parsing fails
	IsVerified bool   `json:"isVerified,omitempty"` // Indicates if the address has been verified
}

// Validate checks if either the structured address or the raw address is provided and verifies the format.
func (add *MailAddress) Validate() error {
	// Check if both structured address and raw address are empty or both are provided, which is not allowed.
	if add.IsEmpty() && add.Raw == "" {
		return fmt.Errorf("either structured address or raw address must be provided")
	}
	if !add.IsEmpty() && add.Raw != "" {
		return fmt.Errorf("only one of structured address or raw address should be provided")
	}

	// If the structured address is provided, validate it.
	if !add.IsEmpty() {
		// Trim whitespace and validate required fields.
		add.Line1 = strings.TrimSpace(add.Line1)
		add.CountryCode = strings.TrimSpace(add.CountryCode)
		if add.Line1 == "" || add.CountryCode == "" {
			return fmt.Errorf("structured address is incomplete")
		}
		// Additional validation logic can be implemented here.
		add.IsVerified = true // Structured address is considered verified if it passes validation.
	} else {
		// If only the raw address is provided, store it and mark the address as not verified.
		add.Raw = strings.TrimSpace(add.Raw)
		add.IsVerified = false
	}
	return nil
}

// GetVerify returns the verification status of the address.
func (add *MailAddress) GetVerify() bool {
	return add.IsVerified
}

// VerifyWithFields checks if the specified fields are present and sets the verification status.
func (add *MailAddress) VerifyWithFields(langType string, reqFields ...address.Field) error {
	// Set the locale based on the provided language type.
	locale := address.NewLocale(langType)
	if locale.IsEmpty() {
		locale = address.NewLocale("en") // Default to English if the locale is empty.
	}

	// Reset verification status.
	add.IsVerified = false

	// Check each required field.
	for _, field := range reqFields {
		value := getField(add, field)
		if strings.TrimSpace(value) == "" {
			add.IsVerified = false
			return fmt.Errorf("%s is required", fieldDescription(field))
		}
	}

	// If all required fields are present, set the address as verified.
	add.IsVerified = true
	return nil
}

// getField returns the value of the specified field from the address.
func getField(add *MailAddress, field address.Field) string {
	switch field {
	case address.FieldLine1:
		return add.Line1
	case address.FieldLine2:
		return add.Line2
	case address.FieldLine3:
		return add.Line3
	case address.FieldSublocality:
		return add.Sublocality
	case address.FieldLocality:
		return add.Locality
	case address.FieldRegion:
		return add.Region
	case address.FieldPostalCode:
		return add.PostalCode
	default:
		return ""
	}
}

// fieldDescription provides a human-readable description of the address field.
func fieldDescription(field address.Field) string {
	switch field {
	case address.FieldLine1:
		return "address line 1"
	case address.FieldLine2:
		return "address line 2"
	case address.FieldLine3:
		return "address line 3"
	case address.FieldSublocality:
		return "sublocality (neighborhood/suburb/district)"
	case address.FieldLocality:
		return "locality (city/village/town)"
	case address.FieldRegion:
		return "region (state/province/prefecture)"
	case address.FieldPostalCode:
		return "postal code (zip/pin code)"
	default:
		return "unknown field"
	}
}

// IsEmpty checks if all the address fields are empty.
func (add *MailAddress) IsEmpty() bool {
	return add.Line1 == "" &&
		add.Line2 == "" &&
		add.Line3 == "" &&
		add.Sublocality == "" &&
		add.Locality == "" &&
		add.Region == "" &&
		add.PostalCode == "" &&
		add.CountryCode == ""
}

// ToHTML formats the address into an HTML string based on the specified language type.
func (add *MailAddress) ToHTML(langType string) string {
	// Set the locale based on the provided language type.
	locale := address.NewLocale(langType)
	if locale.IsEmpty() {
		locale = address.NewLocale("en") // Default to English if the locale is empty.
	}

	// Use the address formatter to generate the HTML representation of the address.
	formatter := address.NewFormatter(locale)
	return formatter.Format(add.Address)
}

func (add *MailAddress) ToLines() []string {
	if add == nil {
		return nil
	}

	var lines []string
	appendIf := func(s string) {
		s = strings.TrimSpace(s)
		if s != "" {
			lines = append(lines, s)
		}
	}

	appendIf(add.Line1)
	appendIf(add.Line2)
	appendIf(add.Line3)

	var cityParts []string
	for _, s := range []string{add.Locality, add.Region, add.PostalCode} {
		if strings.TrimSpace(s) != "" {
			cityParts = append(cityParts, strings.TrimSpace(s))
		}
	}
	if len(cityParts) > 0 {
		lines = append(lines, strings.Join(cityParts, ", "))
	}

	appendIf(add.CountryCode)

	if len(lines) == 0 && add.Raw != "" {
		lines = append(lines, strings.TrimSpace(add.Raw))
	}

	return lines
}

// ToSingleLineLocalRegionPostal collapses the locality, region and postal into a single line.
func (add *MailAddress) ToSingleLineLocalRegionPostal() string {
	if add == nil {
		return ""
	}

	var sb strings.Builder

	if add.Locality != "" {
		sb.WriteString(strings.TrimSpace(add.Locality))
	}

	if add.Region != "" {
		if sb.Len() > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(strings.TrimSpace(add.Region))
	}

	if add.PostalCode != "" {
		if sb.Len() > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(strings.TrimSpace(add.PostalCode))
	}

	return sb.String()
}

//// HasAddressComplete checks if the essential fields of the address are filled.
//func (add *MailAddress) HasAddressComplete() bool {
//	return strings.TrimSpace(add.Addr1) != "" &&
//		strings.TrimSpace(add.City) != "" &&
//		strings.TrimSpace(add.State) != "" &&
//		strings.TrimSpace(add.PostalCode) != ""
//}

//// GetCityStateZipToSingleLine formats the city, state, and ZIP code into a single line.
//func (add *MailAddress) GetCityStateZipToSingleLine() string {
//	var sb strings.Builder
//	appendComma := func() {
//		if sb.Len() > 0 {
//			sb.WriteString(", ")
//		}
//	}
//
//	if city := strings.TrimSpace(add.City); city != "" {
//		appendComma()
//		sb.WriteString(city)
//	}
//
//	if state := strings.TrimSpace(add.State); state != "" {
//		appendComma()
//		sb.WriteString(state)
//	}
//
//	if postalCode := strings.TrimSpace(add.PostalCode); postalCode != "" {
//		if sb.Len() > 0 {
//			sb.WriteString(" ")
//		}
//		sb.WriteString(postalCode)
//	}
//
//	return sb.String()
//}

//// GenerateTextWithOptions formats the address into a text string with options.
//func (add *MailAddress) GenerateTextWithOptions(nameTo string) string {
//	if add == nil {
//		return ""
//	}
//
//	var sb strings.Builder
//
//	// Helper function to append a comma if the builder is not empty.
//	appendComma := func() {
//		if sb.Len() > 0 {
//			sb.WriteString(", ")
//		}
//	}
//
//	// Append the name if provided.
//	if name := strings.TrimSpace(nameTo); name != "" {
//		sb.WriteString(name)
//	}
//
//	// Append the address lines if provided.
//	for _, addrLine := range []string{add.Addr1, add.Addr2, add.Addr3} {
//		if line := strings.TrimSpace(addrLine); line != "" {
//			appendComma()
//			sb.WriteString(line)
//		}
//	}
//
//	// Append the city, state, and postal code if provided.
//	if city := strings.TrimSpace(add.City); city != "" {
//		appendComma()
//		sb.WriteString(city)
//	}
//	if state := strings.TrimSpace(add.State); state != "" {
//		appendComma()
//		sb.WriteString(state)
//	}
//	if postalCode := strings.TrimSpace(add.PostalCode); postalCode != "" {
//		if sb.Len() > 0 {
//			sb.WriteString(" ")
//		}
//		sb.WriteString(postalCode)
//	}
//
//	// Append the country if provided.
//	if country := strings.TrimSpace(add.Country); country != "" {
//		appendComma()
//		sb.WriteString(country)
//	}
//
//	return sb.String()
//}
//
//// GenerateHTMLWithOptions formats the address into an HTML string with options.
//func (add *MailAddress) GenerateHTMLWithOptions(class string, nameTo string, abbreviateState bool, ilocation ILocation) string {
//	if add == nil {
//		return ""
//	}
//
//	var sb strings.Builder
//
//	// Append the provided class along with a default class for styling.
//	class += " zsuite-address"
//	classInside := "zsuite-address"
//
//	// Start the address container div.
//	sb.WriteString(fmt.Sprintf(`<div class="%s">`, class))
//
//	// Check if the address has all the required fields completed.
//	if add.HasAddressComplete() {
//		// Append the name div if provided.
//		if name := strings.TrimSpace(nameTo); name != "" {
//			sb.WriteString(fmt.Sprintf(`<div class="%s-name">%s</div>`, classInside, name))
//		}
//
//		// Append the address line divs.
//		for _, line := range []string{add.Addr1, add.Addr2, add.Addr3} {
//			if addrLine := strings.TrimSpace(line); addrLine != "" {
//				sb.WriteString(fmt.Sprintf(`<div class="%s-addr">%s</div>`, classInside, addrLine))
//			}
//		}
//
//		// Start the group div for city, state, and zip.
//		sb.WriteString(fmt.Sprintf(`<div class="%s-group-poste">`, classInside))
//
//		// Abbreviate the state if requested and a location transformer is provided.
//		state := add.State
//		if abbreviateState && ilocation != nil {
//			abbreviatedState := ilocation.FindDivisionName(add.Country, add.State, true)
//			if abbreviatedState != "" {
//				state = abbreviatedState
//			}
//		}
//
//		// Append the city, state, and zip spans.
//		sb.WriteString(fmt.Sprintf(`<span class="%s-city">%s</span>`, classInside, add.City))
//		sb.WriteString(fmt.Sprintf(`<span class="%s-state">%s</span>`, classInside, state))
//		sb.WriteString(fmt.Sprintf(`<span class="%s-zip">%s</span>`, classInside, add.PostalCode))
//
//		// Close the group div.
//		sb.WriteString(`</div>`)
//
//		// Append the country div if provided.
//		if country := strings.TrimSpace(add.Country); country != "" {
//			sb.WriteString(fmt.Sprintf(`<div class="%s-country">%s</div>`, classInside, country))
//		}
//	}
//
//	// Close the address container div.
//	sb.WriteString(`</div>`)
//
//	return sb.String()
//}

//package acontact
//
//import (
//	"fmt"
//	"strings"
//	//"github.com/bojanz/address"
//)
//
//// MailAddress represents a physical mailing address using the bojanz/address library.
//type MailAddress struct {
//	Addr1      string `json:"addr1,omitempty"`      // First line of the address
//	Addr2      string `json:"addr2,omitempty"`      // Second line of the address (optional)
//	Addr3      string `json:"addr3,omitempty"`      // Third line of the address (optional)
//	City       string `json:"city,omitempty"`       // City or locality
//	State      string `json:"state,omitempty"`      // State, province, or region
//	PostalCode string `json:"postalCode,omitempty"` // Postal or ZIP code
//	Country    string `json:"country,omitempty"`    // Country name or code
//
//	// address.Address // Address struct from bojanz/address library
//}
//
////// NewMailAddress creates a new MailAddress with standardized country and state names.
////func (add *MailAddress) NewMailAddress() (*MailAddress, error) {
////	// Use the bojanz/address library to create a new address
////	newAddress := address.Address{
////		AddressLine1:       add.Addr1,
////		AddressLine2:       add.Addr2,
////		Locality:           add.City,
////		AdministrativeArea: add.State,
////		PostalCode:         add.PostalCode,
////		CountryCode:        add.Country,
////	}
////
////	// Validate and format the address using the bojanz/address library
////	formattedAddress, err := address.Format(newAddress)
////	if err != nil {
////		return nil, fmt.Errorf("failed to format address: %v", err)
////	}
////
////	// Return the formatted address
////	return &MailAddress{
////		Address:    formattedAddress,
////		Addr1:      formattedAddress.AddressLine1,
////		Addr2:      formattedAddress.AddressLine2,
////		Addr3:      add.Addr3, // Addr3 is not handled by bojanz/address and is kept as is
////		City:       formattedAddress.Locality,
////		State:      formattedAddress.AdministrativeArea,
////		PostalCode: formattedAddress.PostalCode,
////		Country:    formattedAddress.CountryCode,
////	}, nil
////}
//
//// NewMailAddress creates a new MailAddress with standardized country and state names.
//func (add *MailAddress) NewMailAddress(ilocation ILocation, abbreviateState bool, abbreviateCountry bool) (*MailAddress, error) {
//	if ilocation == nil {
//		return nil, fmt.Errorf("location transformer is nil")
//	}
//
//	// Standardize the country name or abbreviation
//	country := strings.TrimSpace(add.Country)
//	if country != "" {
//		country = ilocation.FindCountryName(country)
//		if abbreviateCountry {
//			country = ilocation.FindCountryShort(country)
//		}
//	}
//
//	// Standardize the state name or abbreviation
//	state := strings.TrimSpace(add.State)
//	if state != "" && country != "" {
//		state = ilocation.FindDivisionName(country, state, abbreviateState)
//	}
//
//	// Return the new MailAddress with standardized fields
//	return &MailAddress{
//		Addr1:      strings.TrimSpace(add.Addr1),
//		Addr2:      strings.TrimSpace(add.Addr2),
//		Addr3:      strings.TrimSpace(add.Addr3),
//		City:       strings.TrimSpace(add.City),
//		State:      state,
//		Country:    country,
//		PostalCode: strings.TrimSpace(add.PostalCode),
//	}, nil
//}
//
//// HasAddressComplete checks if the essential fields of the address are filled.
//func (add *MailAddress) HasAddressComplete() bool {
//	return strings.TrimSpace(add.Addr1) != "" &&
//		strings.TrimSpace(add.City) != "" &&
//		strings.TrimSpace(add.State) != "" &&
//		strings.TrimSpace(add.PostalCode) != ""
//}
//
//// IsEmpty checks if all the address fields are empty.
//func (add *MailAddress) IsEmpty() bool {
//	return strings.TrimSpace(add.Addr1) == "" &&
//		strings.TrimSpace(add.Addr2) == "" &&
//		strings.TrimSpace(add.Addr3) == "" &&
//		strings.TrimSpace(add.City) == "" &&
//		strings.TrimSpace(add.State) == "" &&
//		strings.TrimSpace(add.PostalCode) == "" &&
//		strings.TrimSpace(add.Country) == ""
//}
//
//// GetCityStateZipToSingleLine formats the city, state, and ZIP code into a single line.
//func (add *MailAddress) GetCityStateZipToSingleLine() string {
//	var sb strings.Builder
//	appendComma := func() {
//		if sb.Len() > 0 {
//			sb.WriteString(", ")
//		}
//	}
//
//	if city := strings.TrimSpace(add.City); city != "" {
//		appendComma()
//		sb.WriteString(city)
//	}
//
//	if state := strings.TrimSpace(add.State); state != "" {
//		appendComma()
//		sb.WriteString(state)
//	}
//
//	if postalCode := strings.TrimSpace(add.PostalCode); postalCode != "" {
//		if sb.Len() > 0 {
//			sb.WriteString(" ")
//		}
//		sb.WriteString(postalCode)
//	}
//
//	return sb.String()
//}
//
//// GenerateTextWithOptions formats the address into a text string with options.
//func (add *MailAddress) GenerateTextWithOptions(nameTo string) string {
//	if add == nil {
//		return ""
//	}
//
//	var sb strings.Builder
//
//	// Helper function to append a comma if the builder is not empty.
//	appendComma := func() {
//		if sb.Len() > 0 {
//			sb.WriteString(", ")
//		}
//	}
//
//	// Append the name if provided.
//	if name := strings.TrimSpace(nameTo); name != "" {
//		sb.WriteString(name)
//	}
//
//	// Append the address lines if provided.
//	for _, addrLine := range []string{add.Addr1, add.Addr2, add.Addr3} {
//		if line := strings.TrimSpace(addrLine); line != "" {
//			appendComma()
//			sb.WriteString(line)
//		}
//	}
//
//	// Append the city, state, and postal code if provided.
//	if city := strings.TrimSpace(add.City); city != "" {
//		appendComma()
//		sb.WriteString(city)
//	}
//	if state := strings.TrimSpace(add.State); state != "" {
//		appendComma()
//		sb.WriteString(state)
//	}
//	if postalCode := strings.TrimSpace(add.PostalCode); postalCode != "" {
//		if sb.Len() > 0 {
//			sb.WriteString(" ")
//		}
//		sb.WriteString(postalCode)
//	}
//
//	// Append the country if provided.
//	if country := strings.TrimSpace(add.Country); country != "" {
//		appendComma()
//		sb.WriteString(country)
//	}
//
//	return sb.String()
//}
//
//// GenerateHTMLWithOptions formats the address into an HTML string with options.
//func (add *MailAddress) GenerateHTMLWithOptions(class string, nameTo string, abbreviateState bool, ilocation ILocation) string {
//	if add == nil {
//		return ""
//	}
//
//	var sb strings.Builder
//
//	// Append the provided class along with a default class for styling.
//	class += " zsuite-address"
//	classInside := "zsuite-address"
//
//	// Start the address container div.
//	sb.WriteString(fmt.Sprintf(`<div class="%s">`, class))
//
//	// Check if the address has all the required fields completed.
//	if add.HasAddressComplete() {
//		// Append the name div if provided.
//		if name := strings.TrimSpace(nameTo); name != "" {
//			sb.WriteString(fmt.Sprintf(`<div class="%s-name">%s</div>`, classInside, name))
//		}
//
//		// Append the address line divs.
//		for _, line := range []string{add.Addr1, add.Addr2, add.Addr3} {
//			if addrLine := strings.TrimSpace(line); addrLine != "" {
//				sb.WriteString(fmt.Sprintf(`<div class="%s-addr">%s</div>`, classInside, addrLine))
//			}
//		}
//
//		// Start the group div for city, state, and zip.
//		sb.WriteString(fmt.Sprintf(`<div class="%s-group-poste">`, classInside))
//
//		// Abbreviate the state if requested and a location transformer is provided.
//		state := add.State
//		if abbreviateState && ilocation != nil {
//			abbreviatedState := ilocation.FindDivisionName(add.Country, add.State, true)
//			if abbreviatedState != "" {
//				state = abbreviatedState
//			}
//		}
//
//		// Append the city, state, and zip spans.
//		sb.WriteString(fmt.Sprintf(`<span class="%s-city">%s</span>`, classInside, add.City))
//		sb.WriteString(fmt.Sprintf(`<span class="%s-state">%s</span>`, classInside, state))
//		sb.WriteString(fmt.Sprintf(`<span class="%s-zip">%s</span>`, classInside, add.PostalCode))
//
//		// Close the group div.
//		sb.WriteString(`</div>`)
//
//		// Append the country div if provided.
//		if country := strings.TrimSpace(add.Country); country != "" {
//			sb.WriteString(fmt.Sprintf(`<div class="%s-country">%s</div>`, classInside, country))
//		}
//	}
//
//	// Close the address container div.
//	sb.WriteString(`</div>`)
//
//	return sb.String()
//}
