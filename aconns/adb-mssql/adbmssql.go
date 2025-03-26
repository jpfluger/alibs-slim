package adb_mssql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb" // MSSQL driver for Go's database/sql package.
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/uptrace/bun"                      // Bun is a SQL-first Golang ORM for PostgreSQL, MySQL, MSSQL, and SQLite.
	"github.com/uptrace/bun/dialect/mssqldialect" // MSSQL dialect for Bun ORM.
	"net/url"
	"strconv"
	"strings"
	"sync"
)

const (
	ADAPTERTYPE_MSSQL                 aconns.AdapterType = "mssql"
	MSSQL_DEFAULT_PORT                                   = 1433
	MSSQL_CONNECTION_TIMEOUT                             = 30
	MSSQL_ENCRYPT_FALSE                                  = "false"
	MSSQL_TRUSTSERVERCERTIFICATE_TRUE                    = "true"
)

// ADBMSSql represents an MSSQL database adapter.
type ADBMSSql struct {
	aconns.ADBAdapterBase

	ConnectionTimeout      int    `json:"connectionTimeout,omitempty"`
	Encrypt                string `json:"encrypt,omitempty"`
	TrustServerCertificate string `json:"trustServerCertificate,omitempty"`

	sqldb *sql.DB
	db    *bun.DB

	mu sync.RWMutex
}

// validate checks if the ADBMSSql object is valid.
func (cn *ADBMSSql) validate() error {
	if err := cn.ADBAdapterBase.Validate(); err != nil {
		return err
	}

	if cn.ConnectionTimeout <= 0 {
		cn.ConnectionTimeout = MSSQL_CONNECTION_TIMEOUT
	}

	cn.Host = strings.TrimSpace(cn.Host)
	if cn.Host == "" {
		cn.Host = "localhost"
	}

	if cn.Port <= 0 {
		cn.Port = MSSQL_DEFAULT_PORT
	}

	cn.Encrypt = strings.TrimSpace(cn.Encrypt)
	if cn.Encrypt == "" {
		cn.Encrypt = MSSQL_ENCRYPT_FALSE
	}

	cn.TrustServerCertificate = strings.TrimSpace(cn.TrustServerCertificate)
	if cn.TrustServerCertificate == "" {
		cn.TrustServerCertificate = MSSQL_TRUSTSERVERCERTIFICATE_TRUE
	}

	return nil
}

// Validate checks if the ADBMSSql object is valid.
func (cn *ADBMSSql) Validate() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.validate()
}

// Test attempts to validate the ADBMSSql, open a connection if necessary, and test the connection.
func (cn *ADBMSSql) Test() (bool, aconns.TestStatus, error) {
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

// OpenConnection opens a connection to the MSSQL database.
func (cn *ADBMSSql) OpenConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.openConnection()
}

// openConnection opens a connection to the MSSQL database.
func (cn *ADBMSSql) openConnection() error {
	if err := cn.validate(); err != nil {
		return err
	}

	connector, err := cn.getConnConfig()
	if err != nil {
		return err
	}

	sqldb := sql.OpenDB(connector)
	if err = cn.testConnection(sqldb); err != nil {
		return err
	}

	cn.sqldb = sqldb
	cn.db = bun.NewDB(sqldb, mssqldialect.New())
	return nil
}

// getConnString generates the connection string for the MSSQL database.
func (cn *ADBMSSql) getConnString() string {
	query := url.Values{}
	query.Add("database", cn.Database)
	query.Add("encrypt", fmt.Sprintf("%v", cn.Encrypt))
	query.Add("connection timeout", fmt.Sprintf("%d", cn.ConnectionTimeout))

	if cn.TrustServerCertificate == MSSQL_TRUSTSERVERCERTIFICATE_TRUE {
		query.Add("TrustServerCertificate", fmt.Sprintf("%v", cn.TrustServerCertificate))
	}

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(cn.Username, cn.Password),
		Host:     fmt.Sprintf("%s:%d", cn.Host, cn.Port),
		RawQuery: query.Encode(),
	}

	return u.String()
}

// getConnConfig generates the connection configuration for the MSSQL database.
func (cn *ADBMSSql) getConnConfig() (driver.Connector, error) {
	connString := cn.getConnString()
	return mssql.NewConnector(connString)
}

// GetAddress returns the address of the MSSQL server.
func (cn *ADBMSSql) GetAddress() string {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.getAddress()
}

// getAddress returns the address of the MSSQL server.
func (cn *ADBMSSql) getAddress() string {
	port := cn.GetPort()
	return fmt.Sprintf("%s:%s", cn.GetHost(), strconv.Itoa(port))
}

// GetConnectionTimeout returns the connection timeout.
func (cn *ADBMSSql) GetConnectionTimeout() int {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.getConnectionTimeout()
}

// getConnectionTimeout returns the connection timeout.
func (cn *ADBMSSql) getConnectionTimeout() int {
	return cn.ConnectionTimeout
}

// GetEncrypt returns the encryption setting.
func (cn *ADBMSSql) GetEncrypt() string {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.getEncrypt()
}

// getEncrypt returns the encryption setting.
func (cn *ADBMSSql) getEncrypt() string {
	return cn.Encrypt
}

// testConnection tests the connection to the MSSQL database.
func (cn *ADBMSSql) testConnection(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("no mssql db has been created where host=%s", cn.GetHost())
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging db: %v", err)
	}
	return nil
}

// CloseConnection closes the MSSQL connection.
func (cn *ADBMSSql) CloseConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	var err error

	if cn.db != nil {
		if closeErr := cn.db.Close(); closeErr != nil {
			err = fmt.Errorf("error when closing bun db connection where host=%s; %v", cn.GetHost(), closeErr)
		}
		cn.db = nil
	}
	if cn.sqldb != nil {
		if closeErr := cn.sqldb.Close(); closeErr != nil {
			if err != nil {
				err = fmt.Errorf("%v; %v", err, closeErr)
			} else {
				err = fmt.Errorf("error when closing mssql db connection where host=%s; %v", cn.GetHost(), closeErr)
			}
		}
		cn.sqldb = nil
	}

	return err
}

// DB returns the bun.DB instance.
func (cn *ADBMSSql) DB() *bun.DB {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.db
}

// SQLDB returns the sql.DB instance.
func (cn *ADBMSSql) SQLDB() *sql.DB {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.sqldb
}

// GetSandboxAdapter returns a sandbox adapter for the MSSQL database.
// It ensures that the ADBMSSql instance and its SQLDB are properly initialized.
func (cn *ADBMSSql) GetSandboxAdapter() (aconns.ISBAdapter, error) {
	if cn == nil {
		return nil, fmt.Errorf("no mssql db has been created")
	}
	if cn.SQLDB() == nil {
		return nil, fmt.Errorf("no mssql db has been created where host=%s", cn.GetHost())
	}
	return aconns.NewSBAdapterSql(cn, cn.SQLDB()), nil
}

// ADBMSSqls represents a slice of ADBMSSql pointers.
type ADBMSSqls []*ADBMSSql
