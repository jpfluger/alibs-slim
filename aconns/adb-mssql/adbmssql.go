package adb_mssql

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb" // MSSQL driver for Go's database/sql package.
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/uptrace/bun"                      // Bun is a SQL-first Golang ORM for PostgreSQL, MySQL, MSSQL, and SQLite.
	"github.com/uptrace/bun/dialect/mssqldialect" // MSSQL dialect for Bun ORM.
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
	return cn.test()
}

// Test attempts to validate the ADBMSSql, open a connection if necessary, and test the connection.
func (cn *ADBMSSql) test() (bool, aconns.TestStatus, error) {
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
	if err := cn.testConnectionWithCtx(ctx, cn.sqldb); err != nil {
		status := aconns.HEALTHSTATUS_PING_FAILED
		if context.DeadlineExceeded == err {
			status = aconns.HEALTHSTATUS_TIMEOUT
		} else if strings.Contains(err.Error(), "network") || strings.Contains(err.Error(), "connection refused") {
			status = aconns.HEALTHSTATUS_NETWORK_ERROR
		}
		cn.UpdateHealth(status)
		//alog.LOGGER(alog.LOGGER_APP).Warn().Err(err).Msg("MSSQL test failed")
		return false, aconns.TESTSTATUS_FAILED, err
	}

	cn.UpdateHealth(aconns.HEALTHSTATUS_HEALTHY)
	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

// Refresh refreshes the MSSQL connection by closing the existing one (if any) and opening a new one.
func (cn *ADBMSSql) Refresh() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	if cn.db != nil {
		cn.db.Close()
		cn.db = nil
	}
	if cn.sqldb != nil {
		cn.sqldb.Close()
		cn.sqldb = nil
	}
	return cn.openConnection()
}

// openConnection opens a connection to the MSSQL database.
func (cn *ADBMSSql) openConnection() error {
	connString := cn.getConnString()

	sqldb, err := sql.Open("sqlserver", connString)
	if err != nil {
		return fmt.Errorf("could not open conn for mssql where host=%s; %v", cn.GetHost(), err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cn.ConnectionTimeout)*time.Second)
	defer cancel()
	if err = sqldb.PingContext(ctx); err != nil {
		sqldb.Close()
		return fmt.Errorf("could not ping new conn where host=%s; %v", cn.GetHost(), err)
	}

	db := bun.NewDB(sqldb, mssqldialect.New())

	if err := cn.testConnectionWithCtx(ctx, sqldb); err != nil {
		return err
	}

	cn.sqldb = sqldb
	cn.db = db
	return nil
}

// getConnString returns the connection string for the MSSQL database.
func (cn *ADBMSSql) getConnString() string {
	query := url.Values{}
	query.Add("database", cn.GetDatabase())
	query.Add("encrypt", cn.getEncrypt())
	query.Add("TrustServerCertificate", cn.TrustServerCertificate)
	query.Add("connection timeout", strconv.Itoa(cn.getConnectionTimeout()))

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(cn.GetUsername(), cn.GetPassword()),
		Host:     cn.getAddress(),
		RawQuery: query.Encode(),
	}
	return u.String()
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

// testConnectionWithCtx tests the connection to the MSSQL database using a provided context.
func (cn *ADBMSSql) testConnectionWithCtx(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("no mssql db has been created where host=%s", cn.GetHost())
	}
	if err := db.PingContext(ctx); err != nil {
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
	cn.UpdateHealth(aconns.HEALTHSTATUS_CLOSED)

	return err
}

// DB returns the bun.DB instance.
func (cn *ADBMSSql) DB() *bun.DB {
	cn.mu.RLock()
	if cn.IsHealthy() && !cn.GetHealth().IsStale(5*time.Minute) {
		defer cn.mu.RUnlock()
		return cn.db
	}
	cn.mu.RUnlock()

	// Upgrade to write lock for refresh
	cn.mu.Lock()
	defer cn.mu.Unlock()
	cn.test() // Refresh and test
	return cn.db
}

