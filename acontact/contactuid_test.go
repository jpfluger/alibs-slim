package acontact

import (
	"github.com/stretchr/testify/assert"
	"github.com/jpfluger/alibs-slim/auser"
	"testing"
)

func TestContactUIDs_SetContactUID(t *testing.T) {
	contacts := ContactUIDs{}
	uid := auser.NewUID()
	contactUID := &ContactUID{
		UID: uid,
		Contact: Contact{
			CID: NewCID(),
		},
	}

	// Add contactUID
	err := contacts.SetContactUID(contactUID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(contacts))

	// Update contactUID
	contactUIDUpdated := &ContactUID{
		UID: uid,
		Contact: Contact{
			CID: contactUID.Contact.CID,
			ContactCore: ContactCore{
				Name: Name{First: "John", Last: "Doe"},
			},
		},
	}
	err = contacts.SetContactUID(contactUIDUpdated)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(contacts))
	assert.Equal(t, "John", contacts[0].Contact.Name.First)
}

func TestContactUIDs_RemoveContactUID(t *testing.T) {
	contacts := ContactUIDs{}
	uid := auser.NewUID()
	contactUID := &ContactUID{UID: uid}

	// Add and remove contactUID
	contacts.SetContactUID(contactUID)
	assert.Equal(t, 1, len(contacts))
	err := contacts.RemoveContactUID(uid)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(contacts))

	// Attempt to remove non-existent contactUID
	err = contacts.RemoveContactUID(uid)
	assert.Error(t, err)
}

func TestContactUIDs_FindContactUID(t *testing.T) {
	contacts := ContactUIDs{}
	uid := auser.NewUID()
	contactUID := &ContactUID{UID: uid}

	// Add and find contactUID
	contacts.SetContactUID(contactUID)
	foundContact, found := contacts.FindContactUID(uid)
	assert.True(t, found)
	assert.Equal(t, contactUID, foundContact)

	// Attempt to find non-existent contactUID
	_, found = contacts.FindContactUID(auser.NewUID())
	assert.False(t, found)
}

func TestContactUIDs_ToMap(t *testing.T) {
	uid1 := auser.NewUID()
	contactUIDs := ContactUIDs{
		&ContactUID{
			UID: uid1,
			Contact: Contact{
				CID: NewCID(),
			},
		},
		&ContactUID{
			UID: uid1,
			Contact: Contact{
				CID: NewCID(),
			},
		},
	}
	contactUIDMap := contactUIDs.ToMap()
	assert.Equal(t, 1, len(contactUIDMap))
	assert.Equal(t, 2, len(contactUIDMap[uid1]))
}

func TestContactUIDMap_FindByUID(t *testing.T) {
	uid1 := auser.NewUID()
	cid1 := NewCID()
	contactUIDMap := ContactUIDMap{
		uid1: ContactsMap{
			cid1: &Contact{CID: cid1},
		},
	}
	contacts, found := contactUIDMap.FindByUID(uid1)
	assert.True(t, found)
	assert.NotNil(t, contacts)

	_, found = contactUIDMap.FindByUID(auser.NewUID())
	assert.False(t, found)
}

func TestContactUIDMap_FindContactByUIDAndCID(t *testing.T) {
	uid1 := auser.NewUID()
	cid1 := NewCID()
	contactUIDMap := ContactUIDMap{
		uid1: ContactsMap{
			cid1: &Contact{CID: cid1},
		},
	}
	contact, found := contactUIDMap.FindContactByUIDAndCID(uid1, cid1)
	assert.True(t, found)
	assert.NotNil(t, contact)

	_, found = contactUIDMap.FindContactByUIDAndCID(uid1, NewCID())
	assert.False(t, found)
}

func TestContactUID_GettersAndSetters(t *testing.T) {
	uid1 := auser.NewUID()
	cid1 := NewCID()
	contactUID := &ContactUID{
		UID: uid1,
		Contact: Contact{
			CID: cid1,
		},
		IsDefault:   true,
		EntityTypes: EntityTypes{"Business"},
		PersonTypes: PersonTypes{"Employee"},
	}

	// Test GetUID
	assert.Equal(t, uid1, contactUID.GetUID())

	// Test GetCID
	assert.Equal(t, cid1, contactUID.GetCID())

	// Test IsDefaultContact
	assert.True(t, contactUID.IsDefaultContact())

	// Test GetEntityTypes
	assert.Equal(t, EntityTypes{"Business"}, contactUID.GetEntityTypes())

	// Test GetPersonTypes
	assert.Equal(t, PersonTypes{"Employee"}, contactUID.GetPersonTypes())

	// Test SetDefault
	contactUID.SetDefault(false)
	assert.False(t, contactUID.IsDefaultContact())

	// Test SetEntityTypes
	newEntityTypes := EntityTypes{"Vendor"}
	contactUID.SetEntityTypes(newEntityTypes)
	assert.Equal(t, newEntityTypes, contactUID.GetEntityTypes())

	// Test SetPersonTypes
	newPersonTypes := PersonTypes{"Manager"}
	contactUID.SetPersonTypes(newPersonTypes)
	assert.Equal(t, newPersonTypes, contactUID.GetPersonTypes())
}
