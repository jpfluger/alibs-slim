package aclient_redis

import (
	"testing"

	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/stretchr/testify/assert"
)

func TestMapRedis_GetSetDelete(t *testing.T) {
	mr := NewMapRedis()

	// Define an AdapterName for testing
	name := aconns.AdapterName("test_adapter")

	// Test setting an AClientRedis
	client := &AClientRedis{}
	mr.Set(name, client)

	// Test getting the AClientRedis
	retrievedClient, exists := mr.Get(name)
	assert.True(t, exists)
	assert.Equal(t, client, retrievedClient)

	// Test getting a non-existent AClientRedis
	_, exists = mr.Get(aconns.AdapterName("non_existent_adapter"))
	assert.False(t, exists)

	// Test deleting the AClientRedis
	mr.Delete(name)
	_, exists = mr.Get(name)
	assert.False(t, exists)
}
