package adb_mssql

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
)

// Global variables for connection details
var (
	testHost     = "localhost"
	testPort     = 1433
	testDatabase = "testdb"
	testUser     = "testuser"
	testPassword = "testpass"
	testService  = "testservice"
	testTimeout  = 30
	testEncrypt  = "false"
)

func TestADBMSSql_Validate(t *testing.T) {
	mssql := &ADBMSSql{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("mssql"),
				Name: aconns.AdapterName("test_mssql"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		ConnectionTimeout: testTimeout,
		Encrypt:           testEncrypt,
	}

	err := mssql.Validate()
	assert.NoError(t, err)
}

func TestADBMSSql_Test(t *testing.T) {
	mssql := &ADBMSSql{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("mssql"),
				Name: aconns.AdapterName("test_mssql"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		ConnectionTimeout: testTimeout,
		Encrypt:           testEncrypt,
	}

	ok, status, err := mssql.Test()
	assert.False(t, ok)
	assert.Equal(t, aconns.TESTSTATUS_FAILED, status)
	assert.Error(t, err)
}

func TestADBMSSql_OpenConnection(t *testing.T) {
	mssql := &ADBMSSql{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("mssql"),
				Name: aconns.AdapterName("test_mssql"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		ConnectionTimeout: testTimeout,
		Encrypt:           testEncrypt,
	}

	err := mssql.OpenConnection()
	assert.Error(t, err)
}

func TestADBMSSql_CloseConnection(t *testing.T) {
	mssql := &ADBMSSql{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("mssql"),
				Name: aconns.AdapterName("test_mssql"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		ConnectionTimeout: testTimeout,
		Encrypt:           testEncrypt,
	}

	// Mock the database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Simulate opening a connection
	mssql.sqldb = db
	mssql.db = bun.NewDB(db, mssqldialect.New())

	// Perform a query to meet the expectation
	mock.ExpectQuery("SELECT version()").WillReturnRows(sqlmock.NewRows([]string{"version()"}).AddRow("8.0.23"))
	rows, err := mssql.db.Query("SELECT version()")
	assert.NoError(t, err)
	assert.NoError(t, rows.Close())

	// Expect the close operation
	mock.ExpectClose()

	err = mssql.CloseConnection()
	assert.NoError(t, err)
	assert.Nil(t, mssql.sqldb)
	assert.Nil(t, mssql.db)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestADBMSSql_GetAddress(t *testing.T) {
	mssql := &ADBMSSql{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("mssql"),
				Name: aconns.AdapterName("test_mssql"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		ConnectionTimeout: testTimeout,
		Encrypt:           testEncrypt,
	}

	address := mssql.GetAddress()
	assert.Equal(t, "localhost:1433", address)
}

func TestADBMSSql_GetConnectionTimeout(t *testing.T) {
	mssql := &ADBMSSql{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("mssql"),
				Name: aconns.AdapterName("test_mssql"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		ConnectionTimeout: testTimeout,
		Encrypt:           testEncrypt,
	}

	timeout := mssql.GetConnectionTimeout()
	assert.Equal(t, testTimeout, timeout)
}

func TestADBMSSql_GetEncrypt(t *testing.T) {
	mssql := &ADBMSSql{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("mssql"),
				Name: aconns.AdapterName("test_mssql"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		ConnectionTimeout: testTimeout,
		Encrypt:           testEncrypt,
	}

	encrypt := mssql.GetEncrypt()
	assert.Equal(t, testEncrypt, encrypt)
}
