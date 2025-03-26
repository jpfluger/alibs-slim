package acontact

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/auser"
)

type IContactUID interface {
	IContact
	GetUID() auser.UID
	IsDefaultContact() bool
	GetEntityTypes() EntityTypes
	GetPersonTypes() PersonTypes
}

// ContactUID represents a Contact with a unique UID and additional metadata.
type ContactUID struct {
	UID         auser.UID   `json:"uid,omitempty"`
	Contact     Contact     `json:"contact"`
	IsDefault   bool        `json:"isDefault,omitempty"`
	EntityTypes EntityTypes `json:"entityTypes,omitempty"` // Relationship-defined types for Entity
	PersonTypes PersonTypes `json:"personTypes,omitempty"` // Relationship-defined types for Person
}

// Validate ensures that the ContactUID is valid.
func (c *ContactUID) Validate() error {
	if c == nil {
		return fmt.Errorf("contact is nil")
	}
	if c.UID.IsNil() {
		return fmt.Errorf("uid is nil")
	}
	if err := c.Contact.Validate(); err != nil {
		return err
	}
	return nil
}

// GetUID returns the UID of the contact.
func (c *ContactUID) GetUID() auser.UID {
	return c.UID
}

// IsDefaultContact returns whether the contact is marked as default.
func (c *ContactUID) IsDefaultContact() bool {
	return c.IsDefault
}

// GetEntityTypes returns the entity types associated with the contact.
func (c *ContactUID) GetEntityTypes() EntityTypes {
	return c.EntityTypes
}

// GetPersonTypes returns the person types associated with the contact.
func (c *ContactUID) GetPersonTypes() PersonTypes {
	return c.PersonTypes
}

// SetDefault sets the contact's default status.
func (c *ContactUID) SetDefault(isDefault bool) {
	c.IsDefault = isDefault
}

// SetEntityTypes sets the entity types associated with the contact.
func (c *ContactUID) SetEntityTypes(types EntityTypes) {
	c.EntityTypes = types
}

// SetPersonTypes sets the person types associated with the contact.
func (c *ContactUID) SetPersonTypes(types PersonTypes) {
	c.PersonTypes = types
}

// ContactUIDs represents a collection of ContactUID pointers.
type ContactUIDs []*ContactUID

// SetContactUID adds or updates a ContactUID in the ContactUIDs array by UID.
// If the UID already exists, it updates the existing ContactUID; otherwise, it adds a new one.
func (cs *ContactUIDs) SetContactUID(contactUID *ContactUID) error {
	if contactUID == nil {
		return fmt.Errorf("contactUID is nil")
	}
	if contactUID.UID.IsNil() {
		return fmt.Errorf("contactUID UID is nil")
	}

	for i, existing := range *cs {
		if existing.UID == contactUID.UID {
			(*cs)[i] = contactUID // Update existing contactUID
			return nil
		}
	}

	// Add new contactUID
	*cs = append(*cs, contactUID)
	return nil
}

// RemoveContactUID removes a ContactUID from the ContactUIDs array by UID.
func (cs *ContactUIDs) RemoveContactUID(uid auser.UID) error {
	if uid.IsNil() {
		return fmt.Errorf("uid is nil")
	}

	for i, contactUID := range *cs {
		if contactUID.UID == uid {
			// Remove contactUID by slicing
			*cs = append((*cs)[:i], (*cs)[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("contactUID with UID %v not found", uid)
}

// RemoveContactByCID removes a ContactUID from the ContactUIDs array by CID.
func (cs *ContactUIDs) RemoveContactByCID(cid CID) error {
	if cid.IsNil() {
		return fmt.Errorf("cid is nil")
	}

	for i, contactUID := range *cs {
		if contactUID.Contact.CID == cid {
			// Remove contactUID by slicing
			*cs = append((*cs)[:i], (*cs)[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("contactUID with CID %v not found", cid)
}

// FindContactUID finds a ContactUID in the ContactUIDs array by UID.
// Returns the ContactUID and a boolean indicating whether it was found.
func (cs *ContactUIDs) FindContactUID(uid auser.UID) (*ContactUID, bool) {
	if uid.IsNil() {
		return nil, false
	}

	for _, contactUID := range *cs {
		if contactUID.UID == uid {
			return contactUID, true
		}
	}

	return nil, false
}

// FindContactByCID finds a ContactUID in the ContactUIDs array by CID.
// Returns the ContactUID and a boolean indicating whether it was found.
func (cs *ContactUIDs) FindContactByCID(cid CID) (*ContactUID, bool) {
	if cid.IsNil() {
		return nil, false
	}

	for _, contactUID := range *cs {
		if contactUID.Contact.CID == cid {
			return contactUID, true
		}
	}

	return nil, false
}

// FindContactsByUIDSubset finds all ContactUIDs in the ContactUIDs array that match the provided UIDs.
// Returns a subset of ContactUIDs.
func (cs *ContactUIDs) FindContactsByUIDSubset(uids []auser.UID) ContactUIDs {
	subset := ContactUIDs{}

	uidSet := make(map[auser.UID]bool)
	for _, uid := range uids {
		uidSet[uid] = true
	}

	for _, contactUID := range *cs {
		if uidSet[contactUID.UID] {
			subset = append(subset, contactUID)
		}
	}

	return subset
}

// ToMap converts a ContactUIDs slice to a ContactUIDMap keyed by UID.
func (cs *ContactUIDs) ToMap() ContactUIDMap {
	contactUIDMap := make(ContactUIDMap)
	for _, contactUID := range *cs {
		if !contactUID.UID.IsNil() {
			if _, exists := contactUIDMap[contactUID.UID]; !exists {
				contactUIDMap[contactUID.UID] = make(ContactsMap)
			}
			contactUIDMap[contactUID.UID][contactUID.Contact.CID] = &contactUID.Contact
		}
	}
	return contactUIDMap
}

type ContactUIDMap map[auser.UID]ContactsMap

// FindByUID finds a ContactsMap for a specific UID in the ContactUIDMap.
func (cum ContactUIDMap) FindByUID(uid auser.UID) (ContactsMap, bool) {
	if uid.IsNil() {
		return nil, false
	}
	if contacts, exists := cum[uid]; exists {
		return contacts, true
	}
	return nil, false
}

// FindContactByUIDAndCID finds a specific Contact in the ContactUIDMap by UID and CID.
func (cum ContactUIDMap) FindContactByUIDAndCID(uid auser.UID, cid CID) (*Contact, bool) {
	if uid.IsNil() {
		return nil, false
	}
	if cid.IsNil() {
		return nil, false
	}
	if contactsMap, exists := cum.FindByUID(uid); exists {
		return contactsMap.FindByCID(cid)
	}
	return nil, false
}
