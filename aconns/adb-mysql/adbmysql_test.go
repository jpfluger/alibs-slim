package adb_mysql

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"

	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/stretchr/testify/assert"
)

// Global variables for connection details
var (
	testHost     = "localhost"
	testPort     = 3306
	testDatabase = "testdb"
	testUser     = "testuser"
	testPassword = "testpass"
	testTimeout  = 30
)

func TestADBMysql_Validate(t *testing.T) {
	mysql := &ADBMysql{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("mysql"),
				Name: aconns.AdapterName("test_mysql"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		ConnectionTimeout: testTimeout,
	}

	err := mysql.Validate()
	assert.NoError(t, err)
}

func TestADBMysql_Test(t *testing.T) {
	mysql := &ADBMysql{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("mysql"),
				Name: aconns.AdapterName("test_mysql"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		ConnectionTimeout: testTimeout,
	}

	ok, status, err := mysql.Test()
	assert.False(t, ok)
	assert.Equal(t, aconns.TESTSTATUS_FAILED, status)
	assert.Error(t, err)
}

func TestADBMysql_CloseConnection(t *testing.T) {
	mysql := &ADBMysql{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("mysql"),
				Name: aconns.AdapterName("test_mysql"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		ConnectionTimeout: testTimeout,
	}

	// Mock the database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Expect the version query
	mock.ExpectQuery("SELECT version()").WillReturnRows(sqlmock.NewRows([]string{"version()"}).AddRow("8.0.23"))

	// Simulate opening a connection
	mysql.sqldb = db
	mysql.db = bun.NewDB(db, mysqldialect.New())

	// Expect the close operation
	mock.ExpectClose()

	err = mysql.CloseConnection()
	assert.NoError(t, err)
	assert.Nil(t, mysql.sqldb)
	assert.Nil(t, mysql.db)
}

func TestADBMysql_GetAddress(t *testing.T) {
	mysql := &ADBMysql{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("mysql"),
				Name: aconns.AdapterName("test_mysql"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		ConnectionTimeout: testTimeout,
	}

	address := mysql.GetAddress()
	assert.Equal(t, "localhost:3306", address)
}

func TestADBMysql_GetConnectionTimeout(t *testing.T) {
	mysql := &ADBMysql{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("mysql"),
				Name: aconns.AdapterName("test_mysql"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		ConnectionTimeout: testTimeout,
	}

	timeout := mysql.GetConnectionTimeout()
	assert.Equal(t, testTimeout, timeout)
}
