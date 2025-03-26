package acontact

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContacts_SetContact(t *testing.T) {
	contacts := Contacts{}
	cid := NewCID()
	contact := &Contact{
		CID: cid,
		ContactCore: ContactCore{
			Name: Name{First: "John", Last: "Doe"},
		},
	}

	// Add contact
	err := contacts.SetContact(contact)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(contacts))

	// Update contact
	contactUpdated := &Contact{
		CID: cid,
		ContactCore: ContactCore{
			Name: Name{First: "Johnathan", Last: "Doe"},
		},
	}
	err = contacts.SetContact(contactUpdated)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(contacts))
	assert.Equal(t, "Johnathan", contacts[0].Name.First)
}

func TestContacts_RemoveContact(t *testing.T) {
	contacts := Contacts{}
	cid := NewCID()
	contact := &Contact{CID: cid}

	// Add and remove contact
	contacts.SetContact(contact)
	assert.Equal(t, 1, len(contacts))
	err := contacts.RemoveContact(cid)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(contacts))

	// Attempt to remove non-existent contact
	err = contacts.RemoveContact(cid)
	assert.Error(t, err)
}

func TestContacts_FindContact(t *testing.T) {
	contacts := Contacts{}
	cid := NewCID()
	contact := &Contact{CID: cid}

	// Add and find contact
	contacts.SetContact(contact)
	foundContact, found := contacts.FindContact(cid)
	assert.True(t, found)
	assert.Equal(t, contact, foundContact)

	// Attempt to find non-existent contact
	_, found = contacts.FindContact(NewCID())
	assert.False(t, found)
}

func TestContacts_ToMap(t *testing.T) {
	cid1 := NewCID()
	cid2 := NewCID()
	contacts := Contacts{
		&Contact{CID: cid1},
		&Contact{CID: cid2},
	}
	contactsMap := contacts.ToMap()
	assert.Equal(t, 2, len(contactsMap))
	assert.NotNil(t, contactsMap[cid1])
	assert.NotNil(t, contactsMap[cid2])
}

func TestContactsMap_FindByCID(t *testing.T) {
	cid1 := NewCID()
	cid2 := NewCID()
	contactsMap := ContactsMap{
		cid1: &Contact{CID: cid1},
		cid2: &Contact{CID: cid2},
	}
	contact, found := contactsMap.FindByCID(cid1)
	assert.True(t, found)
	assert.Equal(t, cid1, contact.CID)

	_, found = contactsMap.FindByCID(NewCID())
	assert.False(t, found)
}
