package adb_mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // MySQL driver for Go's database/sql package.
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/uptrace/bun"                      // Bun is a SQL-first Golang ORM for PostgreSQL, MySQL, MSSQL, and SQLite.
	"github.com/uptrace/bun/dialect/mysqldialect" // MySQL dialect for Bun ORM.
	"strconv"
	"strings"
	"sync"
)

const (
	ADAPTERTYPE_MYSQL        aconns.AdapterType = "mysql"
	ADAPTERTYPE_MARIA        aconns.AdapterType = "maria"
	MYSQL_DEFAULT_PORT                          = 3306
	MYSQL_CONNECTION_TIMEOUT                    = 30
)

// ADBMysql represents a MySQL database adapter.
type ADBMysql struct {
	aconns.ADBAdapterBase

	ConnectionTimeout int `json:"connectionTimeout,omitempty"`

	sqldb *sql.DB
	db    *bun.DB

	mu sync.RWMutex
}

// validate checks if the ADBMysql object is valid.
func (cn *ADBMysql) validate() error {
	if err := cn.ADBAdapterBase.Validate(); err != nil {
		return err
	}

	if cn.ConnectionTimeout <= 0 {
		cn.ConnectionTimeout = MYSQL_CONNECTION_TIMEOUT
	}

	cn.Host = strings.TrimSpace(cn.Host)
	if cn.Host == "" {
		cn.Host = "localhost"
	}

	if cn.Port <= 0 {
		cn.Port = MYSQL_DEFAULT_PORT
	}

	return nil
}

// Validate checks if the ADBMysql object is valid.
func (cn *ADBMysql) Validate() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.validate()
}

// Test attempts to validate the ADBMysql, open a connection if necessary, and test the connection.
func (cn *ADBMysql) Test() (bool, aconns.TestStatus, error) {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	if err := cn.validate(); err != nil {
		return false, aconns.TESTSTATUS_FAILED, err
	}

	if cn.sqldb != nil {
		if err := cn.testConnection(cn.sqldb); err == nil {
			return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
		}
	}

	if err := cn.openConnection(); err != nil {
		return false, aconns.TESTSTATUS_FAILED, err
	}

	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

// OpenConnection opens a connection to the MySQL database.
func (cn *ADBMysql) OpenConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.openConnection()
}

// openConnection opens a connection to the MySQL database.
func (cn *ADBMysql) openConnection() error {
	connString := cn.getConnString()

	sqldb, err := sql.Open("mysql", connString)
	if err != nil {
		return fmt.Errorf("could not open conn for mysql where host=%s; %v", cn.GetHost(), err)
	}

	if err = cn.testConnection(sqldb); err != nil {
		return err
	}

	cn.sqldb = sqldb
	cn.db = bun.NewDB(sqldb, mysqldialect.New())

	return nil
}

// getConnString returns the connection string for the MySQL database.
func (cn *ADBMysql) getConnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%ds",
		cn.GetUsername(), cn.GetPassword(), cn.GetHost(), cn.GetPort(), cn.GetDatabase(), cn.getConnectionTimeout())
}

// GetAddress returns the address of the MySQL server.
func (cn *ADBMysql) GetAddress() string {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.getAddress()
}

// getAddress returns the address of the MySQL server.
func (cn *ADBMysql) getAddress() string {
	port := cn.GetPort()
	return fmt.Sprintf("%s:%s", cn.GetHost(), strconv.Itoa(port))
}

// GetConnectionTimeout returns the connection timeout.
func (cn *ADBMysql) GetConnectionTimeout() int {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.getConnectionTimeout()
}

// getConnectionTimeout returns the connection timeout.
func (cn *ADBMysql) getConnectionTimeout() int {
	return cn.ConnectionTimeout
}

// testConnection tests the connection to the MySQL database.
func (cn *ADBMysql) testConnection(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("no mysql db has been created where host=%s", cn.GetHost())
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging db: %v", err)
	}
	return nil
}

// CloseConnection closes the connection to the MySQL database.
func (cn *ADBMysql) CloseConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	// Close the bun.DB connection first
	if cn.db != nil {
		if err := cn.db.Close(); err != nil {
			return fmt.Errorf("error when closing bun db connection where host=%s; %v", cn.GetHost(), err)
		}
		cn.db = nil
	}

	// Close the sql.DB connection
	if cn.sqldb != nil {
		if err := cn.sqldb.Close(); err != nil {
			return fmt.Errorf("error when closing sql db connection where host=%s; %v", cn.GetHost(), err)
		}
		cn.sqldb = nil
	}

	return nil
}

// DB returns the bun.DB instance.
func (cn *ADBMysql) DB() *bun.DB {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.db
}

// SQLDB returns the sql.DB instance.
func (cn *ADBMysql) SQLDB() *sql.DB {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.sqldb
}

// GetSandboxAdapter returns a sandbox adapter for the MySQL database.
// It ensures that the ADBMysql instance and its SQLDB are properly initialized.
func (cn *ADBMysql) GetSandboxAdapter() (aconns.ISBAdapter, error) {
	if cn == nil {
		return nil, fmt.Errorf("no mysql db has been created")
	}
	if cn.SQLDB() == nil {
		return nil, fmt.Errorf("no mysql db has been created where host=%s", cn.GetHost())
	}
	return aconns.NewSBAdapterSql(cn, cn.SQLDB()), nil
}

// ADBMysqls represents a slice of ADBMysql pointers.
type ADBMysqls []*ADBMysql
