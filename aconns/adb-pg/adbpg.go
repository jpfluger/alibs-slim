package adb_pg

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/jpfluger/alibs-slim/autils"
	"github.com/uptrace/bun"                   // Bun is a SQL-first Golang ORM for PostgreSQL, MySQL, MSSQL, and SQLite.
	"github.com/uptrace/bun/dialect/pgdialect" // PostgreSQL dialect for Bun ORM.
	"github.com/uptrace/bun/driver/pgdriver"   // PostgreSQL driver for Bun ORM.
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ADAPTERTYPE_PG        aconns.AdapterType = "pg"
	POSTGRES_DEFAULT_PORT                    = 5432
)

// ADBPG represents a PostgreSQL adapter.
type ADBPG struct {
	aconns.ADBAdapterBase

	DialTimeout  int `json:"dialTimeout,omitempty"`
	ReadTimeout  int `json:"readTimeout,omitempty"`
	WriteTimeout int `json:"writeTimeout,omitempty"`
	PingTimeOut  int `json:"pingTimeOut,omitempty"`

	TLSType string `json:"tlsType,omitempty"`

	QueryHook PGQueryHook `json:"queryHook,omitempty"`

	db *bun.DB

	queryDebug *QueryHook

	mu sync.RWMutex
}

// validate checks if the ADBPG object is valid.
func (cn *ADBPG) validate() error {
	if err := cn.ADBAdapterBase.Validate(); err != nil {
		return err
	}

	if cn.DialTimeout <= 0 {
		cn.DialTimeout = 5
	}
	if cn.ReadTimeout <= 0 {
		cn.ReadTimeout = 30
	}
	if cn.WriteTimeout <= 0 {
		cn.WriteTimeout = 30
	}
	if cn.PingTimeOut <= 0 {
		cn.PingTimeOut = 5
	}

	cn.Host = strings.TrimSpace(cn.Host)
	if cn.Host == "" {
		cn.Host = "localhost"
	}

	if cn.Port <= 0 {
		cn.Port = POSTGRES_DEFAULT_PORT
	}

	cn.Database = autils.ToStringTrimLower(cn.Database)

	return nil
}

// Validate checks if the ADBPG object is valid.
func (cn *ADBPG) Validate() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.validate()
}

func (cn *ADBPG) WaitForPostgresReady(retries int, delay time.Duration) error {
	for i := 0; i < retries; i++ {
		if _, _, err := cn.Test(); err == nil {
			return nil
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("PostgreSQL is not ready after %d retries", retries)
}

// Test attempts to validate the ADBPG, open a connection if necessary, and test the connection.
func (cn *ADBPG) Test() (bool, aconns.TestStatus, error) {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	if err := cn.validate(); err != nil {
		return false, aconns.TESTSTATUS_FAILED, err
	}

	if cn.db != nil {
		if err := cn.testConnection(cn.db); err == nil {
			return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
		}
	}

	if err := cn.openConnection(); err != nil {
		return false, aconns.TESTSTATUS_FAILED, err
	}

	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

// OpenConnection opens a connection to the PostgreSQL database.
func (cn *ADBPG) OpenConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.openConnection()
}

// openConnection opens a connection to the PostgreSQL database.
func (cn *ADBPG) openConnection() error {
	sqldb := sql.OpenDB(cn.getConnConfig())

	ii := 0
	var errPing error
	for ii < cn.PingTimeOut {
		time.Sleep(1 * time.Second)
		ii++
		if err := sqldb.Ping(); err != nil {
			errPing = fmt.Errorf("could not ping new conn where host=%s; %v", cn.GetHost(), err)
			continue
		} else {
			errPing = nil
			break
		}
	}

	if errPing != nil {
		return errPing
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	if err := cn.testConnection(db); err != nil {
		return err
	}

	cn.queryDebug = NewQueryHook(
		QueryHookOptionWithEnabled(cn.QueryHook.IsEnabled),
		QueryHookOptionWithVerbose(cn.QueryHook.IsVerbose),
	)

	db.AddQueryHook(cn.queryDebug)

	cn.db = db
	return nil
}

// getConnConfig returns the PostgreSQL connection configuration.
func (cn *ADBPG) getConnConfig() *pgdriver.Connector {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(cn.getAddress()),
		pgdriver.WithUser(cn.GetUsername()),
		pgdriver.WithPassword(cn.GetPassword()),
		pgdriver.WithDatabase(cn.GetDatabase()),
	)

	pgdriver.WithDialTimeout(time.Duration(cn.DialTimeout) * time.Second)(pgconn.Config())
	pgdriver.WithReadTimeout(time.Duration(cn.ReadTimeout) * time.Second)(pgconn.Config())
	pgdriver.WithWriteTimeout(time.Duration(cn.WriteTimeout) * time.Second)(pgconn.Config())

	switch autils.ToStringTrimLower(cn.TLSType) {
	case "verify-ca", "verify-full":
		pgdriver.WithTLSConfig(new(tls.Config))(pgconn.Config())
	case "allow", "prefer", "require":
		pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true})(pgconn.Config())
	default:
		pgdriver.WithInsecure(true)(pgconn.Config())
	}

	return pgconn
}

