package acontact

import (
	"errors"
	"fmt"
)

type IContact interface {
	IContactCore
	GetCID() CID
}

// Contact represents an individual contact.
type Contact struct {
	CID CID `json:"cid,omitempty"`
	ContactCore
}

// Validate ensures that the Contact is valid.
func (c *Contact) Validate() error {
	if c == nil {
		return fmt.Errorf("contact is nil")
	}
	if c.CID.IsNil() {
		return fmt.Errorf("cid is nil")
	}
	if err := c.ContactCore.Validate(); err != nil {
		return err
	}
	return nil
}

// GetCID returns the CID of the contact.
func (c *ContactUID) GetCID() CID {
	return c.Contact.CID
}

// Contacts represents a collection of Contact pointers.
type Contacts []*Contact

// SetContact adds or updates a contact in the Contacts array by CID.
// If the CID already exists, it updates the existing Contact; otherwise, it adds a new one.
func (cs *Contacts) SetContact(contact *Contact) error {
	if contact == nil {
		return errors.New("contact is nil")
	}
	if contact.CID.IsNil() {
		return errors.New("contact CID is nil")
	}

	for i, existing := range *cs {
		if existing.CID == contact.CID {
			(*cs)[i] = contact // Update existing contact
			return nil
		}
	}

	// Add new contact
	*cs = append(*cs, contact)
	return nil
}

// RemoveContact removes a contact from the Contacts array by CID.
func (cs *Contacts) RemoveContact(cid CID) error {
	if cid.IsNil() {
		return errors.New("cid is nil")
	}

	for i, contact := range *cs {
		if contact.CID == cid {
			// Remove contact by slicing
			*cs = append((*cs)[:i], (*cs)[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("contact with CID %v not found", cid)
}

// FindContact finds a contact in the Contacts array by CID.
// Returns the contact and a boolean indicating whether it was found.
func (cs *Contacts) FindContact(cid CID) (*Contact, bool) {
	if cid.IsNil() {
		return nil, false
	}

	for _, contact := range *cs {
		if contact.CID == cid {
			return contact, true
		}
	}

	return nil, false
}

// ToMap converts a Contacts slice to a ContactsMap keyed by CID.
func (cs *Contacts) ToMap() ContactsMap {
	contactsMap := make(ContactsMap)
	for _, contact := range *cs {
		if !contact.CID.IsNil() {
			contactsMap[contact.CID] = contact
		}
	}
	return contactsMap
}

type ContactsMap map[CID]*Contact

// FindByCID finds a Contact in the ContactsMap by CID.
func (cm ContactsMap) FindByCID(cid CID) (*Contact, bool) {
	if cid.IsNil() {
		return nil, false
	}
	if contact, exists := cm[cid]; exists {
		return contact, true
	}
	return nil, false
}
