package aclient_redis

import (
	"testing"

	"github.com/alexedwards/scs/redisstore"
	"github.com/stretchr/testify/assert"
)

func TestRedisStoreMap_GetSet(t *testing.T) {
	rsm := NewRedisStoreMap()

	// Test setting a RedisStore
	store := &redisstore.RedisStore{}
	rsm.Set("test_prefix", store)

	// Test getting the RedisStore
	retrievedStore, exists := rsm.Get("test_prefix")
	assert.True(t, exists)
	assert.Equal(t, store, retrievedStore)

	// Test getting a non-existent RedisStore
	_, exists = rsm.Get("non_existent_prefix")
	assert.False(t, exists)
}
