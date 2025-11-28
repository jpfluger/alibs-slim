package adb_pg

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// Global variables for connection details
var (
	testHost     = "localhost"
	testPort     = 5432
	testDatabase = "testdb"
	testUser     = "testuser"
	testPassword = "testpass"
	testTimeout  = 5
)

func TestADBPG_Validate(t *testing.T) {
	pg := &ADBPG{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("pg"),
				Name: aconns.AdapterName("test_pg"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		DialTimeout:  testTimeout,
		ReadTimeout:  testTimeout,
		WriteTimeout: testTimeout,
		PingTimeOut:  testTimeout,
	}

	err := pg.Validate()
	assert.NoError(t, err)
}

func TestADBPG_Test(t *testing.T) {
	pg := &ADBPG{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("pg"),
				Name: aconns.AdapterName("test_pg"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		DialTimeout:  testTimeout,
		ReadTimeout:  testTimeout,
		WriteTimeout: testTimeout,
		PingTimeOut:  testTimeout,
	}

	ok, status, err := pg.Test()
	assert.False(t, ok)
	assert.Equal(t, aconns.TESTSTATUS_FAILED, status)
	assert.Error(t, err)
}

func TestADBPG_OpenConnection(t *testing.T) {
	pg := &ADBPG{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("pg"),
				Name: aconns.AdapterName("test_pg"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		DialTimeout:  testTimeout,
		ReadTimeout:  testTimeout,
		WriteTimeout: testTimeout,
		PingTimeOut:  testTimeout,
	}

	err := pg.OpenConnection()
	assert.Error(t, err)
}

func TestADBPG_CloseConnection(t *testing.T) {
	pg := &ADBPG{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("pg"),
				Name: aconns.AdapterName("test_pg"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		DialTimeout:  testTimeout,
		ReadTimeout:  testTimeout,
		WriteTimeout: testTimeout,
		PingTimeOut:  testTimeout,
	}

	// Mock the database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Simulate opening a connection
	sqldb := sql.OpenDB(pg.getConnConfig())
	pg.db = bun.NewDB(sqldb, pgdialect.New())

	// Expect the close operation
	mock.ExpectClose()

	err = pg.CloseConnection()
	assert.NoError(t, err)
	assert.Nil(t, pg.db)
}

func TestADBPG_GetAddress(t *testing.T) {
	pg := &ADBPG{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("pg"),
				Name: aconns.AdapterName("test_pg"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		DialTimeout:  testTimeout,
		ReadTimeout:  testTimeout,
		WriteTimeout: testTimeout,
		PingTimeOut:  testTimeout,
	}

	address := pg.GetAddress()
	assert.Equal(t, "localhost:5432", address)
}

func TestADBPG_Count(t *testing.T) {
	pg := &ADBPG{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("pg"),
				Name: aconns.AdapterName("test_pg"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		DialTimeout:  testTimeout,
		ReadTimeout:  testTimeout,
		WriteTimeout: testTimeout,
		PingTimeOut:  testTimeout,
	}

	// Mock the database connection with exact query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	// Use the mocked sql.DB connection directly
	pg.db = bun.NewDB(db, pgdialect.New())

	// Define a model for the test
	type TestModel struct {
		ID int
	}

	// Mock the count query
	mock.ExpectQuery(`SELECT count(*) FROM "test_models" AS "test_model"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Use the model in the Count method
	count, err := pg.Count((*TestModel)(nil))
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestADBPG_SelectAll(t *testing.T) {
	pg := &ADBPG{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("pg"),
				Name: aconns.AdapterName("test_pg"),
				Host: testHost,
				Port: testPort,
			},
			Database: testDatabase,
			Username: testUser,
			Password: testPassword,
		},
		DialTimeout:  testTimeout,
		ReadTimeout:  testTimeout,
		WriteTimeout: testTimeout,
		PingTimeOut:  testTimeout,
	}

	// Mock the database connection with exact query matching
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	// Use the mocked sql.DB connection directly
	pg.db = bun.NewDB(db, pgdialect.New())

	// Define a model for the test
	type TestModel struct {
		ID int `bun:"id"`
	}

	// Mock the select query
	mock.ExpectQuery(`SELECT "test_model"."id" FROM "test_models" AS "test_model"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Initialize a slice to hold the results
	var models []TestModel

	// Use the model in the SelectAll method
	err = pg.SelectAll(&models)
	assert.NoError(t, err)
	assert.Len(t, models, 1)
	assert.Equal(t, 1, models[0].ID)
}
