package aclient_redis

import (
	"github.com/gomodule/redigo/redis"
	"testing"

	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/stretchr/testify/assert"
)

func TestConnMapGlobal_GetSetRemove(t *testing.T) {
	// Initialize the global connection map
	connMap = &connMapGlobal{Map: NewMapRedis()}

	// Define an AdapterName for testing
	name := aconns.AdapterName("test_adapter")

	// Mock the Redis connection pool
	pool := &redis.Pool{}

	// Test setting an AClientRedis
	client := &AClientRedis{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Name: name,
			},
		},
		pool: pool,
	}
	err := connMap.Set(client)
	assert.NoError(t, err)

	// Test getting the AClientRedis
	retrievedClient := connMap.Get(name)
	assert.NotNil(t, retrievedClient)
	assert.Equal(t, client, retrievedClient)

	// Test getting a non-existent AClientRedis
	nonExistentClient := connMap.Get(aconns.AdapterName("non_existent_adapter"))
	assert.Nil(t, nonExistentClient)

	// Test removing the AClientRedis
	connMap.Remove(name)
	removedClient := connMap.Get(name)
	assert.Nil(t, removedClient)
}

func TestConnMapGlobal_SetErrors(t *testing.T) {
	// Initialize the global connection map
	connMap = &connMapGlobal{Map: NewMapRedis()}

	// Test setting a nil AClientRedis
	err := connMap.Set(nil)
	assert.Error(t, err)
	assert.Equal(t, "connMapGlobal is nil", err.Error())

	// Test setting an AClientRedis with an empty name
	client := &AClientRedis{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Name: "",
			},
		},
	}
	err = connMap.Set(client)
	assert.Error(t, err)
	assert.Equal(t, "name is empty", err.Error())

	// Test setting an AClientRedis with a nil pool
	client = &AClientRedis{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Name: aconns.AdapterName("test_adapter"),
			},
		},
	}
	err = connMap.Set(client)
	assert.Error(t, err)
	assert.Equal(t, "pool is nil", err.Error())
}

func TestConnMapGlobal_GetRedisStore(t *testing.T) {
	// Initialize the global connection map
	connMap = &connMapGlobal{Map: NewMapRedis()}

	// Define an AdapterName for testing
	name := aconns.AdapterName("test_adapter")

	// Mock the Redis connection pool
	pool := &redis.Pool{}
	client := &AClientRedis{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Name: name,
			},
		},
		pool: pool,
	}
	err := connMap.Set(client)
	assert.NoError(t, err)

	// Test getting a RedisStore
	store := connMap.GetRedisStore(name, "test_prefix")
	assert.NotNil(t, store)

	// Test getting a RedisStore for a non-existent client
	nonExistentStore := connMap.GetRedisStore(aconns.AdapterName("non_existent_adapter"), "test_prefix")
	assert.Nil(t, nonExistentStore)
}

func TestNewRedisStoreSCS(t *testing.T) {
	// Initialize the global connection map
	connMap = &connMapGlobal{Map: NewMapRedis()}

	// Mock the Redis connection pool
	pool := &redis.Pool{}
	client := &AClientRedis{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Name: REDIS_MASTER,
			},
		},
		pool: pool,
	}
	err := connMap.Set(client)
	assert.NoError(t, err)

	// Test creating a new RedisStore for session management
	store := NewRedisStoreSCS()
	assert.NotNil(t, store)
}
