package aclient_ldap

import (
	"testing"

	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/stretchr/testify/assert"
)

func TestMapLDAP_GetSetDelete(t *testing.T) {
	ml := NewMapLDAP()

	// Define an AdapterName for testing
	name := aconns.AdapterName("test_adapter")

	// Test setting an AClientLDAP
	client := &AClientLDAP{}
	ml.Set(name, client)

	// Test getting the AClientLDAP
	retrievedClient, exists := ml.Get(name)
	assert.True(t, exists)
	assert.Equal(t, client, retrievedClient)

	// Test getting a non-existent AClientLDAP
	_, exists = ml.Get(aconns.AdapterName("non_existent_adapter"))
	assert.False(t, exists)

	// Test deleting the AClientLDAP
	ml.Delete(name)
	_, exists = ml.Get(name)
	assert.False(t, exists)
}
