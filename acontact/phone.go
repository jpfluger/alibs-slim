package acontact

import (
	"encoding/json"
	"fmt"
	"github.com/nyaruka/phonenumbers"
	"sort"
	"strings"
)

// https://github.com/nyaruka/phonenumbers

// Phone represents a phone number with its type, country, and default status.
type Phone struct {
	Type      PhoneType `json:"type,omitempty"`      // The type of phone number (e.g., mobile, home)
	Number    string    `json:"number,omitempty"`    // The actual phone number
	Country   string    `json:"country,omitempty"`   // The country code of the phone number
	IsDefault bool      `json:"isDefault,omitempty"` // Indicates if this is the default phone number
}

// Validate checks if the Phone fields are valid and tries to parse the number.
func (p *Phone) Validate() error {
	if p.Type.IsEmpty() {
		return fmt.Errorf("phone type is empty")
	}
	p.Number = strings.TrimSpace(p.Number)
	p.Country = strings.TrimSpace(p.Country)

	if p.Number == "" {
		return fmt.Errorf("phone number is empty")
	}

	// Attempt to parse the phone number and determine the country code.
	// Parsed number gets standardized in E.164 format.
	// First check if the number contains a '+' sign
	if strings.Contains(p.Number, "+") {
		num, region, err := p.guessNumberCountry()
		if err != nil {
			return err
		}
		p.Number = num
		p.Country = region
	} else if p.Country != "" {
		// If the number does not contain a '+' sign but a country is provided, use the country as the region code
		num, err := phonenumbers.Parse(p.Number, p.Country)
		if err != nil {
			// If parsing fails with the provided country, the number is invalid
			return fmt.Errorf("invalid phone number: %v", err)
		}
		// No need to update the country as it was already provided
		p.Number = num.String()
	} else {
		num, region, err := p.guessNumberCountry()
		if err != nil {
			return err
		}
		p.Number = num
		p.Country = region
	}

	return nil
}

func (p *Phone) guessNumberCountry() (phone, country string, err error) {
	num, err := phonenumbers.Parse(p.Number, "")
	if err != nil {
		// If parsing fails, the number is invalid
		return "", "", fmt.Errorf("invalid phone number: %v", err)
	}
	// If parsing is successful, update the Country field with the determined country code
	regionCode := phonenumbers.GetRegionCodeForNumber(num)
	if regionCode == "" {
		// If the country code could not be determined, set it as unknown
		regionCode = "Unknown"
	}
	return num.String(), regionCode, nil
}

// Phones is a slice of Phone pointers, representing a collection of phone numbers.
type Phones []*Phone

// FindByType searches for a phone number by its type.
func (ps Phones) FindByType(phoneType PhoneType) *Phone {
	return ps.findByType(phoneType, false)
}

// FindByTypeOrDefault searches for a phone number by its type or returns the default number.
func (ps Phones) FindByTypeOrDefault(phoneType PhoneType) *Phone {
	return ps.findByType(phoneType, true)
}

// findByType is a helper function that searches for a phone number by type and optionally returns the default number.
func (ps Phones) findByType(phoneType PhoneType, checkDefault bool) *Phone {
	var defaultPhone *Phone
	for _, p := range ps {
		if p.Type.ToStringTrimLower() == phoneType.ToStringTrimLower() {
			return p
		}
		if p.IsDefault {
			defaultPhone = p
		}
	}
	if checkDefault {
		return defaultPhone
	}
	return nil
}

// FindByNumber searches for a phone number by its number.
func (ps Phones) FindByNumber(number string) *Phone {
	for _, p := range ps {
		if p.Number == number {
			return p
		}
	}
	return nil
}

func (ps Phones) HasTypeWithDefault(phoneType PhoneType, allowDefault bool) bool {
	return ps.findByType(phoneType, allowDefault) != nil
}

// HasType checks if a phone number of the specified type exists in the collection.
func (ps Phones) HasType(phoneType PhoneType) bool {
	return ps.FindByType(phoneType) != nil
}

// HasTypeOrDefault checks if a phone number of the specified type exists, or if there's a default number.
func (ps Phones) HasTypeOrDefault(phoneType PhoneType) bool {
	return ps.FindByTypeOrDefault(phoneType) != nil
}

// HasNumber checks if a phone number with the specified number exists in the collection.
func (ps Phones) HasNumber(number string) bool {
	return ps.FindByNumber(number) != nil
}

// Clone creates a deep copy of the Phones collection.
func (ps Phones) Clone() Phones {
	b, err := json.Marshal(ps)
	if err != nil {
		return nil
	}
	var clone Phones
	if err := json.Unmarshal(b, &clone); err != nil {
		return nil
	}
	return clone
}

// MergeFrom adds phone numbers from another collection that are not already present.
func (ps *Phones) MergeFrom(target Phones) {
	if ps == nil || target == nil {
		return
	}
	for _, t := range target {
		if t.Type.IsEmpty() {
			continue
		}
		isFound := false
		for _, p := range *ps {
			if p.Type.ToStringTrimLower() == t.Type.ToStringTrimLower() {
				isFound = true
				break
			}
		}
		if !isFound {
			*ps = append(*ps, t)
		}
	}
}

// Set adds or updates a phone number in the collection.
func (ps *Phones) Set(phone *Phone) {
	if phone == nil || phone.Type.IsEmpty() || strings.TrimSpace(phone.Number) == "" {
		return
	}
	// Create a new slice for the updated phone numbers
	newPhones := Phones{}
	for _, p := range *ps {
		if p.Type.ToStringTrimLower() == phone.Type.ToStringTrimLower() {
			continue // Skip the phone number of the same type to replace it
		} else if p.IsDefault && phone.IsDefault {
			p.IsDefault = false // Unset the default if the new phone number is the default
		}
		newPhones = append(newPhones, p)
	}
	newPhones = append(newPhones, phone) // Add the new phone number

	// Sort the phone numbers, placing the default number at the top
	sort.SliceStable(newPhones, func(i, j int) bool {
		return newPhones[i].IsDefault || newPhones[i].Type < newPhones[j].Type
	})

	*ps = newPhones // Update the original collection
}

// Remove deletes a phone number of the specified type from the collection.
func (ps *Phones) Remove(phoneType PhoneType) {
	if phoneType.IsEmpty() {
		return
	}
	newPhones := Phones{}
	for _, p := range *ps {
		if p.Type.ToStringTrimLower() == phoneType.ToStringTrimLower() {
			continue // Skip the phone number of the type to be removed
		}
		newPhones = append(newPhones, p)
	}
	*ps = newPhones // Update the original collection with the remaining phone numbers
}
