package adb_oracle

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jpfluger/alibs-slim/aconns"
	go_ora "github.com/sijms/go-ora/v2" // go-ora is a pure Go Oracle driver.
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
	//TLSType           string `json:"tlsType,omitempty"` // "disable", "require", etc.

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

	//cn.TLSType = strings.TrimSpace(strings.ToLower(cn.TLSType))
	//if cn.TLSType == "" {
	//	cn.TLSType = "disable" // Default no TLS
	//}

	return nil
}

// Validate checks if the ADBOracle object is valid.
func (cn *ADBOracle) Validate() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.validate()
}

func (cn *ADBOracle) Test() (bool, aconns.TestStatus, error) {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.test()
}

func (cn *ADBOracle) test() (bool, aconns.TestStatus, error) {
	if err := cn.validate(); err != nil {
		cn.UpdateHealth(aconns.HEALTHSTATUS_VALIDATE_FAILED)
		return false, aconns.TESTSTATUS_FAILED, err
	}

	if cn.sqldb == nil {
		if err := cn.openConnection(); err != nil {
			cn.UpdateHealth(aconns.HEALTHSTATUS_OPEN_FAILED)
			return false, aconns.TESTSTATUS_FAILED, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Ping timeout
	defer cancel()
	if err := cn.sqldb.PingContext(ctx); err != nil {
		status := aconns.HEALTHSTATUS_PING_FAILED
		if context.DeadlineExceeded == err {
			status = aconns.HEALTHSTATUS_TIMEOUT
		} else if strings.Contains(err.Error(), "network") || strings.Contains(err.Error(), "connection refused") {
			status = aconns.HEALTHSTATUS_NETWORK_ERROR
		}
		cn.UpdateHealth(status)
		//alog.LOGGER(alog.LOGGER_APP).Warn().Err(err).Msg("Oracle ping failed")
		return false, aconns.TESTSTATUS_FAILED, err
	}

	cn.UpdateHealth(aconns.HEALTHSTATUS_HEALTHY)
	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

// OpenConnection opens a connection to the Oracle database.
func (cn *ADBOracle) OpenConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.openConnection()
}

func (cn *ADBOracle) openConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cn.ConnectionTimeout)*time.Second)
	defer cancel()

	connString := cn.getConnString() // Add TLS params to BuildUrl if needed (go-ora supports)

	//var tlsConfig *tls.Config
	//switch cn.TLSType {
	//case "require":
	//	tlsConfig = &tls.Config{} // Basic; add certs if needed
	//case "verify-ca":
	//	// Load CA certs, etc.
	//}

	// If go-ora supports TLS, integrate; otherwise, note limitation.

	sqldb, err := sql.Open("oracle", connString)
	if err != nil {
		return fmt.Errorf("could not open conn for oracle where host=%s; %v", cn.GetHost(), err)
	}

	if err = sqldb.PingContext(ctx); err != nil { // Use context for timeout
		sqldb.Close()
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

// CloseConnection closes the connection to the Oracle database.
func (cn *ADBOracle) CloseConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	if cn.sqldb != nil {
		if err := cn.sqldb.Close(); err != nil {
			return fmt.Errorf("error when closing oracle db connection where host=%s; %v", cn.GetHost(), err)
		}
		cn.sqldb = nil
		cn.UpdateHealth(aconns.HEALTHSTATUS_CLOSED)
	}
	return nil
}

func (cn *ADBOracle) SQLDB() *sql.DB {
	cn.mu.RLock()
	if cn.IsHealthy() && !cn.GetHealth().IsStale(5*time.Minute) {
		defer cn.mu.RUnlock()
		return cn.sqldb
	}
	cn.mu.RUnlock()

	// Upgrade to write lock for refresh
	cn.test() // Refresh and test
	return cn.sqldb
}

func (cn *ADBOracle) Refresh() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	if cn.sqldb != nil {
		cn.sqldb.Close() // Close old
		cn.sqldb = nil
	}
	return cn.openConnection()
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