// SQLDB returns the sql.DB instance.
func (cn *ADBMSSql) SQLDB() *sql.DB {
	cn.mu.RLock()
	if cn.IsHealthy() && !cn.GetHealth().IsStale(5*time.Minute) {
		defer cn.mu.RUnlock()
		return cn.sqldb
	}
	cn.mu.RUnlock()

	// Upgrade to write lock for refresh
	cn.mu.Lock()
	defer cn.mu.Unlock()
	cn.test() // Refresh and test
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

//package adb_mssql
//
//import (
//	"database/sql"
//	"database/sql/driver"
//	"fmt"
//	"net/url"
//	"strconv"
//	"strings"
//	"sync"
//
//	_ "github.com/denisenkom/go-mssqldb" // MSSQL driver for Go's database/sql package.
//	mssql "github.com/denisenkom/go-mssqldb"
//	"github.com/jpfluger/alibs-slim/aconns"
//	"github.com/uptrace/bun"                      // Bun is a SQL-first Golang ORM for PostgreSQL, MySQL, MSSQL, and SQLite.
//	"github.com/uptrace/bun/dialect/mssqldialect" // MSSQL dialect for Bun ORM.
//)
//
//const (
//	ADAPTERTYPE_MSSQL                 aconns.AdapterType = "mssql"
//	MSSQL_DEFAULT_PORT                                   = 1433
//	MSSQL_CONNECTION_TIMEOUT                             = 30
//	MSSQL_ENCRYPT_FALSE                                  = "false"
//	MSSQL_TRUSTSERVERCERTIFICATE_TRUE                    = "true"
//)
//
//// ADBMSSql represents an MSSQL database adapter.
//type ADBMSSql struct {
//	aconns.ADBAdapterBase
//
//	ConnectionTimeout      int    `json:"connectionTimeout,omitempty"`
//	Encrypt                string `json:"encrypt,omitempty"`
//	TrustServerCertificate string `json:"trustServerCertificate,omitempty"`
//
//	sqldb *sql.DB
//	db    *bun.DB
//
//	mu sync.RWMutex
//}
//
//// validate checks if the ADBMSSql object is valid.
//func (cn *ADBMSSql) validate() error {
//	if err := cn.ADBAdapterBase.Validate(); err != nil {
//		return err
//	}
//
//	if cn.ConnectionTimeout <= 0 {
//		cn.ConnectionTimeout = MSSQL_CONNECTION_TIMEOUT
//	}
//
//	cn.Host = strings.TrimSpace(cn.Host)
//	if cn.Host == "" {
//		cn.Host = "localhost"
//	}
//
//	if cn.Port <= 0 {
//		cn.Port = MSSQL_DEFAULT_PORT
//	}
//
//	cn.Encrypt = strings.TrimSpace(cn.Encrypt)
//	if cn.Encrypt == "" {
//		cn.Encrypt = MSSQL_ENCRYPT_FALSE
//	}
//
//	cn.TrustServerCertificate = strings.TrimSpace(cn.TrustServerCertificate)
//	if cn.TrustServerCertificate == "" {
//		cn.TrustServerCertificate = MSSQL_TRUSTSERVERCERTIFICATE_TRUE
//	}
//
//	return nil
//}
//
//// Validate checks if the ADBMSSql object is valid.
//func (cn *ADBMSSql) Validate() error {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//	return cn.validate()
//}
//
//// Test attempts to validate the ADBMSSql, open a connection if necessary, and test the connection.
//func (cn *ADBMSSql) Test() (bool, aconns.TestStatus, error) {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//
//	if err := cn.validate(); err != nil {
//		return false, aconns.TESTSTATUS_FAILED, err
//	}
//
//	if cn.sqldb != nil {
//		if err := cn.testConnection(cn.sqldb); err == nil {
//			return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
//		}
//	}
//
//	if err := cn.openConnection(); err != nil {
//		return false, aconns.TESTSTATUS_FAILED, err
//	}
//
//	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
//}
//
//// OpenConnection opens a connection to the MSSQL database.
//func (cn *ADBMSSql) OpenConnection() error {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//	return cn.openConnection()
//}
//
//// openConnection opens a connection to the MSSQL database.
//func (cn *ADBMSSql) openConnection() error {
//	if err := cn.validate(); err != nil {
//		return err
//	}
//
//	connector, err := cn.getConnConfig()
//	if err != nil {
//		return err
//	}
//
//	sqldb := sql.OpenDB(connector)
//	if err = cn.testConnection(sqldb); err != nil {
//		return err
//	}
//
//	cn.sqldb = sqldb
//	cn.db = bun.NewDB(sqldb, mssqldialect.New())
//	return nil
//}
//
//// getConnString generates the connection string for the MSSQL database.
//func (cn *ADBMSSql) getConnString() string {
//	query := url.Values{}
//	query.Add("database", cn.Database)
//	query.Add("encrypt", fmt.Sprintf("%v", cn.Encrypt))
//	query.Add("connection timeout", fmt.Sprintf("%d", cn.ConnectionTimeout))
//
//	if cn.TrustServerCertificate == MSSQL_TRUSTSERVERCERTIFICATE_TRUE {
//		query.Add("TrustServerCertificate", fmt.Sprintf("%v", cn.TrustServerCertificate))
//	}
//
//	u := &url.URL{
//		Scheme:   "sqlserver",
//		User:     url.UserPassword(cn.Username, cn.Password),
//		Host:     fmt.Sprintf("%s:%d", cn.Host, cn.Port),
//		RawQuery: query.Encode(),
//	}
//
//	return u.String()
//}
//
//// getConnConfig generates the connection configuration for the MSSQL database.
//func (cn *ADBMSSql) getConnConfig() (driver.Connector, error) {
//	connString := cn.getConnString()
//	return mssql.NewConnector(connString)
//}
//
//// GetAddress returns the address of the MSSQL server.
//func (cn *ADBMSSql) GetAddress() string {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//	return cn.getAddress()
//}
//
//// getAddress returns the address of the MSSQL server.
//func (cn *ADBMSSql) getAddress() string {
//	port := cn.GetPort()
//	return fmt.Sprintf("%s:%s", cn.GetHost(), strconv.Itoa(port))
//}
//
//// GetConnectionTimeout returns the connection timeout.
//func (cn *ADBMSSql) GetConnectionTimeout() int {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//	return cn.getConnectionTimeout()
//}
//
//// getConnectionTimeout returns the connection timeout.
//func (cn *ADBMSSql) getConnectionTimeout() int {
//	return cn.ConnectionTimeout
//}
//
//// GetEncrypt returns the encryption setting.
//func (cn *ADBMSSql) GetEncrypt() string {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//	return cn.getEncrypt()
//}
//
//// getEncrypt returns the encryption setting.
//func (cn *ADBMSSql) getEncrypt() string {
//	return cn.Encrypt
//}
//
//// testConnection tests the connection to the MSSQL database.
//func (cn *ADBMSSql) testConnection(db *sql.DB) error {
//	if db == nil {
//		return fmt.Errorf("no mssql db has been created where host=%s", cn.GetHost())
//	}
//	if err := db.Ping(); err != nil {
//		return fmt.Errorf("error pinging db: %v", err)
//	}
//	return nil
//}
//
//// CloseConnection closes the MSSQL connection.
//func (cn *ADBMSSql) CloseConnection() error {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//
//	var err error
//
//	if cn.db != nil {
//		if closeErr := cn.db.Close(); closeErr != nil {
//			err = fmt.Errorf("error when closing bun db connection where host=%s; %v", cn.GetHost(), closeErr)
//		}
//		cn.db = nil
//	}
//	if cn.sqldb != nil {
//		if closeErr := cn.sqldb.Close(); closeErr != nil {
//			if err != nil {
//				err = fmt.Errorf("%v; %v", err, closeErr)
//			} else {
//				err = fmt.Errorf("error when closing mssql db connection where host=%s; %v", cn.GetHost(), closeErr)
//			}
//		}
//		cn.sqldb = nil
//	}
//
//	return err
//}
//
//// DB returns the bun.DB instance.
//func (cn *ADBMSSql) DB() *bun.DB {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//	return cn.db
//}
//
//// SQLDB returns the sql.DB instance.
//func (cn *ADBMSSql) SQLDB() *sql.DB {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//	return cn.sqldb
//}
//
//// GetSandboxAdapter returns a sandbox adapter for the MSSQL database.
//// It ensures that the ADBMSSql instance and its SQLDB are properly initialized.
//func (cn *ADBMSSql) GetSandboxAdapter() (aconns.ISBAdapter, error) {
//	if cn == nil {
//		return nil, fmt.Errorf("no mssql db has been created")
//	}
//	if cn.SQLDB() == nil {
//		return nil, fmt.Errorf("no mssql db has been created where host=%s", cn.GetHost())
//	}
//	return aconns.NewSBAdapterSql(cn, cn.SQLDB()), nil
//}
//
//// ADBMSSqls represents a slice of ADBMSSql pointers.
//type ADBMSSqls []*ADBMSSql
