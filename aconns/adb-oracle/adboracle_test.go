package adb_oracle

import (
	"database/sql"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Global variables for connection details
var (
	testHost     = "localhost"
	testPort     = 1521
	testDatabase = "testdb"
	testUser     = "testuser"
	testPassword = "testpass"
	testService  = "testservice"
	testTimeout  = 30
)

func TestADBOracle_Validate(t *testing.T) {
	oracle := &ADBOracle{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("oracle"),
				Name: aconns.AdapterName("test_oracle"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		Service:           testService,
		ConnectionTimeout: testTimeout,
	}

	err := oracle.Validate()
	assert.NoError(t, err)
}

func TestADBOracle_Test(t *testing.T) {
	oracle := &ADBOracle{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("oracle"),
				Name: aconns.AdapterName("test_oracle"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		Service:           testService,
		ConnectionTimeout: testTimeout,
	}

	ok, status, err := oracle.Test()
	assert.False(t, ok)
	assert.Equal(t, aconns.TESTSTATUS_FAILED, status)
	assert.Error(t, err)
}

func TestADBOracle_OpenConnection(t *testing.T) {
	oracle := &ADBOracle{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("oracle"),
				Name: aconns.AdapterName("test_oracle"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		Service:           testService,
		ConnectionTimeout: testTimeout,
	}

	err := oracle.OpenConnection()
	assert.Error(t, err)
}

func TestADBOracle_CloseConnection(t *testing.T) {
	oracle := &ADBOracle{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("oracle"),
				Name: aconns.AdapterName("test_oracle"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		Service:           testService,
		ConnectionTimeout: testTimeout,
	}

	// Simulate opening a connection
	oracle.sqldb, _ = sql.Open("oracle", oracle.getConnString())

	err := oracle.CloseConnection()
	assert.NoError(t, err)
	assert.Nil(t, oracle.sqldb)
}

func TestADBOracle_GetAddress(t *testing.T) {
	oracle := &ADBOracle{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("oracle"),
				Name: aconns.AdapterName("test_oracle"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		Service:           testService,
		ConnectionTimeout: testTimeout,
	}

	address := oracle.GetAddress()
	assert.Equal(t, "localhost:1521", address)
}

func TestADBOracle_GetService(t *testing.T) {
	oracle := &ADBOracle{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("oracle"),
				Name: aconns.AdapterName("test_oracle"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		Service:           testService,
		ConnectionTimeout: testTimeout,
	}

	service := oracle.GetService()
	assert.Equal(t, testService, service)
}

func TestADBOracle_GetConnectionTimeout(t *testing.T) {
	oracle := &ADBOracle{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("oracle"),
				Name: aconns.AdapterName("test_oracle"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		Service:           testService,
		ConnectionTimeout: testTimeout,
	}

	timeout := oracle.GetConnectionTimeout()
	assert.Equal(t, testTimeout, timeout)
}
