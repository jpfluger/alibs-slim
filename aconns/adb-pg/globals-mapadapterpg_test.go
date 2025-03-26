package adb_pg

import (
	"testing"

	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/stretchr/testify/assert"
)

func TestConnMapGlobal_GetSetRemove(t *testing.T) {
	// Initialize the global connection map
	gAdapterMap = &adapterMapGlobal{Map: NewMapAdapterPG()}

	// Define an AdapterName for testing
	name := aconns.AdapterName("test_adapter")

	// Mock the ADBPG instance
	client := &ADBPG{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Name: name,
			},
		},
	}

	// Test setting an ADBPG
	err := gAdapterMap.Set(client)
	assert.NoError(t, err)

	// Test getting the ADBPG
	retrievedClient := gAdapterMap.Get(name)
	assert.NotNil(t, retrievedClient)
	assert.Equal(t, client, retrievedClient)

	// Test getting a non-existent ADBPG
	nonExistentClient := gAdapterMap.Get(aconns.AdapterName("non_existent_adapter"))
	assert.Nil(t, nonExistentClient)

	// Test removing the ADBPG
	gAdapterMap.Remove(name)
	removedClient := gAdapterMap.Get(name)
	assert.Nil(t, removedClient)
}

func TestConnMapGlobal_SetErrors(t *testing.T) {
	// Initialize the global connection map
	gAdapterMap = &adapterMapGlobal{Map: NewMapAdapterPG()}

	// Test setting a nil ADBPG
	err := gAdapterMap.Set(nil)
	assert.Error(t, err)
	assert.Equal(t, "adapterMapGlobal is nil", err.Error())

	// Test setting an ADBPG with an empty name
	client := &ADBPG{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Name: "",
			},
		},
	}
	err = gAdapterMap.Set(client)
	assert.Error(t, err)
	assert.Equal(t, "name is empty", err.Error())
}
