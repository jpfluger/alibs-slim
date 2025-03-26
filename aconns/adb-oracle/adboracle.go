package adb_oracle

import (
	"database/sql"
	"fmt"
	"github.com/jpfluger/alibs-slim/aconns"
	go_ora "github.com/sijms/go-ora/v2" // go-ora is a pure Go Oracle driver.
	"strconv"
	"strings"
	"sync"
)

const (
	ADAPTERTYPE_ORACLE        aconns.AdapterType = "oracle"
	ORACLE_DEFAULT_PORT                          = 1521
	ORACLE_CONNECTION_TIMEOUT                    = 30
)

// ADBOracle represents an Oracle database adapter.
type ADBOracle struct {
	aconns.ADBAdapterBase

	Service           string `json:"service,omitempty"`
	ConnectionTimeout int    `json:"connectionTimeout,omitempty"`

	sqldb *sql.DB

	mu sync.RWMutex
}

// validate checks if the ADBOracle object is valid.
func (cn *ADBOracle) validate() error {
	if err := cn.ADBAdapterBase.Validate(); err != nil {
		if err != aconns.ErrDatabaseIsEmpty {
			return err
		}
	}

	cn.Service = strings.TrimSpace(cn.Service)
	if cn.Service == "" {
		return fmt.Errorf("service is empty")
	}

	if cn.ConnectionTimeout <= 0 {
		cn.ConnectionTimeout = ORACLE_CONNECTION_TIMEOUT
	}

	cn.Host = strings.TrimSpace(cn.Host)
	if cn.Host == "" {
		cn.Host = "localhost"
	}

	if cn.Port <= 0 {
		cn.Port = ORACLE_DEFAULT_PORT
	}

	return nil
}

// Validate checks if the ADBOracle object is valid.
func (cn *ADBOracle) Validate() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.validate()
}

// Test attempts to validate the ADBOracle, open a connection if necessary, and test the connection.
func (cn *ADBOracle) Test() (bool, aconns.TestStatus, error) {
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

// OpenConnection opens a connection to the Oracle database.
func (cn *ADBOracle) OpenConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.openConnection()
}

// openConnection opens a connection to the Oracle database.
func (cn *ADBOracle) openConnection() error {
	connString := cn.getConnString()

	sqldb, err := sql.Open("oracle", connString)
	if err != nil {
		return fmt.Errorf("could not open conn for oracle where host=%s; %v", cn.GetHost(), err)
	}

	if err = cn.testConnection(sqldb); err != nil {
		return err
	}

	cn.sqldb = sqldb
	return nil
}

// getConnString returns the connection string for the Oracle database.
func (cn *ADBOracle) getConnString() string {
	return go_ora.BuildUrl(cn.GetHost(), cn.GetPort(), cn.getService(), cn.GetUsername(), cn.GetPassword(), nil)
}

// GetAddress returns the address of the Oracle server.
func (cn *ADBOracle) GetAddress() string {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.getAddress()
}

// getAddress returns the address of the Oracle server.
func (cn *ADBOracle) getAddress() string {
	port := cn.GetPort()
	return fmt.Sprintf("%s:%s", cn.GetHost(), strconv.Itoa(port))
}

// GetService returns the service name of the Oracle database.
func (cn *ADBOracle) GetService() string {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.getService()
}

// getService returns the service name of the Oracle database.
func (cn *ADBOracle) getService() string {
	return cn.Service
}

// GetConnectionTimeout returns the connection timeout.
func (cn *ADBOracle) GetConnectionTimeout() int {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.getConnectionTimeout()
}

// getConnectionTimeout returns the connection timeout.
func (cn *ADBOracle) getConnectionTimeout() int {
	return cn.ConnectionTimeout
}

// testConnection tests the connection to the Oracle database.
func (cn *ADBOracle) testConnection(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("no oracle db has been created where host=%s", cn.GetHost())
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging db: %v", err)
	}
	return nil
}

// CloseConnection closes the connection to the Oracle database.
func (cn *ADBOracle) CloseConnection() error {
	if cn.sqldb != nil {
		if err := cn.sqldb.Close(); err != nil {
			return fmt.Errorf("error when closing oracle db connection where host=%s; %v", cn.GetHost(), err)
		}
		cn.sqldb = nil
	}
	return nil
}

// SQLDB returns the sql.DB instance.
func (cn *ADBOracle) SQLDB() *sql.DB {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.sqldb
}

// GetSandboxAdapter returns a sandbox adapter for the Oracle database.
// It ensures that the ADBOracle instance and its SQLDB are properly initialized.
func (cn *ADBOracle) GetSandboxAdapter() (aconns.ISBAdapter, error) {
	if cn == nil {
		return nil, fmt.Errorf("no oracle db has been created")
	}
	if cn.SQLDB() == nil {
		return nil, fmt.Errorf("no oracle db has been created where host=%s", cn.GetHost())
	}
	return aconns.NewSBAdapterSql(cn, cn.SQLDB()), nil
}

// ADBOracles represents a slice of ADBOracle pointers.
type ADBOracles []*ADBOracle
