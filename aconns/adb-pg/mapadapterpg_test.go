package adb_pg

import (
	"testing"

	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/stretchr/testify/assert"
)

func TestMapPG_GetSetDelete(t *testing.T) {
	mpg := NewMapAdapterPG()

	// Define an AdapterName for testing
	name := aconns.AdapterName("test_adapter")

	// Test setting an ADBPG
	client := &ADBPG{}
	mpg.Set(name, client)

	// Test getting the ADBPG
	retrievedClient, exists := mpg.Get(name)
	assert.True(t, exists)
	assert.Equal(t, client, retrievedClient)

	// Test getting a non-existent ADBPG
	_, exists = mpg.Get(aconns.AdapterName("non_existent_adapter"))
	assert.False(t, exists)

	// Test deleting the ADBPG
	mpg.Delete(name)
	_, exists = mpg.Get(name)
	assert.False(t, exists)
}
