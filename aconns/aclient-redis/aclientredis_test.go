package aclient_redis

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/stretchr/testify/assert"
)

// Global variables for connection details
var (
	testHost     = "localhost"
	testPort     = 6379
	testDatabase = "0"
	testUser     = "testuser"
	testPassword = "testpass"
	testTimeout  = 5 * time.Second
)

func TestAClientRedis_Validate(t *testing.T) {
	redisClient := &AClientRedis{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("redis"),
				Name: aconns.AdapterName("test_redis"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		MaxIdleConnections: REDIS_MAXIDLE_CONNECTIONS,
		IdleTimeout:        REDIS_IDLETIMEOUT,
	}

	err := redisClient.Validate()
	assert.NoError(t, err)
}

func TestAClientRedis_Test(t *testing.T) {
	redisClient := &AClientRedis{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("redis"),
				Name: aconns.AdapterName("test_redis"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		MaxIdleConnections: REDIS_MAXIDLE_CONNECTIONS,
		IdleTimeout:        REDIS_IDLETIMEOUT,
	}

	ok, status, err := redisClient.Test()
	assert.False(t, ok)
	assert.Equal(t, aconns.TESTSTATUS_FAILED, status)
	assert.Error(t, err)
}

func TestAClientRedis_OpenConnection(t *testing.T) {
	redisClient := &AClientRedis{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("redis"),
				Name: aconns.AdapterName("test_redis"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		MaxIdleConnections: REDIS_MAXIDLE_CONNECTIONS,
		IdleTimeout:        REDIS_IDLETIMEOUT,
	}

	err := redisClient.OpenConnection()
	assert.Error(t, err)
}

func TestAClientRedis_CloseConnection(t *testing.T) {
	redisClient := &AClientRedis{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("redis"),
				Name: aconns.AdapterName("test_redis"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		MaxIdleConnections: REDIS_MAXIDLE_CONNECTIONS,
		IdleTimeout:        REDIS_IDLETIMEOUT,
	}

	// Mock the Redis connection pool
	pool := &redis.Pool{}
	redisClient.pool = pool

	err := redisClient.CloseConnection()
	assert.NoError(t, err)
	assert.Nil(t, redisClient.pool)
}

func TestAClientRedis_GetAddress(t *testing.T) {
	redisClient := &AClientRedis{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("redis"),
				Name: aconns.AdapterName("test_redis"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		MaxIdleConnections: REDIS_MAXIDLE_CONNECTIONS,
		IdleTimeout:        REDIS_IDLETIMEOUT,
	}

	address := redisClient.GetAddress()
	assert.Equal(t, "localhost:6379", address)
}

func TestAClientRedis_GetRedisStore(t *testing.T) {
	redisClient := &AClientRedis{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("redis"),
				Name: aconns.AdapterName("test_redis"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		MaxIdleConnections: REDIS_MAXIDLE_CONNECTIONS,
		IdleTimeout:        REDIS_IDLETIMEOUT,
	}

	// Mock the Redis connection pool
	pool := &redis.Pool{}
	redisClient.pool = pool

	store := redisClient.GetRedisStore("test_prefix")
	assert.NotNil(t, store)
}
