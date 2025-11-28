package aclient_badger

import (
	"testing"

	"github.com/alexedwards/scs/badgerstore"
	"github.com/stretchr/testify/assert"
)

func TestBadgerStoreMap_GetSet(t *testing.T) {
	bsm := NewBadgerStoreMap()

	// Test setting a BadgerStore
	store := &badgerstore.BadgerStore{}
	bsm.Set("test_prefix", store)

	// Test getting the BadgerStore
	retrievedStore, exists := bsm.Get("test_prefix")
	assert.True(t, exists)
	assert.Equal(t, store, retrievedStore)

	// Test getting a non-existent BadgerStore
	_, exists = bsm.Get("non_existent_prefix")
	assert.False(t, exists)
}
