package aconns

import (
	"database/sql" // Package sql provides a generic interface around SQL (or SQL-like) databases.
	"fmt"
	"strings"
)

// ISBAdapterSql is for sandboxed adapters with SQL capability.
type ISBAdapterSql interface {
	ISBAdapter

	// SupportsModels returns 'true' if Model(s) are supported,
	// otherwise use Query to return a RowScanner.
	SupportsModels() bool

	// Query operates similar to "database/sql".
	Query(query string) (ISBAdapterSqlRows, error)
	QueryArgs(query string, args ...interface{}) (ISBAdapterSqlRows, error)

	// QueryModel operates similar to "github.com/uptrace/bun".
	QueryModel(query string, model interface{}) error
	QueryModelArgs(query string, model interface{}, args ...interface{}) error
}

// SBAdapterSql provides SQL capabilities for sandboxed adapters.
type SBAdapterSql struct {
	adapter IAdapterDB
	db      *sql.DB
}

// NewSBAdapterSql creates a new SBAdapterSql instance.
func NewSBAdapterSql(adapter IAdapterDB, db *sql.DB) *SBAdapterSql {
	return &SBAdapterSql{adapter: adapter, db: db}
}

// GetType returns the adapter type.
func (sba *SBAdapterSql) GetType() AdapterType {
	return sba.adapter.GetType()
}

// GetName returns the adapter name.
func (sba *SBAdapterSql) GetName() AdapterName {
	return sba.adapter.GetName()
}

// GetHost returns the adapter host.
func (sba *SBAdapterSql) GetHost() string {
	return sba.adapter.GetHost()
}

// SupportsModels returns false as models are not supported.
func (sba *SBAdapterSql) SupportsModels() bool {
	return false
}

// Query executes a query and returns the result rows.
func (sba *SBAdapterSql) Query(query string) (ISBAdapterSqlRows, error) {
	return sba.QueryArgs(query, nil)
}

// QueryArgs executes a query with arguments and returns the result rows.
func (sba *SBAdapterSql) QueryArgs(query string, args ...interface{}) (result ISBAdapterSqlRows, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()
	if sba == nil || sba.db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	//if err := sba.db.Ping(); err != nil {
	//	return nil, fmt.Errorf("db is not connected: %v", err)
	//}
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("empty query")
	}
	rows, err := sba.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return &SBAdapterSqlRows{rows: rows}, nil
}

// QueryModel returns an error as only QueryArgs is supported.
func (sba *SBAdapterSql) QueryModel(query string, model interface{}) error {
	return fmt.Errorf("only supports QueryArgs")
}

// QueryModelArgs returns an error as only QueryArgs is supported.
func (sba *SBAdapterSql) QueryModelArgs(query string, model interface{}, args ...interface{}) error {
	return fmt.Errorf("only supports QueryArgs")
}
