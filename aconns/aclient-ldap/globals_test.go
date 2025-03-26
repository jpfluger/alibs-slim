package aclient_ldap

import (
	"testing"

	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/stretchr/testify/assert"
)

func TestConnMapGlobal_GetSetRemove(t *testing.T) {
	// Initialize the global connection map
	connMap = &connMapGlobal{Map: NewMapLDAP()}

	// Define an AdapterName for testing
	name := aconns.AdapterName("test_adapter")

	// Mock the AClientLDAP instance
	client := &AClientLDAP{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Name: name,
			},
		},
	}

	// Test setting an AClientLDAP
	err := connMap.Set(client)
	assert.NoError(t, err)

	// Test getting the AClientLDAP
	retrievedClient := connMap.Get(name)
	assert.NotNil(t, retrievedClient)
	assert.Equal(t, client, retrievedClient)

	// Test getting a non-existent AClientLDAP
	nonExistentClient := connMap.Get(aconns.AdapterName("non_existent_adapter"))
	assert.Nil(t, nonExistentClient)

	// Test removing the AClientLDAP
	connMap.Remove(name)
	removedClient := connMap.Get(name)
	assert.Nil(t, removedClient)
}

func TestConnMapGlobal_SetErrors(t *testing.T) {
	// Initialize the global connection map
	connMap = &connMapGlobal{Map: NewMapLDAP()}

	// Test setting a nil AClientLDAP
	err := connMap.Set(nil)
	assert.Error(t, err)
	assert.Equal(t, "connMapGlobal is nil", err.Error())

	// Test setting an AClientLDAP with an empty name
	client := &AClientLDAP{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Name: "",
			},
		},
	}
	err = connMap.Set(client)
	assert.Error(t, err)
	assert.Equal(t, "name is empty", err.Error())
}