// GetAddress returns the address of the PostgreSQL server.
func (cn *ADBPG) GetAddress() string {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.getAddress()
}

// getAddress returns the address of the PostgreSQL server.
func (cn *ADBPG) getAddress() string {
	port := cn.GetPort()
	return fmt.Sprintf("%s:%s", cn.GetHost(), strconv.Itoa(port))
}

// testConnection tests the connection to the PostgreSQL database.
func (cn *ADBPG) testConnection(db *bun.DB) error {
	if db == nil {
		return fmt.Errorf("no postgres db has been created where host=%s", cn.GetHost())
	}

	var isValid int
	if err := db.NewRaw("SELECT 1").Scan(context.Background(), &isValid); err != nil {
		return fmt.Errorf("test connection failed for pg where host=%s; %v", cn.GetHost(), err)
	}

	return nil
}

// CloseConnection closes the connection to the PostgreSQL database.
func (cn *ADBPG) CloseConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	if cn.db != nil {
		cn.db.Close()
		cn.db = nil
	}
	return nil
}

// DB returns the bun.DB instance.
func (cn *ADBPG) DB() *bun.DB {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.db
}

// Escape escapes a string for use in SQL queries.
func Escape(name string) string {
	if name == "" {
		return name
	}
	return fmt.Sprintf(`"%s"`, strings.ReplaceAll(name, `"`, `""`))
}

// Count returns the count of records in the model.
func (cn *ADBPG) Count(model interface{}) (int, error) {
	return cn.DB().NewSelect().Model(model).Count(context.Background())
}

// SelectAll selects all records from the model.
func (cn *ADBPG) SelectAll(model interface{}) error {
	return cn.DB().NewSelect().
		Model(model).
		Scan(context.Background())
}

// Truncate deletes records from the table based on the model.
// If doCascadeAll is true, it adds the CASCADE option to truncate all dependent tables
// otherwise `ON DELETE CASCADE` must be applied to the table itself.
func (cn *ADBPG) Truncate(model interface{}, doCascadeAll bool) (sql.Result, error) {
	truncateQuery := cn.DB().NewTruncateTable().Model(model)

	if doCascadeAll {
		truncateQuery.Cascade()
	}

	return truncateQuery.Exec(context.Background())
}

// GetSandboxAdapter returns a sandbox adapter for the PG database.
func (cn *ADBPG) GetSandboxAdapter() (aconns.ISBAdapter, error) {
	return cn.GetSandboxAdapterWithHelper(nil)
}

func (cn *ADBPG) GetSandboxAdapterWithHelper(helper aconns.ISBAdapterHelper) (aconns.ISBAdapter, error) {
	if cn == nil {
		return nil, fmt.Errorf("no pg db has been created")
	}
	if cn.DB() == nil {
		return nil, fmt.Errorf("no pg db has been created where adapter=%s", cn.GetName().String())
	}
	return &SandboxPGS{
		db:      cn.DB(),
		adapter: cn,
		helper:  helper,
	}, nil
}

// ExecuteSQLFile reads and executes an SQL file as a single command.
func (cn *ADBPG) ExecuteSQLFile(filePath string) error {
	// Read the file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}
	db := cn.DB()
	if db == nil {
		return fmt.Errorf("no pg db has been created")
	}
	// Execute the entire SQL content
	_, err = db.ExecContext(context.Background(), string(data))
	if err != nil {
		return fmt.Errorf("failed to execute SQL file: %w", err)
	}
	return nil
}

// ExecuteSQL executes the given string in a single command.
func (cn *ADBPG) ExecuteSQL(command string) error {
	if command == "" {
		return fmt.Errorf("SQL command is empty")
	}
	db := cn.DB()
	if db == nil {
		return fmt.Errorf("no pg db has been created")
	}
	// Execute the entire SQL content
	_, err := db.ExecContext(context.Background(), command)
	if err != nil {
		return fmt.Errorf("failed to execute SQL command: %w", err)
	}
	return nil
}
