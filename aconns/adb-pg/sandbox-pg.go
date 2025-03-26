package adb_pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/uptrace/bun"
	"strings"
)

type SandboxPGS struct {
	db      *bun.DB
	helper  aconns.ISBAdapterHelper
	adapter *ADBPG
}

// GetType returns the adapter type.
func (sba *SandboxPGS) GetType() aconns.AdapterType {
	return sba.adapter.GetType()
}

// GetName returns the adapter name.
func (sba *SandboxPGS) GetName() aconns.AdapterName {
	return sba.adapter.GetName()
}

// GetHost returns the adapter host.
func (sba *SandboxPGS) GetHost() string {
	return sba.adapter.GetHost()
}

// SupportsModels returns false as models are not supported.
func (sba *SandboxPGS) SupportsModels() bool {
	return true
}

// Query executes a query and returns the result rows.
func (sba *SandboxPGS) Query(query string) (aconns.ISBAdapterSqlRows, error) {
	return nil, fmt.Errorf("only supports QueryModel")
}

// QueryArgs executes a query with arguments and returns the result rows.
func (sba *SandboxPGS) QueryArgs(query string, args ...interface{}) (result aconns.ISBAdapterSqlRows, err error) {
	return nil, fmt.Errorf("only supports QueryModelArgs")
}

// QueryModel returns an error as only QueryArgs is supported.
func (sba *SandboxPGS) QueryModel(query string, model interface{}) error {
	return sba.QueryModelArgs(query, model, nil)
}

// QueryModelArgs returns an error as only QueryArgs is supported.
func (sba *SandboxPGS) QueryModelArgs(query string, model interface{}, args ...interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()
	if sba == nil || sba.db == nil {
		return fmt.Errorf("db is nil")
	}
	query = strings.TrimSpace(query)
	if query == "" {
		return fmt.Errorf("empty query")
	}
	var q *bun.RawQuery
	// the final check ("args[0] == nil") is for users who accidentally add a nil for an arg
	if args == nil || len(args) == 0 || args[0] == nil {
		q = sba.db.NewRaw(query)
	} else {
		q = sba.db.NewRaw(query, args...)
	}
	if err = q.Scan(context.Background(), model); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}

func (sba *SandboxPGS) RunCommand(text string) error {
	if sba == nil || sba.db == nil {
		return fmt.Errorf("db is nil")
	}
	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("empty text")
	}
	_, err := sba.db.Exec(text)
	return err
}

func (sba *SandboxPGS) RunMapByAction() error {
	if sba.helper == nil {
		return fmt.Errorf("helper is nil")
	}
	return sba.RunCommand(sba.helper.MustGetByAction())
}

func (sba *SandboxPGS) RunSqlMapAction(action aconns.ConnActionType) error {
	if sba.helper == nil {
		return fmt.Errorf("helper is nil")
	}
	if action.IsEmpty() {
		return fmt.Errorf("empty action")
	}
	return sba.RunCommand(sba.helper.MustGet(action))
}

func (sba *SandboxPGS) GetAdapterHelper() aconns.ISBAdapterHelper {
	if sba.helper == nil {
		return nil
	}
	return sba.helper
}
