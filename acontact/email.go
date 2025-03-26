package acontact

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/aemail"
	"sort"
	"strings"
)

// Email represents an email address with its type, address, default status, and associated name.
type Email struct {
	Type      EmailType           `json:"type,omitempty"`      // Type of the email (e.g., mobile, home)
	Address   aemail.EmailAddress `json:"address,omitempty"`   // The actual email address
	IsDefault bool                `json:"isDefault,omitempty"` // Indicates if this is the default email address
	Name      string              `json:"name,omitempty"`      // The name associated with the email address
}

func (e *Email) Validate() error {
	if e.Type.IsEmpty() {
		return fmt.Errorf("email type is empty")
	}
	e.Name = strings.TrimSpace(e.Name)
	if e.Address.IsEmpty() {
		return fmt.Errorf("email address is empty")
	}
	return e.Address.Validate()
}

// Emails is a slice of Email pointers, representing a collection of email addresses.
type Emails []*Email

// FindByType searches for an email by its type.
func (es Emails) FindByType(emailType EmailType) *Email {
	return es.findByType(emailType, false)
}

// FindByTypeOrDefault searches for an email by its type or returns the default email.
func (es Emails) FindByTypeOrDefault(emailType EmailType) *Email {
	return es.findByType(emailType, true)
}

// findByType is a helper function that searches for an email by type and optionally returns the default email.
func (es Emails) findByType(emailType EmailType, checkDefault bool) *Email {
	emailType = EmailType(emailType.GetType()) // Normalize the email type
	var def *Email                             // Placeholder for the default email
	for _, e := range es {
		if e.Type.ToStringTrimLower() == emailType.ToStringTrimLower() {
			return e // Return the matching email
		}
		if e.IsDefault {
			def = e // Remember the default email
		}
	}
	if checkDefault {
		return def // Return the default email if allowed and no match was found
	}
	return nil // No match found
}

// FindByAddress searches for an email by its address.
func (es Emails) FindByAddress(address aemail.EmailAddress) *Email {
	for _, e := range es {
		if e.Address.ToStringTrimLower() == address.ToStringTrimLower() {
			return e // Return the matching email
		}
	}
	return nil // No match found
}

// HasType checks if an email of the specified type exists in the collection.
func (es Emails) HasType(emailType EmailType) bool {
	return es.FindByType(emailType) != nil
}

// HasTypeWithDefault checks if an email of the specified type exists, or if there's a default email.
func (es Emails) HasTypeWithDefault(emailType EmailType, allowDefault bool) bool {
	return es.findByType(emailType, allowDefault) != nil
}

// HasTypeOrDefault checks if an email of the specified type exists, or if there's a default email.
func (es Emails) HasTypeOrDefault(emailType EmailType) bool {
	return es.FindByTypeOrDefault(emailType) != nil
}

// HasAddress checks if an email with the specified address exists in the collection.
func (es Emails) HasAddress(address aemail.EmailAddress) bool {
	return es.FindByAddress(address) != nil
}

// Clone creates a deep copy of the Emails collection.
func (es Emails) Clone() Emails {
	b, err := json.Marshal(es)
	if err != nil {
		return nil // Return nil if marshaling fails
	}
	var clone Emails
	if err := json.Unmarshal(b, &clone); err != nil {
		return nil // Return nil if unmarshaling fails
	}
	return clone // Return the deep copy
}

// MergeFrom adds emails from another collection that are not already present.
func (es *Emails) MergeFrom(target Emails) {
	if es == nil || target == nil {
		return // Do nothing if either collection is nil
	}
	for _, t := range target {
		if t.Type.IsEmpty() {
			continue // Skip empty types
		}
		isFound := false
		for _, e := range *es {
			if e.Type.ToStringTrimLower() == t.Type.ToStringTrimLower() {
				isFound = true
				break // Email type already exists
			}
		}
		if !isFound {
			*es = append(*es, t) // Add the email if it's not found
		}
	}
}

// Set adds or updates an email in the collection.
func (es *Emails) Set(email *Email) {
	if email == nil || email.Type.IsEmpty() || email.Address.IsEmpty() {
		return // Do nothing if the email is nil or has empty fields
	}
	newEmails := Emails{}
	for _, e := range *es {
		if e.Type.ToStringTrimLower() == email.Type.ToStringTrimLower() {
			continue // Skip to replace the email
		} else if e.IsDefault && email.IsDefault {
			e.IsDefault = false // Unset the default if the new email is the default
		}
		newEmails = append(newEmails, e)
	}
	newEmails = append(newEmails, email) // Add the new email

	// Sort the emails, placing the default email at the top
	sort.SliceStable(newEmails, func(ii, jj int) bool {
		return newEmails[ii].IsDefault || newEmails[ii].Type < newEmails[jj].Type
	})

	*es = newEmails // Update the original collection
}

// Remove deletes an email of the specified type from the collection.
func (es *Emails) Remove(emailType EmailType) {
	if emailType.IsEmpty() {
		return // Do nothing if the email type is empty
	}
	newArr := Emails{}
	for _, e := range *es {
		if e.Type.ToStringTrimLower() == emailType.ToStringTrimLower() {
			continue // Skip the email to be removed
		}
		newArr = append(newArr, e)
	}
	*es = newArr // Update the original collection with the remaining emails
}
